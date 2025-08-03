package configwatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// TestConfig represents a sample configuration structure for testing
type TestConfig struct {
	Name     string            `json:"name"`
	Count    int               `json:"count"`
	Settings map[string]string `json:"settings"`
}

func createTempConfigFile(t *testing.T, cfg TestConfig) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "configwatcher_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	configFile := filepath.Join(tmpDir, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0o600); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return configFile
}

func TestNewWatcher(t *testing.T) {
	defaultConfig := TestConfig{
		Name:  "test",
		Count: 42,
		Settings: map[string]string{
			"key1": "value1",
		},
	}

	configFile := createTempConfigFile(t, defaultConfig)

	watcher := NewWatcher(defaultConfig, configFile)
	if watcher == nil {
		t.Fatal("NewWatcher returned nil")
	}

	// Test that the default value is loaded
	got := watcher.Get()
	if got.Name != defaultConfig.Name || got.Count != defaultConfig.Count {
		t.Errorf("Expected %+v, got %+v", defaultConfig, got)
	}
}

func TestNewWatcherWithOptions(t *testing.T) {
	defaultConfig := TestConfig{Name: "test", Count: 1}
	configFile := createTempConfigFile(t, defaultConfig)

	errChan := make(chan error, 10)
	watcher := NewWatcher(defaultConfig, configFile, WithErrorChan[TestConfig](errChan))

	if watcher == nil {
		t.Fatal("NewWatcher returned nil")
	}

	// Verify error channel is set
	if watcher.errChan == nil {
		t.Error("Error channel not set")
	}
}

func TestWatcherGet(t *testing.T) {
	defaultConfig := TestConfig{
		Name:  "initial",
		Count: 100,
	}

	configFile := createTempConfigFile(t, defaultConfig)
	watcher := NewWatcher(defaultConfig, configFile)

	got := watcher.Get()
	if got.Name != "initial" || got.Count != 100 {
		t.Errorf("Expected name='initial', count=100, got name='%s', count=%d", got.Name, got.Count)
	}
}

func TestWatcherSave(t *testing.T) {
	defaultConfig := TestConfig{Name: "test", Count: 1}
	configFile := createTempConfigFile(t, defaultConfig)

	watcher := NewWatcher(defaultConfig, configFile)

	newConfig := TestConfig{
		Name:  "updated",
		Count: 999,
		Settings: map[string]string{
			"new": "setting",
		},
	}

	err := watcher.Save(newConfig)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify the value was updated
	got := watcher.Get()
	if got.Name != "updated" || got.Count != 999 {
		t.Errorf("Expected updated config, got %+v", got)
	}

	// Verify file was written correctly
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	var fileConfig TestConfig
	if err := json.Unmarshal(data, &fileConfig); err != nil {
		t.Fatalf("Failed to unmarshal file config: %v", err)
	}

	if fileConfig.Name != "updated" || fileConfig.Count != 999 {
		t.Errorf("File config not updated correctly: %+v", fileConfig)
	}
}

func TestWatcherSubscribe(t *testing.T) {
	defaultConfig := TestConfig{Name: "test", Count: 1}
	configFile := createTempConfigFile(t, defaultConfig)

	watcher := NewWatcher(defaultConfig, configFile)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateChan := watcher.Subscribe(ctx)

	// Save a new config in a goroutine
	go func() {
		time.Sleep(100 * time.Millisecond)
		newConfig := TestConfig{Name: "updated", Count: 2}
		watcher.Save(newConfig)
	}()

	// Wait for notification
	select {
	case <-updateChan:
		// Success - we got the update notification
		got := watcher.Get()
		if got.Name != "updated" {
			t.Errorf("Expected updated config, got %+v", got)
		}
	case <-ctx.Done():
		t.Fatal("Timeout waiting for config update notification")
	}
}

func TestWatcherFileSystemWatch(t *testing.T) {
	defaultConfig := TestConfig{Name: "test", Count: 1}
	configFile := createTempConfigFile(t, defaultConfig)

	watcher := NewWatcher(defaultConfig, configFile)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	updateChan := watcher.Subscribe(ctx)

	// Modify file externally
	go func() {
		time.Sleep(100 * time.Millisecond)
		newConfig := TestConfig{Name: "external", Count: 999}
		data, _ := json.MarshalIndent(newConfig, "", "  ")
		os.WriteFile(configFile, data, 0o600)
	}()

	// Wait for notification
	select {
	case <-updateChan:
		// Success - filesystem watcher detected the change
		got := watcher.Get()
		if got.Name != "external" || got.Count != 999 {
			t.Errorf("Expected external config update, got %+v", got)
		}
	case <-ctx.Done():
		t.Fatal("Timeout waiting for filesystem update notification")
	}
}

func TestWatcherErrorHandling(t *testing.T) {
	defaultConfig := TestConfig{Name: "test", Count: 1}

	// Use a non-existent directory to trigger errors
	configFile := "/non/existent/path/config.json"

	errChan := make(chan error, 10)
	watcher := NewWatcher(defaultConfig, configFile, WithErrorChan[TestConfig](errChan))

	// Should have received an error for non-existent file
	select {
	case err := <-errChan:
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	case <-time.After(1 * time.Second):
		// It's okay if no error is received immediately
	}

	// Verify default value is still available
	got := watcher.Get()
	if got.Name != "test" {
		t.Errorf("Expected default config, got %+v", got)
	}
}

func TestWatcherConcurrentAccess(t *testing.T) {
	defaultConfig := TestConfig{Name: "concurrent", Count: 0}
	configFile := createTempConfigFile(t, defaultConfig)

	watcher := NewWatcher(defaultConfig, configFile)

	var wg sync.WaitGroup

	// Start multiple goroutines reading and writing
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < 10; j++ {
				// Read
				config := watcher.Get()
				if config.Name != "concurrent" {
					t.Errorf("Unexpected config name: %s", config.Name)
				}

				// Write
				newConfig := TestConfig{
					Name:  "concurrent",
					Count: id*10 + j,
				}
				watcher.Save(newConfig)

				time.Sleep(time.Millisecond)
			}
		}(i)
	}

	wg.Wait()
}

func TestWatcherNonExistentFile(t *testing.T) {
	defaultConfig := TestConfig{Name: "default", Count: 42}

	tmpDir, err := os.MkdirTemp("", "configwatcher_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Use a file that doesn't exist
	configFile := filepath.Join(tmpDir, "nonexistent.json")

	errChan := make(chan error, 10)
	watcher := NewWatcher(defaultConfig, configFile, WithErrorChan[TestConfig](errChan))

	// Should use default value
	got := watcher.Get()
	if got.Name != "default" || got.Count != 42 {
		t.Errorf("Expected default config, got %+v", got)
	}

	// File should be created with default value
	time.Sleep(100 * time.Millisecond) // Give time for file creation

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("Config file should have been created")
	}
}

func TestWatcherMalformedJSON(t *testing.T) {
	defaultConfig := TestConfig{Name: "default", Count: 1}

	tmpDir, err := os.MkdirTemp("", "configwatcher_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "config.json")

	// Write malformed JSON
	malformedJSON := `{"name": "test", "count": invalid}`
	if err := os.WriteFile(configFile, []byte(malformedJSON), 0o600); err != nil {
		t.Fatalf("Failed to write malformed JSON: %v", err)
	}

	errChan := make(chan error, 10)
	watcher := NewWatcher(defaultConfig, configFile, WithErrorChan[TestConfig](errChan))

	// Give some time for the watcher to process the file
	time.Sleep(200 * time.Millisecond)

	// Should either receive unmarshal error or fall back to default
	select {
	case err := <-errChan:
		if err == nil {
			t.Error("Expected unmarshal error")
		}
		// Got the error as expected
	case <-time.After(1 * time.Second):
		// No error received, which is also acceptable if it falls back to default
	} // Should fall back to default value
	got := watcher.Get()
	if got.Name != "default" {
		t.Errorf("Expected default config, got %+v", got)
	}
}

func TestWatcherEmptyFile(t *testing.T) {
	defaultConfig := TestConfig{Name: "default", Count: 1}

	tmpDir, err := os.MkdirTemp("", "configwatcher_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "config.json")

	// Create empty file
	if err := os.WriteFile(configFile, []byte(""), 0o600); err != nil {
		t.Fatalf("Failed to write empty file: %v", err)
	}

	watcher := NewWatcher(defaultConfig, configFile)

	// Should use default value and write it to file
	got := watcher.Get()
	if got.Name != "default" {
		t.Errorf("Expected default config, got %+v", got)
	}

	// File should now contain default config
	time.Sleep(100 * time.Millisecond)
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	if len(data) == 0 {
		t.Error("Config file should have been populated with default config")
	}
}

// Benchmark tests
func BenchmarkWatcherGet(b *testing.B) {
	config := TestConfig{Name: "benchmark", Count: 1}
	tmpDir, _ := os.MkdirTemp("", "configwatcher_bench")
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "config.json")
	data, _ := json.MarshalIndent(config, "", "  ")
	os.WriteFile(configFile, data, 0o600)

	watcher := NewWatcher(config, configFile)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = watcher.Get()
		}
	})
}

func BenchmarkWatcherSave(b *testing.B) {
	config := TestConfig{Name: "benchmark", Count: 1}
	tmpDir, _ := os.MkdirTemp("", "configwatcher_bench")
	defer os.RemoveAll(tmpDir)

	configFile := filepath.Join(tmpDir, "config.json")
	watcher := NewWatcher(config, configFile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newConfig := TestConfig{Name: "benchmark", Count: i}
		watcher.Save(newConfig)
	}
}

// Example test demonstrating usage
func ExampleWatcher() {
	type Config struct {
		AppName string `json:"app_name"`
		Port    int    `json:"port"`
	}

	defaultConfig := Config{
		AppName: "myapp",
		Port:    8080,
	}

	// Create a temporary file for this example
	tmpDir, _ := os.MkdirTemp("", "example")
	defer os.RemoveAll(tmpDir)
	configFile := filepath.Join(tmpDir, "config.json")

	// Create watcher with error channel
	errChan := make(chan error, 10)
	watcher := NewWatcher(defaultConfig, configFile, WithErrorChan[Config](errChan))

	// Get current config
	config := watcher.Get()
	fmt.Printf("Current port: %d\n", config.Port)

	// Save updated config
	config.Port = 9090
	err := watcher.Save(config)
	if err != nil {
		fmt.Printf("Save error: %v\n", err)
	}

	// Get updated config
	updatedConfig := watcher.Get()
	fmt.Printf("Updated port: %d\n", updatedConfig.Port)

	// Output:
	// Current port: 8080
	// Updated port: 9090
}
