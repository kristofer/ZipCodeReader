# ZipCodeReader Architecture

## Overview

ZipCodeReader is a web-based assignment management system built with Go, designed for educational environments. The application follows a clean, layered architecture with clear separation of concerns.

## Architecture Principles

### Design Goals

1. **Simplicity** - Minimize complexity and cognitive load
2. **Maintainability** - Easy to understand and modify
3. **Testability** - Comprehensive test coverage
4. **Scalability** - Can grow without complexity explosion

### Key Decisions

- **Consolidated handlers** - Role-based grouping (instructor.go, student.go)
- **Unified service layer** - Single assignment service handles all operations
- **Route consolidation** - Single source of truth for all routing
- **Minimal abstraction** - Only add complexity when needed

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        Client Browser                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Gin Web Framework                       │
│                                                               │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │   Routes    │─▶│  Middleware  │─▶│   Handlers   │       │
│  │  (routes/)  │  │(middleware/) │  │ (handlers/)  │       │
│  └─────────────┘  └──────────────┘  └──────────────┘       │
│                                              │               │
│                                              ▼               │
│                                      ┌──────────────┐        │
│                                      │   Services   │        │
│                                      │ (services/)  │        │
│                                      └──────────────┘        │
│                                              │               │
│                                              ▼               │
│                                      ┌──────────────┐        │
│                                      │    Models    │        │
│                                      │  (models/)   │        │
│                                      └──────────────┘        │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   SQLite Database (GORM)                     │
└─────────────────────────────────────────────────────────────┘
```

## Layer Responsibilities

### 1. Routes Layer (`routes/`)

**Purpose**: Single source of truth for all application routes

**Key File**: `routes/routes.go`

**Responsibilities**:
- Define all HTTP routes once
- Register middleware
- Group routes by authentication mode (local vs OAuth2)
- Map routes to handlers

**Benefits**:
- No route duplication
- Easy to see all endpoints at a glance
- Works for both authentication modes

### 2. Middleware Layer (`middleware/`)

**Purpose**: Cross-cutting concerns applied to routes

**Key Files**:
- `auth.go` - Authentication and authorization

**Responsibilities**:
- Session management
- Authentication checks
- Role-based access control (RBAC)
- Request validation

**Key Functions**:
- `RequireAuth()` - Ensure user is authenticated
- `RequireAuthWithUser()` - Inject user context into request
- `RequireInstructor()` - Ensure user has instructor role
- `RequireStudent()` - Ensure user has student role

### 3. Handlers Layer (`handlers/`)

**Purpose**: HTTP request handling and response generation

**Key Files**:
- `instructor.go` (~700 lines) - All instructor functionality
- `student.go` (~400 lines) - All student functionality
- `auth.go` (124 lines) - GitHub OAuth2 handlers
- `local_auth.go` (138 lines) - Local authentication handlers

**Responsibilities**:
- Parse HTTP requests
- Validate input
- Call service layer
- Generate HTTP responses (JSON or HTML)
- Error handling

**Design Pattern**: Handler functions follow this pattern:
1. Parse request parameters
2. Validate authorization
3. Call service layer
4. Handle errors
5. Return response

**Example**:
```go
func (h *InstructorHandlers) CreateAssignment(c *gin.Context) {
    // 1. Parse request
    var req CreateAssignmentRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 2. Get user from context
    user := c.MustGet("user").(*models.User)
    
    // 3. Call service
    assignment, err := h.service.CreateAssignment(user.ID, req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // 4. Return response
    c.JSON(201, assignment)
}
```

### 4. Services Layer (`services/`)

**Purpose**: Business logic and data orchestration

**Key Files**:
- `assignment.go` (~700 lines) - Unified assignment management service
- `auth.go` (105 lines) - Authentication logic

**Responsibilities**:
- Implement business rules
- Coordinate multiple models
- Transaction management
- Data validation
- Complex queries

**AssignmentService Methods** (56 total, organized by category):

*Assignment CRUD (17 methods)*:
- CreateAssignment, GetAssignment, UpdateAssignment, DeleteAssignment
- ListInstructorAssignments, ListAllAssignments
- GetAssignmentWithDetails, etc.

*Assignment-Student Relationships (9 methods)*:
- AssignToStudent, AssignToMultipleStudents
- RemoveStudentFromAssignment
- ListAssignedStudents, etc.

*Student Operations (18 methods)*:
- GetStudentAssignments, GetStudentAssignment
- UpdateAssignmentStatus, MarkAsCompleted, MarkAsInProgress
- GetOverdueAssignments, GetUpcomingAssignments, etc.

*Progress Tracking (6 methods)*:
- GetAssignmentProgress, GetDetailedProgress
- GetProgressTrends, GetCompletionAnalytics, etc.

*Due Date Notifications (6 methods)*:
- GetUpcomingDueDates, GetDueDateOverview
- GetDueDateAlerts, NotifyDueDates, etc.

### 5. Models Layer (`models/`)

**Purpose**: Data structures and database operations

**Key Files**:
- `user.go` - User model and authentication methods
- `assignment.go` - Assignment model and queries
- `student_assignment.go` - Assignment-student relationship

**Responsibilities**:
- Define data structures
- Database table mapping (GORM)
- Basic CRUD operations
- Simple queries
- Data validation

**Models**:

```go
// User - Represents instructors and students
type User struct {
    ID           uint
    Username     string
    Email        string
    PasswordHash string
    GithubID     int64
    Role         string  // "instructor" or "student"
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

// Assignment - Reading or programming assignment
type Assignment struct {
    ID               uint
    Title            string
    Description      string
    URL              string
    Category         string
    Type             string      // "reading", "programming", "quiz"
    EstimatedMinutes int
    RepositoryURL    string
    DueDate          *time.Time
    CreatedByID      uint
    CreatedBy        User
    CreatedAt        time.Time
    UpdatedAt        time.Time
    DeletedAt        gorm.DeletedAt
}

// StudentAssignment - Tracks student progress on assignments
type StudentAssignment struct {
    ID              uint
    AssignmentID    uint
    Assignment      Assignment
    StudentID       uint
    Student         User
    Status          string      // "assigned", "in_progress", "completed"
    TimeSpent       int         // Minutes spent
    ProgressPercent int         // 0-100
    SubmissionURL   string      // GitHub PR, Google Doc, etc.
    CompletedAt     *time.Time
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

## Authentication System

### Dual Authentication Mode

ZipCodeReader supports two authentication modes that can be switched via command-line flag:

1. **Local Authentication** (default) - Username/password with bcrypt
2. **GitHub OAuth2** - Login via GitHub account

### Session Management

- Sessions stored server-side using Gin sessions
- Session cookie contains encrypted session ID
- User data loaded from database on each request
- Sessions persist across server restarts (stored in database)

### Role-Based Access Control (RBAC)

Two roles:
- **instructor** - Can create/manage assignments, view all students
- **student** - Can view assigned work, update progress

Middleware enforces role requirements:
```go
// Only instructors can access
instructorRoutes.Use(middleware.RequireInstructor())

// Only students can access
studentRoutes.Use(middleware.RequireStudent())
```

## Data Flow Examples

### Example 1: Instructor Creates Assignment

```
1. Browser → POST /instructor/assignments
2. Routes → Middleware (auth check)
3. Middleware → InstructorHandlers.CreateAssignment
4. Handler → Validate request data
5. Handler → AssignmentService.CreateAssignment
6. Service → Validate business rules
7. Service → Assignment.Create (model)
8. Model → Database INSERT
9. Database → Return new assignment
10. Service → Return assignment to handler
11. Handler → JSON response to browser
```

### Example 2: Student Marks Assignment Complete

```
1. Browser → POST /student/assignments/123/complete
2. Routes → Middleware (auth + student check)
3. Middleware → StudentHandlers.MarkCompleted
4. Handler → Get user from context
5. Handler → AssignmentService.MarkAsCompleted
6. Service → Verify assignment exists
7. Service → Verify assignment is assigned to student
8. Service → Update StudentAssignment status
9. Model → Database UPDATE
10. Database → Return updated record
11. Service → Return to handler
12. Handler → JSON response to browser
```

## Database Schema

### Tables

**users**
- Primary key: `id`
- Unique: `username`, `github_id`
- Indexes: `username`, `github_id`, `role`

**assignments**
- Primary key: `id`
- Foreign key: `created_by_id` → `users(id)`
- Indexes: `created_by_id`, `deleted_at`, `category`, `type`
- Soft deletes supported

**student_assignments**
- Primary key: `id`
- Foreign keys: `assignment_id` → `assignments(id)`, `student_id` → `users(id)`
- Unique: `(assignment_id, student_id)`
- Indexes: `assignment_id`, `student_id`, `status`

### Relationships

```
User (instructor) ──< Assignment
                     (one-to-many: instructor creates many assignments)

Assignment >──< User (student)
            (many-to-many through student_assignments)

StudentAssignment
  ├── belongs to Assignment
  └── belongs to User (student)
```

## File Organization

```
zipcodereader/
├── main.go                     # Entry point (76 lines)
│   ├── Parse command-line flags
│   ├── Initialize database
│   ├── Setup Gin router
│   ├── Register routes
│   └── Start server
│
├── routes/
│   └── routes.go              # All route definitions (180 lines)
│
├── handlers/
│   ├── instructor.go          # Instructor endpoints (700 lines)
│   ├── student.go             # Student endpoints (400 lines)
│   ├── auth.go               # OAuth2 handlers (124 lines)
│   └── local_auth.go         # Local auth handlers (138 lines)
│
├── services/
│   ├── assignment.go          # Assignment service (700 lines)
│   └── auth.go               # Auth service (105 lines)
│
├── models/
│   ├── user.go               # User model (169 lines)
│   ├── assignment.go         # Assignment model (110 lines)
│   └── student_assignment.go # StudentAssignment model (200 lines)
│
├── middleware/
│   └── auth.go               # Auth middleware (200 lines)
│
├── config/
│   └── config.go             # Configuration (100 lines)
│
├── database/
│   └── migrations.go         # Database setup (108 lines)
│
├── templates/                 # HTML templates
│   ├── base.html
│   ├── instructor_assignments.html
│   ├── student_assignments.html
│   └── ...
│
└── static/                    # CSS, JS, images
    ├── css/
    ├── js/
    └── images/
```

## Testing Strategy

### Unit Tests

**Handler Tests** (`handlers/*_test.go`):
- Test HTTP endpoints with mock requests
- Verify response codes and data
- Test authentication and authorization
- 16 comprehensive tests

**Service Tests** (`services/*_test.go`):
- Test business logic
- Verify data transformations
- Test error handling

**Model Tests** (`models/*_test.go`):
- Test database operations
- Verify relationships
- Test validation rules

### Test Database

- Uses SQLite in-memory database for tests
- Fresh database for each test
- No test data pollution

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./handlers -v

# With coverage
make test-coverage
```

## Configuration

### Environment Variables

- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - Database file path (default: ./zipcodereader.db)
- `SESSION_SECRET` - Session encryption key
- `GITHUB_CLIENT_ID` - GitHub OAuth2 client ID
- `GITHUB_CLIENT_SECRET` - GitHub OAuth2 client secret
- `GITHUB_CALLBACK_URL` - OAuth2 callback URL

### Command-Line Flags

```bash
# Local authentication (default)
./zipcodereader

# GitHub OAuth2
./zipcodereader --use_oauth2

# Custom port
./zipcodereader --port=3000
```

## Performance Considerations

### Database

- SQLite3 with WAL mode for better concurrency
- Proper indexes on foreign keys and search fields
- Soft deletes to preserve data integrity
- Connection pooling via GORM

### Caching

- Session data cached in server memory
- User data loaded once per request
- Static files served with caching headers

### Scalability

Current architecture supports:
- 100s of concurrent users
- 1000s of assignments
- 10,000s of student assignments

For larger scale, consider:
- PostgreSQL instead of SQLite
- Redis for session storage
- CDN for static assets
- Load balancing multiple instances

## Security

### Authentication

- bcrypt password hashing (cost factor 10)
- Secure session cookies (HTTP-only, SameSite)
- GitHub OAuth2 with state validation
- Session timeout (configurable)

### Authorization

- Role-based access control (RBAC)
- Ownership verification for resources
- Instructor can only modify their own assignments
- Students can only access assigned work

### Input Validation

- Request body validation with struct tags
- URL validation for assignment links
- SQL injection prevention (GORM parameterized queries)
- XSS prevention (HTML escaping in templates)

## Error Handling

### HTTP Status Codes

- `200 OK` - Successful GET
- `201 Created` - Successful POST
- `400 Bad Request` - Invalid input
- `401 Unauthorized` - Not authenticated
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

### Error Responses

JSON format:
```json
{
  "error": "Assignment not found"
}
```

HTML format:
- Error pages with user-friendly messages
- Redirect to login for auth errors
- Flash messages for form errors

## Future Architecture Improvements

### When to Consider Changes

**Add microservices** when:
- Team size > 10 developers
- Need independent scaling of components
- Different deployment schedules for features

**Add API gateway** when:
- Multiple client applications
- Need rate limiting per client
- Complex authentication flows

**Add caching layer** when:
- Database queries become bottleneck
- Read:write ratio > 10:1
- Need sub-100ms response times

**Add message queue** when:
- Long-running background tasks
- Need asynchronous processing
- Event-driven architecture needed

### Principles to Maintain

1. **Keep it simple** - Only add complexity when justified
2. **Consolidate related code** - Related functionality stays together
3. **Single source of truth** - No duplication
4. **Test everything** - Comprehensive test coverage
5. **Document decisions** - Explain why, not just what

## Resources

- [Go Gin Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [Go Testing Best Practices](https://golang.org/doc/code.html#Testing)
- [12 Factor App](https://12factor.net/)
- [REST API Design](https://restfulapi.net/)

---

**Last Updated**: October 2025  
**Version**: 1.0 (Post-Refactoring)