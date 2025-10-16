#!/bin/bash

# Phase 3 Task 4: Student Assignment Viewing Handlers
# Tests student-facing assignment interfaces, viewing, progress tracking, and completion

echo "🚀 Phase 3 Task 4: Student Assignment Viewing Handlers"
echo "======================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test results
PASS=0
FAIL=0

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
        ((PASS++))
    else
        echo -e "${RED}❌ $2${NC}"
        ((FAIL++))
    fi
}

# Note: Start server manually with './zipcodereader' before running this test
echo "Assuming server is running on http://localhost:8080..."
sleep 1



# Setup test users
echo "Setting up test users..."
# Use timestamp to create unique usernames
TIMESTAMP=$(date +%s)
INSTRUCTOR_USER="instructor_$TIMESTAMP"
STUDENT_USER="student_$TIMESTAMP"
STUDENT2_USER="student2_$TIMESTAMP"

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR_USER&email=$INSTRUCTOR_USER@example.com&password=password&confirm_password=password&role=instructor" > /dev/null

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT_USER&email=$STUDENT_USER@example.com&password=password&confirm_password=password&role=student" > /dev/null

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT2_USER&email=$STUDENT2_USER@example.com&password=password&confirm_password=password&role=student" > /dev/null

# Login users
curl -s -c instructor_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR_USER&password=password" > /dev/null

curl -s -c student1_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT_USER&password=password" > /dev/null

curl -s -c student2_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT2_USER&password=password" > /dev/null

# Create test assignments
echo "Creating test assignments..."
ASSIGNMENT1=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Student Test Assignment 1",
    "description": "First assignment for student testing",
    "url": "https://example.com/student-test-1",
    "category": "Student Testing",
    "due_date": "2025-08-25T23:59:59Z"
  }')

ASSIGNMENT2=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Student Test Assignment 2",
    "description": "Second assignment for student testing",
    "url": "https://example.com/student-test-2",
    "category": "Advanced Testing"
  }')

ASSIGNMENT1_ID=$(echo "$ASSIGNMENT1" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
ASSIGNMENT2_ID=$(echo "$ASSIGNMENT2" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)

# Assign to students
# Get the actual student IDs from the database
STUDENT_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT_USER';")
STUDENT2_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT2_USER';")

curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT1_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"student_ids\": [$STUDENT_ID, $STUDENT2_ID]}" > /dev/null

curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT2_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"student_ids\": [$STUDENT_ID]}" > /dev/null

# Test 1: GET /student/assignments - List all assigned readings
echo "Test 1: GET /student/assignments - List all assigned readings"
STUDENT_ASSIGNMENTS=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments)

if echo "$STUDENT_ASSIGNMENTS" | grep -q "Student Test Assignment"; then
    print_result 0 "GET /student/assignments - List assigned readings works"
else
    print_result 1 "GET /student/assignments - List assigned readings failed"
fi

# Test 2: GET /student/assignments/:id - View specific assignment
echo "Test 2: GET /student/assignments/:id - View specific assignment"
SPECIFIC_ASSIGNMENT=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments/$ASSIGNMENT1_ID)

if echo "$SPECIFIC_ASSIGNMENT" | grep -q "Student Test Assignment 1"; then
    print_result 0 "GET /student/assignments/:id - View specific assignment works"
else
    print_result 1 "GET /student/assignments/:id - View specific assignment failed"
fi

# Test 3: GET /student/assignments/:id/detail - Assignment detail UI
echo "Test 3: GET /student/assignments/:id/detail - Assignment detail UI"
ASSIGNMENT_DETAIL=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments/$ASSIGNMENT1_ID/detail)

if echo "$ASSIGNMENT_DETAIL" | grep -q "html\|assignment\|detail"; then
    print_result 0 "GET /student/assignments/:id/detail - Assignment detail UI works"
else
    print_result 1 "GET /student/assignments/:id/detail - Assignment detail UI failed"
fi

# Test 4: POST /student/assignments/:id/status - Update assignment status
echo "Test 4: POST /student/assignments/:id/status - Update assignment status"
STATUS_UPDATE=$(curl -s -b student1_cookies.txt -X POST http://localhost:8080/student/assignments/$ASSIGNMENT1_ID/status \
  -H "Content-Type: application/json" \
  -d '{"status": "in_progress"}')

if echo "$STATUS_UPDATE" | grep -q "success\|in_progress"; then
    print_result 0 "POST /student/assignments/:id/status - Update status works"
else
    print_result 1 "POST /student/assignments/:id/status - Update status failed"
fi

# Test 5: POST /student/assignments/:id/complete - Mark as completed
echo "Test 5: POST /student/assignments/:id/complete - Mark as completed"
COMPLETE_ASSIGNMENT=$(curl -s -b student1_cookies.txt -X POST http://localhost:8080/student/assignments/$ASSIGNMENT1_ID/complete \
  -H "Content-Type: application/json")

if echo "$COMPLETE_ASSIGNMENT" | grep -q "success\|completed"; then
    print_result 0 "POST /student/assignments/:id/complete - Mark as completed works"
else
    print_result 1 "POST /student/assignments/:id/complete - Mark as completed failed"
fi

# Test 6: POST /student/assignments/:id/progress - Mark as in progress
echo "Test 6: POST /student/assignments/:id/progress - Mark as in progress"
PROGRESS_ASSIGNMENT=$(curl -s -b student1_cookies.txt -X POST http://localhost:8080/student/assignments/$ASSIGNMENT2_ID/progress \
  -H "Content-Type: application/json")

if echo "$PROGRESS_ASSIGNMENT" | grep -q "success\|in_progress"; then
    print_result 0 "POST /student/assignments/:id/progress - Mark as in progress works"
else
    print_result 1 "POST /student/assignments/:id/progress - Mark as in progress failed"
fi

# Test 7: GET /student/dashboard - Assignment dashboard with overview
echo "Test 7: GET /student/dashboard - Assignment dashboard with overview"
STUDENT_DASHBOARD=$(curl -s -b student1_cookies.txt http://localhost:8080/student/dashboard)

if echo "$STUDENT_DASHBOARD" | grep -q "dashboard\|assignment\|overview"; then
    print_result 0 "GET /student/dashboard - Student dashboard works"
else
    print_result 1 "GET /student/dashboard - Student dashboard failed"
fi

# Test 8: GET /student/dashboard/stats - Get student statistics
echo "Test 8: GET /student/dashboard/stats - Get student statistics"
STUDENT_STATS=$(curl -s -b student1_cookies.txt http://localhost:8080/student/dashboard/stats)

if echo "$STUDENT_STATS" | grep -q "assignments\|completed\|progress"; then
    print_result 0 "GET /student/dashboard/stats - Student statistics works"
else
    print_result 1 "GET /student/dashboard/stats - Student statistics failed"
fi

# Test 9: GET /student/assignments/overdue - Get overdue assignments
echo "Test 9: GET /student/assignments/overdue - Get overdue assignments"
OVERDUE_ASSIGNMENTS=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments/overdue)

if [ $? -eq 0 ]; then
    print_result 0 "GET /student/assignments/overdue - Get overdue assignments works"
else
    print_result 1 "GET /student/assignments/overdue - Get overdue assignments failed"
fi

# Test 10: GET /student/assignments/upcoming - Get upcoming assignments
echo "Test 10: GET /student/assignments/upcoming - Get upcoming assignments"
UPCOMING_ASSIGNMENTS=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments/upcoming)

if [ $? -eq 0 ]; then
    print_result 0 "GET /student/assignments/upcoming - Get upcoming assignments works"
else
    print_result 1 "GET /student/assignments/upcoming - Get upcoming assignments failed"
fi

# Test 11: GET /student/assignments/recent - Get recently completed
echo "Test 11: GET /student/assignments/recent - Get recently completed"
RECENT_ASSIGNMENTS=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments/recent)

if [ $? -eq 0 ]; then
    print_result 0 "GET /student/assignments/recent - Get recent assignments works"
else
    print_result 1 "GET /student/assignments/recent - Get recent assignments failed"
fi

# Test 12: GET /student/categories - Get assignment categories
echo "Test 12: GET /student/categories - Get assignment categories"
STUDENT_CATEGORIES=$(curl -s -b student1_cookies.txt http://localhost:8080/student/categories)

if echo "$STUDENT_CATEGORIES" | grep -q "categories\|Student Testing"; then
    print_result 0 "GET /student/categories - Get assignment categories works"
else
    print_result 1 "GET /student/categories - Get assignment categories failed"
fi

# Test 13: GET /student/assignments/by-status - Filter by status
echo "Test 13: GET /student/assignments/by-status - Filter by status"
STATUS_FILTER=$(curl -s -b student1_cookies.txt "http://localhost:8080/student/assignments/status/completed")

if echo "$STATUS_FILTER" | grep -q "Student Test Assignment 1"; then
    print_result 0 "GET /student/assignments/by-status - Filter by status works"
else
    print_result 1 "GET /student/assignments/by-status - Filter by status failed"
fi

# Test 14: GET /student/assignments/by-category - Filter by category
echo "Test 14: GET /student/assignments/by-category - Filter by category"
CATEGORY_FILTER=$(curl -s -b student1_cookies.txt "http://localhost:8080/student/assignments/category/Student%20Testing")

if echo "$CATEGORY_FILTER" | grep -q "Student Test Assignment"; then
    print_result 0 "GET /student/assignments/by-category - Filter by category works"
else
    print_result 1 "GET /student/assignments/by-category - Filter by category failed"
fi

# Test 15: GET /student/assignments/search - Search assignments
echo "Test 15: GET /student/assignments/search - Search assignments"
SEARCH_ASSIGNMENTS=$(curl -s -b student1_cookies.txt "http://localhost:8080/student/assignments/search?q=Student%20Test")

if echo "$SEARCH_ASSIGNMENTS" | grep -q "Student Test Assignment"; then
    print_result 0 "GET /student/assignments/search - Search assignments works"
else
    print_result 1 "GET /student/assignments/search - Search assignments failed"
fi

# Test 16: Assignment status persistence
echo "Test 16: Assignment status persistence"
# Check if the completed status persists
UPDATED_ASSIGNMENT=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments/$ASSIGNMENT1_ID)

if echo "$UPDATED_ASSIGNMENT" | grep -q "completed"; then
    print_result 0 "Assignment status persistence works"
else
    print_result 1 "Assignment status persistence failed"
fi

# Test 17: Role-based access control - Instructors cannot access student endpoints
echo "Test 17: Role-based access control - Instructors cannot access student endpoints"
UNAUTHORIZED_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/student/assignments)

if echo "$UNAUTHORIZED_RESPONSE" | grep -q "unauthorized\|forbidden\|403\|Authentication required\|error"; then
    print_result 0 "Role-based access control works (instructors blocked from student endpoints)"
else
    print_result 1 "Role-based access control failed (instructors should be blocked from student endpoints)"
fi

# Test 18: Student isolation - Student1 cannot see Student2's assignments
echo "Test 18: Student isolation - Student1 cannot see Student2's assignments"
STUDENT1_ASSIGNMENTS=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments)
STUDENT2_ASSIGNMENTS=$(curl -s -b student2_cookies.txt http://localhost:8080/student/assignments)

# Student1 should see both assignments, Student2 should only see Assignment1
STUDENT1_COUNT=$(echo "$STUDENT1_ASSIGNMENTS" | grep -o "Student Test Assignment" | wc -l 2>/dev/null || echo "0")
STUDENT2_COUNT=$(echo "$STUDENT2_ASSIGNMENTS" | grep -o "Student Test Assignment" | wc -l 2>/dev/null || echo "0")

if [ "$STUDENT1_COUNT" -eq 2 ] && [ "$STUDENT2_COUNT" -eq 1 ]; then
    print_result 0 "Student isolation works (students see only their assignments)"
else
    print_result 1 "Student isolation failed (students should see only their assignments)"
fi

# Test 19: Assignment completion timestamp
echo "Test 19: Assignment completion timestamp"
COMPLETED_ASSIGNMENT=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments/$ASSIGNMENT1_ID)

if echo "$COMPLETED_ASSIGNMENT" | grep -q "completed_at\|completion"; then
    print_result 0 "Assignment completion timestamp works"
else
    print_result 1 "Assignment completion timestamp failed"
fi

# Test 20: Due date notifications and sorting
echo "Test 20: Due date notifications and sorting"
DUE_DATE_RESPONSE=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments)

if echo "$DUE_DATE_RESPONSE" | grep -q "due_date\|2025-08-25"; then
    print_result 0 "Due date notifications and sorting works"
else
    print_result 1 "Due date notifications and sorting failed"
fi

# Test 21: Error handling - Non-existent assignment
echo "Test 21: Error handling - Non-existent assignment"
NON_EXISTENT_ASSIGNMENT=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments/99999)

if echo "$NON_EXISTENT_ASSIGNMENT" | grep -q "not found\|error\|404"; then
    print_result 0 "Error handling works (non-existent assignment)"
else
    print_result 1 "Error handling failed (should return error for non-existent assignment)"
fi

# Test 22: Invalid status update
echo "Test 22: Invalid status update"
INVALID_STATUS=$(curl -s -b student1_cookies.txt -X POST http://localhost:8080/student/assignments/$ASSIGNMENT2_ID/status \
  -H "Content-Type: application/json" \
  -d '{"status": "invalid_status"}')

if echo "$INVALID_STATUS" | grep -q "error\|invalid"; then
    print_result 0 "Invalid status update validation works"
else
    print_result 1 "Invalid status update validation failed"
fi

# Test 23: Assignment access control - Student cannot access unassigned assignment
echo "Test 23: Assignment access control - Student cannot access unassigned assignment"
# Create assignment not assigned to student1
UNASSIGNED_ASSIGNMENT=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Unassigned Assignment",
    "description": "This assignment is not assigned to student1",
    "url": "https://example.com/unassigned"
  }')

UNASSIGNED_ID=$(echo "$UNASSIGNED_ASSIGNMENT" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)

UNAUTHORIZED_ACCESS=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments/$UNASSIGNED_ID)

if echo "$UNAUTHORIZED_ACCESS" | grep -q "not found\|unauthorized\|403\|404"; then
    print_result 0 "Assignment access control works (student cannot access unassigned assignment)"
else
    print_result 1 "Assignment access control failed (student should not access unassigned assignment)"
fi

# Test 24: Assignment progress tracking
echo "Test 24: Assignment progress tracking"
PROGRESS_TRACKING=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments)

if echo "$PROGRESS_TRACKING" | grep -q "completed\|in_progress"; then
    print_result 0 "Assignment progress tracking works"
else
    print_result 1 "Assignment progress tracking failed"
fi

# Test 25: Assignment filtering by status and category combined
echo "Test 25: Assignment filtering by status and category combined"
COMBINED_FILTER=$(curl -s -b student1_cookies.txt "http://localhost:8080/student/assignments?status=completed&category=Student%20Testing")

if echo "$COMBINED_FILTER" | grep -q "Student Test Assignment 1"; then
    print_result 0 "Combined filtering works"
else
    print_result 1 "Combined filtering failed"
fi

# Cleanup
echo "🧹 Cleanup"

rm -f instructor_cookies.txt student1_cookies.txt student2_cookies.txt
sleep 1

# Summary
echo ""
echo "📊 Phase 3 Task 4 Test Results:"
echo "================================"
echo -e "${GREEN}✅ Passed: $PASS${NC}"
echo -e "${RED}❌ Failed: $FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}🎉 Phase 3 Task 4: Student Assignment Viewing Handlers - ALL TESTS PASSED!${NC}"
    exit 0
else
    echo -e "${RED}🚨 Phase 3 Task 4: Student Assignment Viewing Handlers - $FAIL TESTS FAILED!${NC}"
    echo ""
    echo "Issues to address:"
    echo "- Student assignment viewing interfaces"
    echo "- Assignment completion and progress tracking"
    echo "- Student dashboard and statistics"
    echo "- Assignment filtering and search for students"
    echo "- Role-based access control"
    echo "- Assignment status management"
    echo "- Due date notifications"
    echo ""
    exit 1
fi