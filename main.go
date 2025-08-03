package configwatcher

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/blackorder/chanhub"
	"github.com/fsnotify/fsnotify"
)

// Option configures a Watcher. Use WithErrorChan to receive internal errors.
type Option[T any] func(*Watcher[T])

// WithErrorChan sets an error channel to receive load/save errors.
func WithErrorChan[T any](ch chan<- error) Option[T] {
	return func(w *Watcher[T]) { w.errChan = ch }
}

// Watcher[T] watches a file for type T, broadcasts updates, and reports errors.
type Watcher[T any] struct {
	hub      *chanhub.Hub
	value    atomic.Value
	filename string
	errChan  chan<- error
	fsw      *fsnotify.Watcher
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewWatcher creates a Watcher with defaultVal, file path, and optional settings.
func NewWatcher[T any](defaultVal T, filename string, opts ...Option[T]) *Watcher[T] {
	absFile, _ := filepath.Abs(filename)
	w := &Watcher[T]{
		hub:      chanhub.New(),
		filename: absFile,
	}
	w.value.Store(defaultVal)
	w.load()
	for _, opt := range opts {
		opt(w)
	}

	// start fsnotify watcher
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		w.sendError(err)
	} else {
		w.fsw = fsw
		dir := filepath.Dir(absFile)
		if err := w.fsw.Add(dir); err != nil {
			w.sendError(err)
		}
		w.ctx, w.cancel = context.WithCancel(context.Background())
		go w.watchFS()
	}
	return w
}

// Get returns the current config value.
func (w *Watcher[T]) Get() T {
	return w.value.Load().(T)
}

// Subscribe returns a channel that signals when the config reloads.
func (w *Watcher[T]) Subscribe(ctx context.Context) <-chan struct{} {
	return w.hub.Subscribe(ctx)
}

// Save writes cfg to disk and reloads. Returns any write or marshal error.
func (w *Watcher[T]) Save(cfg T) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		w.sendError(err)
		return err
	}
	if err = os.WriteFile(w.filename, data, 0o600); err != nil {
		w.sendError(err)
		return err
	}
	w.load()
	return nil
}

// watchFS listens for fsnotify events and reloads on relevant changes.
func (w *Watcher[T]) watchFS() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case ev, ok := <-w.fsw.Events:
			if !ok {
				return
			}
			if ev.Name == w.filename && (ev.Op&fsnotify.Write == fsnotify.Write || ev.Op&fsnotify.Create == fsnotify.Create) {
				w.load()
			}
		case err, ok := <-w.fsw.Errors:
			if !ok {
				return
			}
			w.sendError(err)
		}
	}
}

// load reads the file, unmarshals into T, updates on change, and broadcasts.
func (w *Watcher[T]) load() {
	data, err := os.ReadFile(w.filename)
	if err != nil {
		w.sendError(err)
		_ = w.writeFile(w.Get())
		return
	}
	if len(data) == 0 {
		_ = w.writeFile(w.Get())
		return
	}
	var newVal T
	if err := json.Unmarshal(data, &newVal); err != nil {
		w.sendError(err)
		return
	}
	cur := w.Get()
	if !equal(cur, newVal) {
		w.value.Store(newVal)
		w.hub.Broadcast()
	}
}

// writeFile persists cfg without reloading.
func (w *Watcher[T]) writeFile(cfg T) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		w.sendError(err)
		return err
	}
	if err := os.WriteFile(w.filename, data, 0o600); err != nil {
		w.sendError(err)
		return err
	}
	return nil
}

// sendError non-blockingly emits errors to the provided channel.
func (w *Watcher[T]) sendError(err error) {
	if w.errChan == nil || err == nil {
		return
	}
	select {
	case w.errChan <- err:
	default:
	}
}

// equal performs a deep equality check via JSON round-trip.
func equal[T any](a, b T) bool {
	ar, _ := json.Marshal(a)
	br, _ := json.Marshal(b)
	return bytes.Equal(ar, br)
}
