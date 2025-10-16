# ZipCodeReader Quick Start Guide

A simplified, production-ready web application for managing reading assignments for students and instructors.

## 🚀 Quick Start

### Option 1: Using Make (Recommended)

```bash
# Run the application (local auth mode)
make run

# Visit http://localhost:8080
```

### Option 2: Manual

```bash
# Build
go build -o zipcodereader .

# Run
./zipcodereader
```

### Option 3: Development with Hot Reload

```bash
# Automatically rebuilds when you change Go files
make dev
```

The hot reload feature uses a built-in watch script that monitors your Go files and automatically rebuilds and restarts the server when changes are detected. No additional tools required!

## 📋 Common Commands

```bash
make help           # Show all available commands
make run            # Run application (local auth)
make dev            # Run with hot reload (auto-restart on changes)
make test           # Run all tests
make test-coverage  # Generate coverage report
make clean-all      # Clean everything
make seed           # Add test users
```

## 👥 Test Users (after `make seed`)

- **Instructor**: `instructor1` / `password123`
- **Student**: `student1` / `password123`

## 📚 Documentation

- [Makefile Usage](README_MAKEFILE.md) - Build automation guide
- [CLAUDE.md](CLAUDE.md) - Complete development history
- [Codebase-Simplification-Analysis.md](Codebase-Simplification-Analysis.md) - Architecture analysis

## ✨ Recent Improvements

**October 2025 - Major Refactoring:**
- ✅ 50% smaller codebase (4,400 → 2,243 lines)
- ✅ Eliminated all code duplication
- ✅ Consolidated 8 handlers into 2
- ✅ Unified 5 services into 1
- ✅ 16 comprehensive unit tests
- ✅ Complete build automation

## 🏗️ Architecture

```
zipcodereader/
├── main.go                 (76 lines - simple!)
├── routes/                 (Unified routing)
├── handlers/
│   ├── instructor.go       (All instructor endpoints)
│   └── student.go          (All student endpoints)
├── services/
│   └── assignment.go       (Unified service)
└── models/                 (Data models)
```

## 🔧 Development

### Hot Reload Development

```bash
# Start with automatic restart on file changes
make dev
```

The `make dev` command provides hot reload functionality:
- Automatically rebuilds when you save Go files
- Restarts the server with your changes
- No manual restarts needed
- Uses built-in watch script (no dependencies)

**Alternative hot reload tools** (optional):
- Install `air`: `go install github.com/air-verse/air@latest`
- Install `entr`: `brew install entr`

If you have these installed, `make dev` will use them automatically.

### Testing

```bash
# Run tests
make test

# Pre-commit checks (fmt + test)
make check

# Generate coverage report
make test-coverage
```

## 📦 Dependencies

- Go 1.25+
- SQLite3 (built-in)
- Gin web framework
- GORM ORM

Install/update dependencies:
```bash
make deps
# or
go mod tidy
```

**Optional development tools:**
- `air` - For faster hot reload (optional, built-in alternative available)
- `entr` - Alternative file watcher (optional)
- `golangci-lint` - For linting (optional)

## 🎯 Features

- ✅ Dual authentication (local + OAuth2)
- ✅ Role-based access control
- ✅ Assignment management
- ✅ Progress tracking
- ✅ Due date notifications
- ✅ Student dashboard
- ✅ Instructor dashboard

## 📈 Project Stats

- Lines of Code: 4,203
- Test Files: 7 (16 tests)
- Handler Files: 5
- Service Files: 2
- Model Files: 4

## 🚢 Deployment

See [CLAUDE.md](CLAUDE.md) for deployment instructions.

## 📝 License

MIT License - See LICENSE file for details

---

**Made with ❤️ and simplified with AI assistance**
