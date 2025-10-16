# ZipCodeReader

A web-based reading and programming assignment manager for educational environments, built with Go.

## Overview

ZipCodeReader helps instructors assign and track reading materials and programming assignments for students. Students can view their assignments, track progress, and mark work as completed.

## Features

### For Instructors
- Create and manage assignments (reading materials, programming labs)
- Assign work to individual students or groups
- Track student progress and completion rates
- Monitor due dates and overdue assignments
- View detailed analytics and progress reports

### For Students
- View assigned readings and programming tasks
- Track completion progress
- Mark assignments as in-progress or completed
- View due dates and receive notifications
- Submit work via URLs (GitHub PRs, Google Docs, etc.)

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Make (optional, but recommended)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/zipcodereader.git
cd zipcodereader

# Install dependencies
go mod download

# Run the application
make run
```

The application will start on `http://localhost:8080`

### Test Users

After running `make seed`, you can log in with:
- **Instructor**: `instructor1` / `password123`
- **Student**: `student1` / `password123`

## Common Commands

```bash
make help           # Show all available commands
make run            # Run with local authentication
make run-oauth      # Run with GitHub OAuth2
make dev            # Run with hot reload
make test           # Run all tests
make test-coverage  # Generate coverage report
make clean-all      # Clean everything
make seed           # Add test users to database
```

## Authentication Modes

ZipCodeReader supports two authentication modes:

1. **Local Authentication** (default) - Username/password
2. **GitHub OAuth2** - Login with GitHub account

Run with `make run` for local auth or `make run-oauth` for GitHub OAuth2.

## Architecture

```
zipcodereader/
├── main.go                 # Application entry point
├── routes/                 # Route definitions
├── handlers/               # HTTP request handlers
│   ├── instructor.go       # Instructor endpoints
│   ├── student.go          # Student endpoints
│   ├── auth.go            # OAuth2 authentication
│   └── local_auth.go      # Local authentication
├── services/               # Business logic
│   ├── assignment.go       # Assignment management
│   └── auth.go            # Authentication logic
├── models/                 # Data models
│   ├── user.go
│   ├── assignment.go
│   └── student_assignment.go
├── middleware/             # HTTP middleware
├── templates/              # HTML templates
└── static/                 # CSS, JS, images
```

## Data Model

### Core Models

**User**
- ID, Username, Email, Role (student/instructor)
- Password hash (for local auth)
- GitHub integration (for OAuth2)

**Assignment**
- Title, Description, URL, Category
- Type (reading, programming, quiz)
- Estimated time, Due date
- Repository URL (for programming assignments)

**StudentAssignment**
- Assignment → Student relationship
- Status (assigned, in_progress, completed)
- Progress tracking and time spent
- Submission URL

## API Endpoints

### Instructor Endpoints
- `GET/POST /instructor/assignments` - List/create assignments
- `GET/PUT/DELETE /instructor/assignments/:id` - Manage specific assignment
- `POST /instructor/assignments/:id/assign` - Assign to students
- `GET /instructor/assignments/:id/progress` - View progress
- `GET /instructor/dashboard` - Dashboard view

### Student Endpoints
- `GET /student/assignments` - List assigned work
- `GET /student/assignments/:id` - View specific assignment
- `POST /student/assignments/:id/status` - Update status
- `POST /student/assignments/:id/complete` - Mark complete
- `GET /student/dashboard` - Dashboard view

## Development

### Running Tests

```bash
# Run all tests
make test

# Run specific test suites
make test-handlers
make test-services
make test-models

# Generate coverage report
make test-coverage
```

### Hot Reload Development

```bash
# Start with hot reload (auto-restart on file changes)
make dev
```

### Database

ZipCodeReader uses SQLite3 for data storage. The database file is created automatically on first run.

```bash
# Create fresh database
make migrate

# Add test users
make seed

# Backup database
make backup

# Restore from backup
make restore
```

## Documentation

- [Quick Start Guide](docs/quickstart.md)
- [Makefile Commands](docs/makefile.md)
- [Development Log](CLAUDE.md)
- [Architecture History](docs/archive/)

## Project Stats

- **Lines of Code**: ~2,243 (production code)
- **Test Coverage**: 16 comprehensive tests
- **Handler Files**: 4 (instructor, student, auth, local_auth)
- **Service Files**: 2 (assignment, auth)
- **Model Files**: 3 (user, assignment, student_assignment)

## Recent Improvements (October 2025)

Major refactoring completed:
- ✅ 50% reduction in codebase size (4,400 → 2,243 lines)
- ✅ Eliminated all code duplication
- ✅ Consolidated 8 handlers into 2 main handlers
- ✅ Unified 5 services into 1 assignment service
- ✅ Added comprehensive test suite
- ✅ Simplified routing system
- ✅ Enhanced data model for multiple assignment types

## Contributing

Contributions are welcome! Please follow these guidelines:
1. Fork the repository
2. Create a feature branch
3. Run tests with `make test`
4. Submit a pull request

## License

MIT License - See [LICENSE](LICENSE) file for details.

## Support

For questions or issues, please open an issue on GitHub.

---

**Built with Go, Gin, GORM, and SQLite**