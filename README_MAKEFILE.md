# Makefile Usage Guide

This project includes a comprehensive Makefile for easy development, testing, and deployment.

## Quick Start

```bash
# See all available commands
make help

# Build and run the application
make run

# Run tests
make test

# Clean everything
make clean-all
```

## Common Tasks

### Development

```bash
# Build the application
make build

# Run with local authentication (default)
make run

# Run with OAuth2 authentication
make run-oauth

# Run with hot reload (auto-restart on file changes)
make dev
```

### Testing

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Run specific test suites
make test-handlers
make test-services
make test-models
```

### Maintenance

```bash
# Format code
make fmt

# Run linters (requires golangci-lint)
make lint

# Run all checks (fmt + test)
make check

# Clean build artifacts
make clean

# Remove database files
make clean-db

# Remove everything (binary + database + temp files)
make clean-all
```

### Database Management

```bash
# Create fresh database with migrations
make migrate

# Seed database with test users
make seed

# Backup database
make backup

# Restore from latest backup
make restore
```

### Dependencies

```bash
# Install dependencies
make install

# Download and verify dependencies
make deps
```

### Project Information

```bash
# Show project statistics
make info
```

## Test Users (after `make seed`)

After running `make seed`, the following test accounts are available:

- **Instructor**: `instructor1` / `password123`
- **Student 1**: `student1` / `password123`
- **Student 2**: `student2` / `password123`

## Development Workflow

### Typical Development Session

```bash
# 1. Clean start
make clean-all

# 2. Build and test
make build
make test

# 3. Run with hot reload
make dev

# 4. Before committing
make check
```

### Pre-Commit Checklist

```bash
# Format, lint, and test in one command
make check
```

## Advanced Usage

### Custom Port

```bash
# Run on custom port
PORT=3000 make run
```

### Hot Reload Setup

The `make dev` command automatically installs and uses [air](https://github.com/cosmtrek/air) for hot reload:

```bash
# Install air manually (if needed)
go install github.com/cosmtrek/air@latest

# Run with hot reload
make dev
```

### Coverage Reports

```bash
# Generate coverage report
make test-coverage

# Open coverage.html in browser
open coverage.html
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Test
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: make deps
      - run: make check
```

## Troubleshooting

### Command not found

If you see "make: command not found":

```bash
# macOS
brew install make

# Ubuntu/Debian
sudo apt-get install build-essential
```

### Port already in use

```bash
# Find and kill process on port 8080
lsof -ti:8080 | xargs kill -9

# Or use custom port
PORT=3000 make run
```

### Database locked

```bash
# Stop any running instances
pkill -f zipcodereader

# Clean and rebuild
make clean-all
make build
```

## Makefile Structure

The Makefile is organized into sections:

1. **Build & Run**: `build`, `run`, `run-oauth`, `dev`
2. **Testing**: `test`, `test-verbose`, `test-coverage`, `test-*`
3. **Cleaning**: `clean`, `clean-db`, `clean-all`
4. **Dependencies**: `install`, `deps`
5. **Code Quality**: `fmt`, `lint`, `check`
6. **Database**: `migrate`, `seed`, `backup`, `restore`
7. **Info**: `help`, `info`

## Tips

- Always run `make help` to see available commands
- Use `make check` before committing code
- Use `make dev` for active development with hot reload
- Run `make clean-all` when switching branches
- Use `make seed` to quickly populate test data
- Run `make info` to see project statistics

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8080 | Server port |
| `GO` | go | Go command |
| `GOFLAGS` | -v | Go build flags |

## Future Enhancements

Coming soon:
- `make docker-build` - Build Docker image
- `make docker-run` - Run in Docker container
- `make deploy` - Deploy to production
- `make bench` - Run benchmarks

---

For more information, see the main [README.md](README.md) or run `make help`.
