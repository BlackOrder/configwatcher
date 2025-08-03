# ConfigWatcher Examples

This directory contains example applications demonstrating various use cases of the ConfigWatcher library.

## Examples

### Basic Example

**Location**: `basic/`

A simple example showing the core functionality of ConfigWatcher:

- Creating a watcher with default configuration
- Subscribing to configuration changes
- Programmatic configuration updates
- Error handling

**Run the example**:
```bash
cd basic
go mod tidy
go run main.go
```

The example will:
1. Create a config watcher for `app-config.json`
2. Display the initial configuration
3. Make programmatic updates
4. Show live updates when you edit the file manually

### Multi-Config Example

**Location**: `multi-config/`

An advanced example demonstrating:

- Multiple configuration watchers
- Different configuration types
- Centralized error handling
- Coordinated configuration updates

**Run the example**:
```bash
cd multi-config
go mod tidy
go run main.go
```

This example manages three separate configuration files:
- `app.json` - Application-level settings
- `server.json` - Server configuration
- `database.json` - Database settings

## Interactive Testing

Both examples create configuration files that you can edit manually while the programs are running to see live configuration reloading in action.

### Try These Experiments

1. **Basic Example**:
   - Start the program: `go run main.go`
   - Edit `app-config.json` in another terminal
   - Watch the console for update notifications

2. **Multi-Config Example**:
   - Start the program: `go run main.go`
   - Edit any of the three JSON files: `app.json`, `server.json`, `database.json`
   - Observe how each file triggers its own update handler

### Configuration File Formats

**Basic Example** (`app-config.json`):
```json
{
  "app_name": "example-app",
  "port": 8080,
  "debug": false,
  "database": {
    "host": "localhost",
    "port": 5432,
    "username": "admin",
    "database": "myapp"
  },
  "features": {
    "feature1": "enabled",
    "feature2": "disabled"
  }
}
```

**Multi-Config Example** (`app.json`):
```json
{
  "environment": "development",
  "log_level": "info",
  "server": {
    "host": "localhost",
    "port": 8080,
    "read_timeout": "30s",
    "write_timeout": "30s",
    "max_clients": 100
  },
  "database": {
    "dsn": "postgres://localhost/myapp",
    "max_conns": 10,
    "max_idle": 5,
    "max_lifetime": "1h"
  }
}
```

## Running Examples

Since this is a library package, examples are meant to be run with `go run`, not built as binaries.

To run all examples:

```bash
# From the root directory
make examples

# Or manually
cd examples/basic && go run main.go
cd examples/multi-config && go run main.go
```

## Cleanup

To remove generated configuration files:

```bash
# From the root directory
make clean

# Or manually
rm -f examples/basic/*.json examples/multi-config/*.json
```

## Learning Points

These examples demonstrate:

1. **Type Safety**: How Go generics provide compile-time type checking
2. **Error Handling**: Different approaches to handling configuration errors
3. **Hot Reloading**: Automatic detection and loading of file changes
4. **Broadcasting**: How multiple subscribers can react to configuration changes
5. **Concurrent Safety**: Safe access from multiple goroutines
6. **Default Handling**: Graceful fallback to default values

## Use in Your Projects

These examples can serve as templates for integrating ConfigWatcher into your own applications. Key patterns to adopt:

- Always provide sensible default configurations
- Set up error channels for production monitoring
- Use context for graceful shutdown of subscribers
- Consider validation after configuration updates
- Structure your configuration types with clear JSON tags
