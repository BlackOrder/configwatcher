# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2025-08-04

### Added
- Initial release of ConfigWatcher
- Type-safe configuration watching with Go generics
- Automatic file system monitoring using fsnotify
- JSON marshaling/unmarshaling support
- Broadcasting of configuration changes
- Error handling via optional error channels
- Thread-safe operations with atomic values
- Comprehensive test suite with 82%+ coverage
- Detailed documentation and examples
- CI/CD pipeline with GitHub Actions
- Linting and code quality checks

### Features
- `NewWatcher[T]()` - Create a new configuration watcher
- `WithErrorChan[T]()` - Option to set error channel
- `Get()` - Retrieve current configuration value
- `Save()` - Save configuration to disk
- `Subscribe()` - Subscribe to configuration changes
- Automatic file creation with default values
- Graceful handling of malformed JSON
- Support for empty and missing files

### Documentation
- Comprehensive README with examples
- GoDoc documentation for all public APIs
- Contributing guidelines
- Multiple usage examples
- Benchmark results

### Testing
- Unit tests for core functionality
- Integration tests for file system watching
- Concurrent access testing
- Error handling testing
- Performance benchmarks
- Race condition detection

### Security
- File permissions set to 0600 for created files
- Input validation and error handling
- Safe concurrent access patterns

## [0.1.0] - 2025-01-XX

### Added
- Initial implementation
- Core watcher functionality
- Basic test suite
- Documentation

[Unreleased]: https://github.com/blackorder/configwatcher/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/blackorder/configwatcher/compare/v0.1.0...v1.0.0
[0.1.0]: https://github.com/blackorder/configwatcher/releases/tag/v0.1.0
