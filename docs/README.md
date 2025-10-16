# ZipCodeReader Documentation

Welcome to the ZipCodeReader documentation. This directory contains all project documentation organized by topic.

## Quick Links

- **[Quick Start Guide](quickstart.md)** - Get up and running in 5 minutes
- **[Makefile Commands](makefile.md)** - Build automation and development commands
- **[Development Log](../CLAUDE.md)** - Complete development history and decisions

## Documentation Structure

### Getting Started

- [Quick Start Guide](quickstart.md) - Installation and basic usage
- [Makefile Commands](makefile.md) - Available make commands and workflows

### Architecture & Design

- [Development Log](../CLAUDE.md) - Comprehensive development history
  - Project phases and milestones
  - Technical decisions and rationale
  - Major refactoring efforts
  - Bug fixes and improvements

### Historical Documents

The [archive/](archive/) directory contains historical analysis and planning documents:

- **[Codebase Simplification Analysis](archive/Codebase-Simplification-Analysis.md)** - Detailed analysis of the October 2025 refactoring that reduced codebase by 50%
- **[Data Model Enhancements](archive/Data-Model-Enhancements.md)** - Proposed and implemented data model improvements
- **[Phase 2 Testing](archive/PHASE2_TESTING.md)** - Authentication system testing documentation
- **[Phase 3 Complete](archive/PHASE3_COMPLETE.md)** - Assignment management system completion summary

These documents provide valuable context for understanding the project's evolution and architectural decisions.

## Project Overview

### What is ZipCodeReader?

ZipCodeReader is a web-based assignment management system for educational environments. It enables:

- **Instructors** to create and assign reading materials and programming tasks
- **Students** to view assignments, track progress, and submit work
- **Progress tracking** with analytics and reporting

### Technology Stack

- **Backend**: Go with Gin web framework
- **Database**: SQLite3 with GORM ORM
- **Authentication**: Dual-mode (Local + GitHub OAuth2)
- **Frontend**: HTML templates with Tailwind CSS
- **Session Management**: Gin sessions

### Key Features

- ✅ Role-based access control (instructor/student)
- ✅ Assignment creation and management
- ✅ Progress tracking and analytics
- ✅ Due date notifications
- ✅ Multiple assignment types (reading, programming, quiz)
- ✅ Bulk assignment capabilities
- ✅ Responsive dashboard interfaces

## Development Workflow

### Common Tasks

```bash
# Start development
make dev              # Run with hot reload

# Testing
make test             # Run all tests
make test-coverage    # Generate coverage report

# Database
make seed             # Add test users
make backup           # Backup database
make migrate          # Reset database

# Cleanup
make clean            # Remove build artifacts
make clean-all        # Complete cleanup
```

### Project Structure

```
zipcodereader/
├── main.go                 # Application entry point (76 lines)
├── routes/                 # Unified routing system
├── handlers/               # HTTP request handlers
│   ├── instructor.go       # All instructor endpoints
│   ├── student.go          # All student endpoints
│   ├── auth.go            # OAuth2 authentication
│   └── local_auth.go      # Local authentication
├── services/               # Business logic
│   ├── assignment.go       # Unified assignment service
│   └── auth.go            # Authentication service
├── models/                 # Data models
│   ├── user.go
│   ├── assignment.go
│   └── student_assignment.go
├── middleware/             # HTTP middleware
├── templates/              # HTML templates
├── static/                 # CSS, JS, images
├── docs/                   # Documentation (you are here)
└── scripts/                # Test and utility scripts
```

## Recent Major Changes

### October 2025 Refactoring

A comprehensive codebase simplification was completed:

- ✅ **50% reduction** in codebase size (4,400 → 2,243 lines)
- ✅ **Eliminated duplication** - Routes defined once, work everywhere
- ✅ **Handler consolidation** - 8 handlers → 2 main handlers (instructor, student)
- ✅ **Service consolidation** - 5 services → 1 unified assignment service
- ✅ **Enhanced testing** - 16 comprehensive unit tests
- ✅ **Data model improvements** - Added support for multiple assignment types

See [archive/Codebase-Simplification-Analysis.md](archive/Codebase-Simplification-Analysis.md) for full details.

## API Documentation

### Authentication Endpoints

#### Local Auth
- `POST /local/register` - Register new user
- `POST /local/login` - Login with username/password
- `GET /local/logout` - Logout

#### OAuth2
- `GET /auth/github` - Initiate GitHub OAuth2 flow
- `GET /auth/github/callback` - OAuth2 callback
- `GET /logout` - Logout

### Instructor Endpoints

- `GET /instructor/dashboard` - Dashboard view
- `GET /instructor/assignments` - List all assignments
- `POST /instructor/assignments` - Create assignment
- `GET /instructor/assignments/:id` - Get specific assignment
- `PUT /instructor/assignments/:id` - Update assignment
- `DELETE /instructor/assignments/:id` - Delete assignment
- `POST /instructor/assignments/:id/assign` - Assign to students
- `GET /instructor/assignments/:id/progress` - View progress
- `GET /instructor/students` - List students
- `GET /instructor/dashboard/stats` - Get statistics

### Student Endpoints

- `GET /student/dashboard` - Dashboard view
- `GET /student/assignments` - List assigned work
- `GET /student/assignments/:id` - View specific assignment
- `POST /student/assignments/:id/status` - Update status
- `POST /student/assignments/:id/complete` - Mark complete
- `POST /student/assignments/:id/progress` - Mark in progress
- `GET /student/dashboard/stats` - Get statistics
- `GET /student/categories` - Get assignment categories

## Contributing

When contributing to ZipCodeReader:

1. Read the [Development Log](../CLAUDE.md) to understand project history
2. Follow the existing code structure and patterns
3. Run `make check` before committing (runs fmt + tests)
4. Update documentation for significant changes
5. Add tests for new features

## Need Help?

- Check the [Quick Start Guide](quickstart.md) for basic setup
- Review the [Makefile Commands](makefile.md) for available commands
- Read the [Development Log](../CLAUDE.md) for technical context
- Check archived documents for historical context

## License

MIT License - See [LICENSE](../LICENSE) file for details.