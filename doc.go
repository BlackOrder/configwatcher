/*
Package configwatcher provides type-safe configuration file watching with automatic reloading.

ConfigWatcher is designed to monitor configuration files for changes and automatically
reload them while providing type-safe access through Go generics. It supports JSON
configuration files and offers broadcasting capabilities for configuration change
notifications.

# Basic Usage

Create a configuration type and use NewWatcher to start monitoring a file:

	type Config struct {
		AppName string `json:"app_name"`
		Port    int    `json:"port"`
		Debug   bool   `json:"debug"`
	}

	defaultConfig := Config{
		AppName: "myapp",
		Port:    8080,
		Debug:   false,
	}

	watcher := configwatcher.NewWatcher(defaultConfig, "config.json")
	config := watcher.Get()

# Configuration Changes

Subscribe to configuration changes using the Subscribe method:

	ctx := context.Background()
	updateChan := watcher.Subscribe(ctx)

	go func() {
		for range updateChan {
			newConfig := watcher.Get()
			fmt.Printf("Config updated: %+v\n", newConfig)
		}
	}()

# Saving Configuration

Update and save configuration programmatically:

	config := watcher.Get()
	config.Port = 9090
	if err := watcher.Save(config); err != nil {
		log.Printf("Failed to save: %v", err)
	}

# Error Handling

Use WithErrorChan to receive error notifications:

	errChan := make(chan error, 10)
	watcher := configwatcher.NewWatcher(
		defaultConfig,
		"config.json",
		configwatcher.WithErrorChan[Config](errChan),
	)

	go func() {
		for err := range errChan {
			log.Printf("Config error: %v", err)
		}
	}()

# Thread Safety

All operations are thread-safe:
  - Get() uses atomic operations for lock-free reads
  - Save() is internally synchronized
  - Subscribe() uses channels for safe concurrent access
  - File watching runs in a separate goroutine

# File Format

Configuration files must be valid JSON matching your struct definition:

	{
		"app_name": "myapp",
		"port": 8080,
		"debug": true
	}

If the file doesn't exist, it will be created with the default configuration.
If the file contains invalid JSON, errors will be reported via the error channel
and the current configuration will be preserved.
*/
package configwatcher
