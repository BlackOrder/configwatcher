# ConfigWatcher

[![CI](https://github.com/blackorder/configwatcher/actions/workflows/ci.yml/badge.svg)](https://github.com/blackorder/configwatcher/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/blackorder/configwatcher/branch/main/graph/badge.svg)](https://codecov.io/gh/blackorder/configwatcher)
[![Go Report Card](https://goreportcard.com/badge/github.com/blackorder/configwatcher)](https://goreportcard.com/report/github.com/blackorder/configwatcher)
[![GoDoc](https://godoc.org/github.com/blackorder/configwatcher?status.svg)](https://godoc.org/github.com/blackorder/configwatcher)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A Go library package for watching configuration files with automatic reloading, type safety, and broadcasting capabilities. Built with generics for type-safe configuration management.

## Features

- **Type-safe configuration**: Uses Go generics for compile-time type safety
- **Automatic file watching**: Monitors configuration files for changes using fsnotify
- **Hot reloading**: Automatically reloads configuration when files change
- **Broadcast notifications**: Subscribe to configuration change events
- **Error handling**: Optional error channel for handling load/save errors
- **JSON support**: Built-in JSON marshaling/unmarshaling
- **Thread-safe**: Concurrent access safe with atomic operations
- **Default value handling**: Gracefully handles missing or malformed files

## Installation

```bash
go get github.com/blackorder/configwatcher
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/blackorder/configwatcher"
)

type AppConfig struct {
    AppName string `json:"app_name"`
    Port    int    `json:"port"`
    Debug   bool   `json:"debug"`
}

func main() {
    // Define default configuration
    defaultConfig := AppConfig{
        AppName: "myapp",
        Port:    8080,
        Debug:   false,
    }

    // Create error channel for handling errors
    errChan := make(chan error, 10)

    // Create watcher
    watcher := configwatcher.NewWatcher(
        defaultConfig,
        "config.json",
        configwatcher.WithErrorChan[AppConfig](errChan),
    )

    // Handle errors in a separate goroutine
    go func() {
        for err := range errChan {
            log.Printf("Config error: %v", err)
        }
    }()

    // Get current configuration
    config := watcher.Get()
    fmt.Printf("Current config: %+v\n", config)

    // Subscribe to configuration changes
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    updateChan := watcher.Subscribe(ctx)
    go func() {
        for range updateChan {
            newConfig := watcher.Get()
            fmt.Printf("Config updated: %+v\n", newConfig)
        }
    }()

    // Update configuration programmatically
    config.Port = 9090
    config.Debug = true
    if err := watcher.Save(config); err != nil {
        log.Printf("Failed to save config: %v", err)
    }

    // Your application logic here...
    time.Sleep(5 * time.Second)
}
```

## API Documentation

### Types

#### `Watcher[T]`

The main type that watches a configuration file for changes and provides type-safe access to the configuration.

```go
type Watcher[T any] struct {
    // internal fields...
}
```

#### `Option[T]`

A function type for configuring a Watcher instance.

```go
type Option[T any] func(*Watcher[T])
```

### Functions

#### `NewWatcher[T any](defaultVal T, filename string, opts ...Option[T]) *Watcher[T]`

Creates a new Watcher instance that monitors the specified file for changes.

**Parameters:**
- `defaultVal`: The default configuration value to use if the file doesn't exist or is invalid
- `filename`: Path to the configuration file to watch
- `opts`: Optional configuration options

**Returns:**
- A new `*Watcher[T]` instance

**Example:**
```go
watcher := configwatcher.NewWatcher(defaultConfig, "config.json")
```

#### `WithErrorChan[T any](ch chan<- error) Option[T]`

Option function that sets an error channel to receive load/save errors.

**Parameters:**
- `ch`: A channel to receive error notifications

**Returns:**
- An `Option[T]` function

**Example:**
```go
errChan := make(chan error, 10)
watcher := configwatcher.NewWatcher(
    defaultConfig,
    "config.json",
    configwatcher.WithErrorChan[AppConfig](errChan),
)
```

### Methods

#### `(w *Watcher[T]) Get() T`

Returns the current configuration value. This method is thread-safe and non-blocking.

**Returns:**
- The current configuration value of type `T`

**Example:**
```go
config := watcher.Get()
fmt.Printf("Current port: %d\n", config.Port)
```

#### `(w *Watcher[T]) Save(cfg T) error`

Saves the configuration to disk and triggers a reload. The configuration is marshaled to JSON before writing.

**Parameters:**
- `cfg`: The configuration value to save

**Returns:**
- An error if the save operation fails

**Example:**
```go
config := watcher.Get()
config.Port = 9090
if err := watcher.Save(config); err != nil {
    log.Printf("Failed to save: %v", err)
}
```

#### `(w *Watcher[T]) Subscribe(ctx context.Context) <-chan struct{}`

Returns a channel that receives notifications when the configuration changes.

**Parameters:**
- `ctx`: Context for cancellation

**Returns:**
- A receive-only channel that signals configuration changes

**Example:**
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

updateChan := watcher.Subscribe(ctx)
for range updateChan {
    newConfig := watcher.Get()
    fmt.Printf("Config updated: %+v\n", newConfig)
}
```

## Configuration File Format

ConfigWatcher uses JSON format for configuration files. The structure must match your configuration type.

**Example config.json:**
```json
{
  "app_name": "myapp",
  "port": 8080,
  "debug": true,
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "mydb"
  }
}
```

## Error Handling

ConfigWatcher provides several mechanisms for error handling:

### Error Channel

Use `WithErrorChan` to receive error notifications:

```go
errChan := make(chan error, 10)
watcher := configwatcher.NewWatcher(
    defaultConfig,
    "config.json",
    configwatcher.WithErrorChan[AppConfig](errChan),
)

go func() {
    for err := range errChan {
        log.Printf("Config error: %v", err)
    }
}()
```

### Save Errors

The `Save` method returns errors directly:

```go
if err := watcher.Save(config); err != nil {
    log.Printf("Failed to save config: %v", err)
}
```

### Common Error Scenarios

- **File not found**: Uses default value and creates the file
- **Permission denied**: Reported via error channel
- **Invalid JSON**: Reported via error channel, keeps current config
- **Empty file**: Uses default value and populates the file

## Advanced Usage

### Custom Configuration Types

```go
type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
    Database string `json:"database"`
}

type AppConfig struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    Features []string       `json:"features"`
}

watcher := configwatcher.NewWatcher(defaultConfig, "app-config.json")
```

### Multiple Configuration Files

```go
// Database configuration
dbWatcher := configwatcher.NewWatcher(defaultDBConfig, "db.json")

// Server configuration  
serverWatcher := configwatcher.NewWatcher(defaultServerConfig, "server.json")

// Application configuration
appWatcher := configwatcher.NewWatcher(defaultAppConfig, "app.json")
```

### Configuration Validation

```go
type Config struct {
    Port int `json:"port"`
    Host string `json:"host"`
}

func (c Config) Validate() error {
    if c.Port < 1 || c.Port > 65535 {
        return fmt.Errorf("invalid port: %d", c.Port)
    }
    if c.Host == "" {
        return fmt.Errorf("host cannot be empty")
    }
    return nil
}

// Subscribe to changes and validate
updateChan := watcher.Subscribe(ctx)
for range updateChan {
    config := watcher.Get()
    if err := config.Validate(); err != nil {
        log.Printf("Invalid config: %v", err)
        // Optionally revert to a known good configuration
    }
}
```

## Thread Safety

ConfigWatcher is designed to be thread-safe:

- `Get()` uses atomic operations for lock-free reads
- `Save()` is synchronized internally
- `Subscribe()` uses channels for safe concurrent access
- File system watching runs in a separate goroutine

## Performance Considerations

- **Memory**: Configuration values are stored in memory for fast access
- **CPU**: File system events trigger reload only when necessary
- **I/O**: JSON marshaling/unmarshaling occurs only during saves and loads
- **Concurrency**: Multiple goroutines can safely call `Get()` simultaneously

## Requirements

- Go 1.21 or later (for generics support)
- Supported operating systems: Linux, macOS, Windows (via fsnotify)

## Dependencies

- [fsnotify](https://github.com/fsnotify/fsnotify) - Cross-platform file system notifications
- [chanhub](https://github.com/blackorder/chanhub) - Channel broadcasting utilities

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

### v0.1.0
- Initial release
- Basic configuration watching functionality
- JSON marshaling/unmarshaling support
- Type-safe generics implementation
- Error handling via channels
- Thread-safe operations
