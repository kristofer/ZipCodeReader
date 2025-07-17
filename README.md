# ZipCodeReader

A web-based reading list manager for ZipCode students and instructors.

> A copilot-assisted app

## Overview

ZipCodeReader is a comprehensive reading assignment management system designed for educational environments. It enables instructors to assign reading materials to students and allows students to track their progress through a user-friendly web interface.

## Features

### Core Functionality

#### For Instructors

- **Assignment Creation**: Create reading assignments using a simple web-based form
- **URL-based Assignments**: Initially support web page URLs as reading materials
- **Student Management**: Assign readings to individual students or groups
- **Progress Tracking**: Monitor student completion rates and reading progress
- **Assignment Organization**: Categorize and organize assignments by course, topic, or difficulty

#### For Students

- **GitHub OAuth2 Authentication**: Secure login using GitHub credentials
- **Personal Dashboard**: View all assigned readings in one place
- **Progress Tracking**: Check off completed assignments and track reading progress
- **Assignment Details**: Access reading materials and view assignment requirements
- **Reading History**: View completed assignments and reading statistics

### Technical Requirements

#### Authentication & Authorization

- GitHub OAuth2 integration for secure user authentication
- Role-based access control (Student vs Instructor)
- Session management and secure logout

#### Data Management

- User profiles linked to GitHub accounts
- Assignment tracking and completion status
- Reading progress persistence
- Assignment metadata (due dates, categories, etc.)

#### User Interface

- Responsive web design for desktop and mobile access
- Intuitive forms for assignment creation
- Clean, accessible student dashboard
- Real-time progress updates

## User Stories

### Instructor Stories

- As an instructor, I want to create reading assignments so that I can assign materials to my students
- As an instructor, I want to assign readings to specific students or groups so that I can customize learning paths
- As an instructor, I want to track student progress so that I can identify who needs additional support
- As an instructor, I want to organize assignments by category so that I can manage my curriculum effectively

### Student Stories

- As a student, I want to log in with my GitHub account so that I can securely access my assignments
- As a student, I want to see all my assigned readings in one place so that I can manage my workload
- As a student, I want to mark assignments as complete so that I can track my progress
- As a student, I want to access reading materials easily so that I can complete my assignments efficiently

## Technical Architecture

### Frontend

- Modern web framework (React, Vue, or similar)
- Responsive CSS framework
- GitHub OAuth2 client integration
- Real-time updates for assignment status

### Backend

- RESTful API or GraphQL endpoint
- GitHub OAuth2 server integration
- Database for user and assignment management
- Session management and security

### Database Schema (Preliminary)

- **Users**: GitHub ID, role (student/instructor), profile information
- **Assignments**: Title, URL, description, due date, category, creator
- **Student_Assignments**: Assignment ID, student ID, completion status, completion date
- **Groups**: Group management for bulk assignment features

## Future Features

### Enhanced Assignment Types

- PDF document support
- Video content integration
- Interactive reading materials
- Multi-media assignments

### Advanced Tracking

- Time spent reading analytics
- Reading comprehension quizzes
- Progress reports and analytics
- Deadline reminders and notifications

### Collaboration Features

- Discussion boards for assignments
- Peer review capabilities
- Group reading assignments
- Instructor feedback system

### Integration Capabilities

- LMS integration (Canvas, Blackboard, etc.)
- Calendar integration
- Email notifications
- Mobile app companion

### Administrative Features

- Bulk assignment creation
- Template management
- Advanced reporting and analytics
- Course management tools

## Installation & Setup

Coming soon - detailed setup instructions will be provided once development begins.

## Contributing

This project is developed with GitHub Copilot assistance. Contributions are welcome following standard GitHub workflow practices.

## License

License information to be determined.

## Implementation Plan

### Technology Stack

- **Backend**: Go with Gin web framework
- **Database**: SQLite3 with GORM ORM
- **Authentication**: GitHub OAuth2 (using go-github and oauth2 packages)
- **Frontend**: HTML templates with Tailwind CSS (initially), later React/Vue integration
- **Session Management**: Gin sessions with Redis/memory store

### Development Phases

#### Phase 1: Project Foundation (Week 1)

**Goal**: Set up basic project structure and dependencies

**Tasks**:

1. Initialize Go module and project structure
2. Set up Gin web server with basic routing
3. Configure SQLite3 database with GORM
4. Create basic HTML templates and static file serving
5. Set up environment configuration and logging
6. Create basic health check endpoint

**Deliverables**:

- Working web server on localhost
- Database connection established
- Basic project structure with proper Go conventions

**Files to Create**:

```text
├── main.go
├── go.mod
├── go.sum
├── config/
│   └── config.go
├── models/
│   └── models.go
├── handlers/
│   └── handlers.go
├── middleware/
│   └── auth.go
├── templates/
│   ├── base.html
│   └── index.html
├── static/
│   ├── css/
│   └── js/
└── database/
    └── migrations.go
```

#### Phase 2: Authentication System (Week 2)

**Goal**: Implement GitHub OAuth2 authentication

**Tasks**:

1. Set up GitHub OAuth2 application
2. Implement OAuth2 flow with GitHub
3. Create user registration and login handlers
4. Implement session management
5. Create user model and database schema
6. Add role-based access control (student/instructor)
7. Create protected routes middleware

**Deliverables**:

- Users can log in with GitHub
- Session persistence
- Role-based access control
- User dashboard skeleton

**Database Schema**:

```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    github_id INTEGER UNIQUE NOT NULL,
    username TEXT NOT NULL,
    email TEXT,
    avatar_url TEXT,
    role TEXT DEFAULT 'student',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Phase 3: Basic Assignment Management (Week 3)

**Goal**: Core assignment creation and viewing functionality

**Tasks**:

1. Create assignment model and database schema
2. Implement assignment creation form for instructors
3. Create assignment listing page
4. Implement basic assignment viewing
5. Add assignment validation and error handling
6. Create instructor dashboard for assignment management

**Deliverables**:

- Instructors can create URL-based assignments
- Students can view assigned readings
- Basic assignment management interface

**Database Schema**:

```sql
CREATE TABLE assignments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    url TEXT NOT NULL,
    due_date DATETIME,
    category TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE student_assignments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    assignment_id INTEGER REFERENCES assignments(id),
    student_id INTEGER REFERENCES users(id),
    completed BOOLEAN DEFAULT FALSE,
    completed_at DATETIME,
    assigned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(assignment_id, student_id)
);
```

#### Phase 4: Student Progress Tracking (Week 4)

**Goal**: Enable students to track and complete assignments

**Tasks**:

1. Create student dashboard with assigned readings
2. Implement assignment completion functionality
3. Add progress tracking and statistics
4. Create assignment detail pages
5. Implement reading history view
6. Add basic progress indicators

**Deliverables**:

- Students can mark assignments as complete
- Progress tracking and completion statistics
- Reading history functionality
- Responsive student dashboard

#### Phase 5: Assignment Distribution (Week 5)

**Goal**: Advanced assignment management and distribution

**Tasks**:

1. Implement group management system
2. Add bulk assignment capabilities
3. Create assignment distribution to specific students/groups
4. Implement assignment categories and organization
5. Add assignment search and filtering
6. Create assignment editing and deletion

**Deliverables**:

- Group-based assignment distribution
- Advanced assignment management
- Search and filtering capabilities
- Assignment organization tools

**Database Schema**:

```sql
CREATE TABLE groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    created_by INTEGER REFERENCES users(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE group_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER REFERENCES groups(id),
    user_id INTEGER REFERENCES users(id),
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(group_id, user_id)
);
```

#### Phase 6: Enhanced UI/UX (Week 6)

**Goal**: Improve user interface and experience

**Tasks**:

1. Implement responsive design with Tailwind CSS
2. Add real-time updates using WebSockets or Server-Sent Events
3. Improve form validation and user feedback
4. Add loading states and error handling
5. Implement dark/light mode toggle
6. Add accessibility improvements

**Deliverables**:

- Modern, responsive UI
- Real-time progress updates
- Improved user experience
- Accessibility compliance

#### Phase 7: Advanced Features (Week 7)

**Goal**: Add advanced tracking and reporting features

**Tasks**:

1. Implement assignment analytics and reporting
2. Add deadline reminders and notifications
3. Create progress reports for instructors
4. Implement assignment templates
5. Add export functionality for reports
6. Create advanced search and filtering

**Deliverables**:

- Comprehensive analytics dashboard
- Automated notifications
- Detailed progress reports
- Template management system

#### Phase 8: API Development (Week 8)

**Goal**: Create RESTful API for future integrations

**Tasks**:

1. Design and implement RESTful API endpoints
2. Add API authentication and rate limiting
3. Create API documentation
4. Implement API versioning
5. Add comprehensive error handling
6. Create API client examples

**Deliverables**:

- Full REST API
- API documentation
- Authentication system for API
- Client integration examples

#### Phase 9: Testing and Quality Assurance (Week 9)

**Goal**: Comprehensive testing and bug fixes

**Tasks**:

1. Write unit tests for all handlers and models
2. Implement integration tests
3. Add end-to-end testing
4. Performance testing and optimization
5. Security audit and improvements
6. Bug fixes and code refactoring

**Deliverables**:

- Comprehensive test suite
- Performance optimizations
- Security hardening
- Code quality improvements

#### Phase 10: Production Deployment (Week 10)

**Goal**: Deploy to production environment

**Tasks**:

1. Set up production environment configuration
2. Implement database migrations
3. Configure SSL/TLS and security headers
4. Set up monitoring and logging
5. Create backup and recovery procedures
6. Write deployment documentation

**Deliverables**:

- Production-ready application
- Deployment documentation
- Monitoring and alerting
- Backup procedures

### Future Enhancements (Post-Launch)

#### Phase 11: Mobile Optimization

- Progressive Web App (PWA) implementation
- Mobile-specific UI improvements
- Offline functionality

#### Phase 12: Advanced Content Types

- PDF document support
- Video content integration
- Interactive reading materials

#### Phase 13: Collaboration Features

- Discussion boards
- Peer review system
- Comment system on assignments

#### Phase 14: External Integrations

- LMS integration (Canvas, Blackboard)
- Calendar integration
- Email notifications
- Slack/Discord integration

### Development Commands

```bash
# Initialize project
go mod init zipcodereader

# Key dependencies
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/google/go-github/v45/github
go get golang.org/x/oauth2
go get github.com/gin-contrib/sessions
go get github.com/gin-contrib/cors

# Run development server
go run main.go

# Build for production
go build -o zipcodereader main.go

# Run tests
go test ./...
```

### Project Structure (Final)

```text
zipcodereader/
├── main.go
├── go.mod
├── go.sum
├── config/
│   ├── config.go
│   └── database.go
├── models/
│   ├── user.go
│   ├── assignment.go
│   ├── group.go
│   └── progress.go
├── handlers/
│   ├── auth.go
│   ├── assignments.go
│   ├── dashboard.go
│   ├── api.go
│   └── admin.go
├── middleware/
│   ├── auth.go
│   ├── cors.go
│   └── logger.go
├── services/
│   ├── github.go
│   ├── assignment.go
│   └── notification.go
├── templates/
│   ├── base.html
│   ├── dashboard.html
│   ├── assignments.html
│   └── admin.html
├── static/
│   ├── css/
│   ├── js/
│   └── images/
├── tests/
│   ├── handlers_test.go
│   ├── models_test.go
│   └── integration_test.go
├── docs/
│   ├── api.md
│   └── deployment.md
└── scripts/
    ├── migrate.go
    └── seed.go
```

