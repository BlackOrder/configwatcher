package configwatcher

import (
	"context"
	"encoding/json"
	"os"
	"sync/atomic"

	"github.com/blackorder/chanhub"
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
}

// NewWatcher creates a Watcher with defaultVal, file path, and optional settings.
func NewWatcher[T any](defaultVal T, filename string, opts ...Option[T]) *Watcher[T] {
	w := &Watcher[T]{
		hub:      chanhub.New(),
		filename: filename,
	}
	w.value.Store(defaultVal)
	w.load()
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Get returns the current config value.
func (w *Watcher[T]) Get() T {
	return w.value.Load().(T)
}

// Subscribe returns a channel that signals when the config reloads.
func (w *Watcher[T]) Subscribe() <-chan struct{} {
	return w.hub.Subscribe(context.Background())
}

// Save writes cfg to disk and reloads. Returns any write or marshal error.
func (w *Watcher[T]) Save(cfg T) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		w.sendError(err)
		return err
	}
	if err = os.WriteFile(w.filename, data, 0644); err != nil {
		w.sendError(err)
		return err
	}
	w.load()
	return nil
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
	if err := os.WriteFile(w.filename, data, 0644); err != nil {
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
	return string(ar) == string(br)
}
