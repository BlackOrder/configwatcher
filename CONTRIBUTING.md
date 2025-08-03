# Contributing to ConfigWatcher

Thank you for your interest in contributing to ConfigWatcher! This document provides guidelines for contributing to the project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## Development Setup

### Prerequisites

- Go 1.21 or later
- Git

### Setup

1. Fork the repository on GitHub
2. Clone your fork:
   ```bash
   git clone https://github.com/your-username/configwatcher.git
   cd configwatcher
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Verify everything works:
   ```bash
   go test ./...
   ```

## Development Guidelines

### Code Style

We follow standard Go conventions:

- Use `gofmt` for formatting
- Follow `golangci-lint` recommendations
- Use meaningful variable and function names
- Add comments for exported functions and types
- Keep functions focused and testable

### Testing

- Write tests for all new functionality
- Maintain or improve test coverage
- Include both unit tests and integration tests
- Test error conditions and edge cases
- Use table-driven tests when appropriate

Example test structure:
```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    Input
        expected Expected
        wantErr  bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Documentation

- Update README.md for user-facing changes
- Add Go doc comments for public APIs
- Include examples in documentation
- Update CHANGELOG.md for releases

### Commit Messages

Use conventional commit format:

- `feat: add new feature`
- `fix: resolve bug in X`
- `docs: update documentation`
- `test: add tests for Y`
- `refactor: simplify Z`
- `perf: improve performance of A`
- `chore: update dependencies`

## Submitting Changes

### Pull Request Process

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and commit:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

3. Run tests and linting:
   ```bash
   go test ./...
   golangci-lint run
   ```

4. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

5. Create a Pull Request on GitHub

### Pull Request Guidelines

- Provide a clear description of changes
- Reference any related issues
- Include tests for new functionality
- Ensure CI passes
- Keep PRs focused and reasonably sized
- Update documentation as needed

### Review Process

- All PRs require review from maintainers
- Address review feedback promptly
- Maintain a respectful tone in discussions
- Be open to suggestions and improvements

## Types of Contributions

### Bug Reports

When reporting bugs, please include:

- Go version and OS
- Minimal reproduction case
- Expected vs actual behavior
- Error messages or logs
- Configuration details

### Feature Requests

For feature requests, please provide:

- Clear description of the problem
- Proposed solution
- Use cases and examples
- Consideration of alternatives

### Code Contributions

We welcome contributions for:

- Bug fixes
- Performance improvements
- New features (discuss first in issues)
- Documentation improvements
- Test coverage improvements

## Development Commands

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Run tests with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Run linter
golangci-lint run

# Format code
go fmt ./...

# Test examples (no build artifacts)
cd examples/basic && go run main.go
cd examples/multi-config && go run main.go
```

## Project Structure

```
configwatcher/
├── main.go           # Core library implementation
├── main_test.go      # Test suite
├── doc.go           # Package documentation
├── README.md        # User documentation
├── CONTRIBUTING.md  # This file
├── CHANGELOG.md     # Version history
├── LICENSE          # MIT license
├── Makefile         # Development commands
├── go.mod           # Go module definition
├── go.sum           # Go module checksums
├── .gitignore       # Git ignore patterns
├── .github/         # GitHub workflows
│   └── workflows/
│       └── ci.yml   # CI configuration
├── .golangci.yml    # Linter configuration
└── examples/        # Usage examples (run with go run)
    ├── README.md    # Examples documentation
    ├── basic/       # Basic usage example
    │   ├── go.mod
    │   ├── go.sum
    │   └── main.go
    └── multi-config/ # Multi-config example
        ├── go.mod
        ├── go.sum
        └── main.go
```

## Release Process

1. Update CHANGELOG.md
2. Tag the release: `git tag v1.0.0`
3. Push tags: `git push --tags`
4. GitHub Actions will handle the release

## Getting Help

- Check existing issues and documentation
- Ask questions in GitHub Discussions
- For bugs, open an issue with reproduction steps

## Recognition

Contributors are recognized in:
- README.md contributors section
- Git commit history
- Release notes for significant contributions

Thank you for contributing to ConfigWatcher!
