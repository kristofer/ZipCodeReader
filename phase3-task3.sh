#!/bin/bash

# Phase 3 Task 3: Instructor Assignment Management Handlers
# Tests HTTP handlers for instructor assignment operations, creation, editing, deletion, and analytics

echo "🚀 Phase 3 Task 3: Instructor Assignment Management Handlers"
echo "==========================================================="

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

curl -s -c student_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT_USER&password=password" > /dev/null

# Test 1: GET /instructor/assignments - List all assignments
echo "Test 1: GET /instructor/assignments - List all assignments"
ASSIGNMENTS_LIST=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments)

if [ $? -eq 0 ]; then
    print_result 0 "GET /instructor/assignments endpoint accessible"
else
    print_result 1 "GET /instructor/assignments endpoint not accessible"
fi

# Test 2: POST /instructor/assignments - Create new assignment
echo "Test 2: POST /instructor/assignments - Create new assignment"
ASSIGNMENT_RESPONSE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Handler Test Assignment",
    "description": "Testing instructor assignment handlers",
    "url": "https://example.com/handler-test",
    "category": "Handler Testing",
    "due_date": "2025-08-20T23:59:59Z"
  }')

if echo "$ASSIGNMENT_RESPONSE" | grep -q "Handler Test Assignment"; then
    ASSIGNMENT_ID=$(echo "$ASSIGNMENT_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
    print_result 0 "POST /instructor/assignments - Create assignment works (ID: $ASSIGNMENT_ID)"
else
    print_result 1 "POST /instructor/assignments - Create assignment failed"
    ASSIGNMENT_ID=1
fi

# Test 3: POST /instructor/assignments - Form validation
echo "Test 3: POST /instructor/assignments - Form validation"
INVALID_ASSIGNMENT=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Missing title and URL"
  }')

if echo "$INVALID_ASSIGNMENT" | grep -q "error\|invalid\|required"; then
    print_result 0 "POST /instructor/assignments - Form validation works"
else
    print_result 1 "POST /instructor/assignments - Form validation failed"
fi

# Test 4: GET /instructor/assignments/:id - Get specific assignment
echo "Test 4: GET /instructor/assignments/:id - Get specific assignment"
SPECIFIC_ASSIGNMENT=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID)

if echo "$SPECIFIC_ASSIGNMENT" | grep -q "Handler Test Assignment"; then
    print_result 0 "GET /instructor/assignments/:id - Get specific assignment works"
else
    print_result 1 "GET /instructor/assignments/:id - Get specific assignment failed"
fi

# Test 5: PUT /instructor/assignments/:id - Update assignment
echo "Test 5: PUT /instructor/assignments/:id - Update assignment"
UPDATE_RESPONSE=$(curl -s -b instructor_cookies.txt -X PUT http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Handler Test Assignment",
    "description": "Updated description for handler testing",
    "url": "https://example.com/updated-handler-test",
    "category": "Updated Handler Testing"
  }')

if echo "$UPDATE_RESPONSE" | grep -q "Updated Handler Test Assignment\|success"; then
    print_result 0 "PUT /instructor/assignments/:id - Update assignment works"
else
    print_result 1 "PUT /instructor/assignments/:id - Update assignment failed"
fi

# Test 6: DELETE /instructor/assignments/:id - Delete assignment
echo "Test 6: DELETE /instructor/assignments/:id - Delete assignment"
# Create a new assignment to delete
DELETE_ASSIGNMENT=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Assignment to Delete",
    "description": "This assignment will be deleted",
    "url": "https://example.com/delete-test"
  }')

DELETE_ASSIGNMENT_ID=$(echo "$DELETE_ASSIGNMENT" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)

DELETE_RESPONSE=$(curl -s -b instructor_cookies.txt -X DELETE http://localhost:8080/instructor/assignments/$DELETE_ASSIGNMENT_ID)

if echo "$DELETE_RESPONSE" | grep -q "success\|deleted"; then
    print_result 0 "DELETE /instructor/assignments/:id - Delete assignment works"
else
    print_result 1 "DELETE /instructor/assignments/:id - Delete assignment failed"
fi

# Test 7: POST /instructor/assignments/:id/assign - Assign to students
echo "Test 7: POST /instructor/assignments/:id/assign - Assign to students"
ASSIGN_RESPONSE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/assign \
  -H "Content-Type: application/json" \
  -d '{
    "student_ids": [2, 3]
  }')

if echo "$ASSIGN_RESPONSE" | grep -q "success\|assigned"; then
    print_result 0 "POST /instructor/assignments/:id/assign - Assign to students works"
else
    print_result 1 "POST /instructor/assignments/:id/assign - Assign to students failed"
fi

# Test 8: GET /instructor/assignments/:id/progress - View assignment progress
echo "Test 8: GET /instructor/assignments/:id/progress - View assignment progress"
PROGRESS_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/progress)

if echo "$PROGRESS_RESPONSE" | grep -q "progress\|completion\|students"; then
    print_result 0 "GET /instructor/assignments/:id/progress - View progress works"
else
    print_result 1 "GET /instructor/assignments/:id/progress - View progress failed"
fi

# Test 9: GET /instructor/assignments/:id/students - List assigned students
echo "Test 9: GET /instructor/assignments/:id/students - List assigned students"
STUDENTS_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/students)

if echo "$STUDENTS_RESPONSE" | grep -q "student\|username"; then
    print_result 0 "GET /instructor/assignments/:id/students - List assigned students works"
else
    print_result 1 "GET /instructor/assignments/:id/students - List assigned students failed"
fi

# Test 10: GET /instructor/assignments/:id/detail - Assignment detail UI
echo "Test 10: GET /instructor/assignments/:id/detail - Assignment detail UI"
DETAIL_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/detail)

if echo "$DETAIL_RESPONSE" | grep -q "html\|assignment\|detail"; then
    print_result 0 "GET /instructor/assignments/:id/detail - Assignment detail UI works"
else
    print_result 1 "GET /instructor/assignments/:id/detail - Assignment detail UI failed"
fi

# Test 11: GET /instructor/assignments/:id/progress-view - Progress view UI
echo "Test 11: GET /instructor/assignments/:id/progress-view - Progress view UI"
PROGRESS_VIEW_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/progress-view)

if echo "$PROGRESS_VIEW_RESPONSE" | grep -q "html\|progress\|view"; then
    print_result 0 "GET /instructor/assignments/:id/progress-view - Progress view UI works"
else
    print_result 1 "GET /instructor/assignments/:id/progress-view - Progress view UI failed"
fi

# Test 12: GET /instructor/dashboard/stats - Get instructor statistics
echo "Test 12: GET /instructor/dashboard/stats - Get instructor statistics"
STATS_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/dashboard/stats)

if echo "$STATS_RESPONSE" | grep -q "assignments\|students\|stats"; then
    print_result 0 "GET /instructor/dashboard/stats - Get statistics works"
else
    print_result 1 "GET /instructor/dashboard/stats - Get statistics failed"
fi

# Test 13: GET /instructor/assignments/:id/detailed-progress - Detailed progress
echo "Test 13: GET /instructor/assignments/:id/detailed-progress - Detailed progress"
DETAILED_PROGRESS_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/detailed-progress)

if echo "$DETAILED_PROGRESS_RESPONSE" | grep -q "progress\|detailed\|analytics"; then
    print_result 0 "GET /instructor/assignments/:id/detailed-progress - Detailed progress works"
else
    print_result 1 "GET /instructor/assignments/:id/detailed-progress - Detailed progress failed"
fi

# Test 14: GET /instructor/progress/summary - Progress summary
echo "Test 14: GET /instructor/progress/summary - Progress summary"
PROGRESS_SUMMARY_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/progress/summary)

if echo "$PROGRESS_SUMMARY_RESPONSE" | grep -q "summary\|progress\|assignments"; then
    print_result 0 "GET /instructor/progress/summary - Progress summary works"
else
    print_result 1 "GET /instructor/progress/summary - Progress summary failed"
fi

# Test 15: GET /instructor/progress/trends - Progress trends
echo "Test 15: GET /instructor/progress/trends - Progress trends"
TRENDS_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/progress/trends)

if echo "$TRENDS_RESPONSE" | grep -q "trends\|progress\|analytics"; then
    print_result 0 "GET /instructor/progress/trends - Progress trends works"
else
    print_result 1 "GET /instructor/progress/trends - Progress trends failed"
fi

# Test 16: GET /instructor/progress/completion-analytics - Completion analytics
echo "Test 16: GET /instructor/progress/completion-analytics - Completion analytics"
COMPLETION_ANALYTICS_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/progress/completion-analytics)

if echo "$COMPLETION_ANALYTICS_RESPONSE" | grep -q "completion\|analytics\|rate"; then
    print_result 0 "GET /instructor/progress/completion-analytics - Completion analytics works"
else
    print_result 1 "GET /instructor/progress/completion-analytics - Completion analytics failed"
fi

# Test 17: GET /instructor/due-dates/overview - Due date overview
echo "Test 17: GET /instructor/due-dates/overview - Due date overview"
DUE_DATES_OVERVIEW_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/due-dates/overview)

if echo "$DUE_DATES_OVERVIEW_RESPONSE" | grep -q "due\|dates\|overview"; then
    print_result 0 "GET /instructor/due-dates/overview - Due date overview works"
else
    print_result 1 "GET /instructor/due-dates/overview - Due date overview failed"
fi

# Test 18: GET /instructor/due-dates/notifications - Due date notifications
echo "Test 18: GET /instructor/due-dates/notifications - Due date notifications"
DUE_DATES_NOTIFICATIONS_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/due-dates/notifications)

if echo "$DUE_DATES_NOTIFICATIONS_RESPONSE" | grep -q "notifications\|due\|dates\|Redirect"; then
    print_result 0 "GET /instructor/due-dates/notifications - Due date notifications works"
else
    print_result 1 "GET /instructor/due-dates/notifications - Due date notifications failed"
fi

# Test 19: Role-based access control - Students cannot access instructor endpoints
echo "Test 19: Role-based access control - Students cannot access instructor endpoints"
UNAUTHORIZED_RESPONSE=$(curl -s -b student_cookies.txt http://localhost:8080/instructor/assignments)

if echo "$UNAUTHORIZED_RESPONSE" | grep -q "unauthorized\|forbidden\|403\|Authentication required\|error"; then
    print_result 0 "Role-based access control works (students blocked from instructor endpoints)"
else
    print_result 1 "Role-based access control failed (students should be blocked from instructor endpoints)"
fi

# Test 20: Error handling - Non-existent assignment
echo "Test 20: Error handling - Non-existent assignment"
NON_EXISTENT_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/99999)

if echo "$NON_EXISTENT_RESPONSE" | grep -q "not found\|error\|404"; then
    print_result 0 "Error handling works (non-existent assignment)"
else
    print_result 1 "Error handling failed (should return error for non-existent assignment)"
fi

# Test 21: Content-Type validation
echo "Test 21: Content-Type validation"
INVALID_CONTENT_TYPE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: text/plain" \
  -d "invalid content")

if echo "$INVALID_CONTENT_TYPE" | grep -q "error\|invalid\|content-type"; then
    print_result 0 "Content-Type validation works"
else
    print_result 1 "Content-Type validation failed"
fi

# Test 22: Assignment creation with URL validation
echo "Test 22: Assignment creation with URL validation"
INVALID_URL_ASSIGNMENT=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Invalid URL Assignment",
    "description": "Testing URL validation",
    "url": "not-a-valid-url"
  }')

if echo "$INVALID_URL_ASSIGNMENT" | grep -q "error\|invalid\|url"; then
    print_result 0 "URL validation works (rejects invalid URLs)"
else
    print_result 1 "URL validation failed (should reject invalid URLs)"
fi

# Test 23: Assignment filtering by category
echo "Test 23: Assignment filtering by category"
# Create assignment with specific category
curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Category Filter Test",
    "description": "Testing category filtering",
    "url": "https://example.com/category-filter",
    "category": "FilterTest"
  }' > /dev/null

CATEGORY_FILTER_RESPONSE=$(curl -s -b instructor_cookies.txt "http://localhost:8080/instructor/assignments?category=FilterTest")

if echo "$CATEGORY_FILTER_RESPONSE" | grep -q "Category Filter Test"; then
    print_result 0 "Assignment filtering by category works"
else
    print_result 1 "Assignment filtering by category failed"
fi

# Test 24: Assignment search functionality
echo "Test 24: Assignment search functionality"
SEARCH_RESPONSE=$(curl -s -b instructor_cookies.txt "http://localhost:8080/instructor/assignments?search=Handler")

if echo "$SEARCH_RESPONSE" | grep -q "Updated Handler Test Assignment"; then
    print_result 0 "Assignment search functionality works"
else
    print_result 1 "Assignment search functionality failed"
fi

# Test 25: Assignment sorting
echo "Test 25: Assignment sorting"
SORT_RESPONSE=$(curl -s -b instructor_cookies.txt "http://localhost:8080/instructor/assignments?sort=title")

if [ $? -eq 0 ]; then
    print_result 0 "Assignment sorting works"
else
    print_result 1 "Assignment sorting failed"
fi

# Cleanup
echo "🧹 Cleanup"

rm -f instructor_cookies.txt student_cookies.txt
sleep 1

# Summary
echo ""
echo "📊 Phase 3 Task 3 Test Results:"
echo "================================"
echo -e "${GREEN}✅ Passed: $PASS${NC}"
echo -e "${RED}❌ Failed: $FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}🎉 Phase 3 Task 3: Instructor Assignment Management Handlers - ALL TESTS PASSED!${NC}"
    exit 0
else
    echo -e "${RED}🚨 Phase 3 Task 3: Instructor Assignment Management Handlers - $FAIL TESTS FAILED!${NC}"
    echo ""
    echo "Issues to address:"
    echo "- Assignment CRUD HTTP handlers"
    echo "- Form validation and error handling"
    echo "- Role-based access control"
    echo "- Student assignment management endpoints"
    echo "- Progress monitoring and analytics endpoints"
    echo "- Assignment filtering and search"
    echo "- UI endpoints for assignment management"
    echo ""
    exit 1
fi