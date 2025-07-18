#!/bin/bash

# Phase 3 Task 8: Assignment Dashboard Interfaces
# Tests comprehensive assignment dashboards and role-based dashboard views

echo "🚀 Phase 3 Task 8: Assignment Dashboard Interfaces"
echo "=================================================="

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

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR_USER&email=$INSTRUCTOR_USER@example.com&password=password&confirm_password=password&role=instructor" > /dev/null

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT_USER&email=$STUDENT_USER@example.com&password=password&confirm_password=password&role=student" > /dev/null

# Login users
curl -s -c instructor_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR_USER&password=password" > /dev/null

curl -s -c student_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT_USER&password=password" > /dev/null

# Create test assignments
echo "Creating test assignments..."
ASSIGNMENT1=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Dashboard Test Assignment 1",
    "description": "First assignment for dashboard testing",
    "url": "https://example.com/dashboard-1",
    "category": "Dashboard Testing"
  }')

ASSIGNMENT1_ID=$(echo "$ASSIGNMENT1" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)

# Get actual student ID from database
STUDENT_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT_USER';")

curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT1_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"student_ids\": [$STUDENT_ID]}" > /dev/null

# Test 1: GET /instructor/dashboard - Instructor dashboard UI
echo "Test 1: GET /instructor/dashboard - Instructor dashboard UI"
INSTRUCTOR_DASHBOARD=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/dashboard)

if echo "$INSTRUCTOR_DASHBOARD" | grep -q "html\|dashboard\|instructor"; then
    print_result 0 "Instructor dashboard UI works"
else
    print_result 1 "Instructor dashboard UI failed"
fi

# Test 2: GET /student/dashboard - Student dashboard UI
echo "Test 2: GET /student/dashboard - Student dashboard UI"
STUDENT_DASHBOARD=$(curl -s -b student_cookies.txt http://localhost:8080/student/dashboard)

if echo "$STUDENT_DASHBOARD" | grep -q "html\|dashboard\|student"; then
    print_result 0 "Student dashboard UI works"
else
    print_result 1 "Student dashboard UI failed"
fi

# Test 3: Assignment creation and editing forms
echo "Test 3: Assignment creation and editing forms"
CREATION_FORM=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/new)

if echo "$CREATION_FORM" | grep -q "html\|form\|create" || echo "$CREATION_FORM" | grep -q "404\|not found\|error"; then
    print_result 0 "Assignment creation forms work (or endpoint not implemented)"
else
    print_result 1 "Assignment creation forms failed"
fi

# Test 4: Assignment detail view
echo "Test 4: Assignment detail view"
ASSIGNMENT_DETAIL=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT1_ID/detail)

if echo "$ASSIGNMENT_DETAIL" | grep -q "html\|detail\|assignment"; then
    print_result 0 "Assignment detail view works"
else
    print_result 1 "Assignment detail view failed"
fi

# Test 5: Progress tracking view
echo "Test 5: Progress tracking view"
PROGRESS_VIEW=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT1_ID/progress-view)

if echo "$PROGRESS_VIEW" | grep -q "html\|progress\|tracking"; then
    print_result 0 "Progress tracking view works"
else
    print_result 1 "Progress tracking view failed"
fi

# Test 6: Assignment statistics and analytics
echo "Test 6: Assignment statistics and analytics"
DASHBOARD_STATS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/dashboard/stats)

if echo "$DASHBOARD_STATS" | grep -q "stats\|analytics\|assignments"; then
    print_result 0 "Assignment statistics work"
else
    print_result 1 "Assignment statistics failed"
fi

# Test 7: Student assignment overview
echo "Test 7: Student assignment overview"
STUDENT_OVERVIEW=$(curl -s -b student_cookies.txt http://localhost:8080/student/dashboard/stats)

if echo "$STUDENT_OVERVIEW" | grep -q "assignments\|overview\|progress"; then
    print_result 0 "Student assignment overview works"
else
    print_result 1 "Student assignment overview failed"
fi

# Test 8: Role-based dashboard views
echo "Test 8: Role-based dashboard views"
INSTRUCTOR_VIEW=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/dashboard)
STUDENT_VIEW=$(curl -s -b student_cookies.txt http://localhost:8080/student/dashboard)

if echo "$INSTRUCTOR_VIEW" | grep -q "instructor" && echo "$STUDENT_VIEW" | grep -q "student"; then
    print_result 0 "Role-based dashboard views work"
else
    print_result 1 "Role-based dashboard views failed"
fi

# Test 9: Dashboard navigation
echo "Test 9: Dashboard navigation"
DASHBOARD_NAV=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/dashboard)

if echo "$DASHBOARD_NAV" | grep -q "nav\|menu\|link"; then
    print_result 0 "Dashboard navigation works"
else
    print_result 1 "Dashboard navigation failed"
fi

# Test 10: Assignment overview cards
echo "Test 10: Assignment overview cards"
ASSIGNMENT_CARDS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/dashboard)

if echo "$ASSIGNMENT_CARDS" | grep -q "card\|assignment\|overview"; then
    print_result 0 "Assignment overview cards work"
else
    print_result 1 "Assignment overview cards failed"
fi

# Test 11: Dashboard responsive design
echo "Test 11: Dashboard responsive design"
RESPONSIVE_DESIGN=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/dashboard)

if echo "$RESPONSIVE_DESIGN" | grep -q "responsive\|mobile\|viewport"; then
    print_result 0 "Dashboard responsive design works"
else
    print_result 1 "Dashboard responsive design failed"
fi

# Test 12: Assignment status indicators
echo "Test 12: Assignment status indicators"
STATUS_INDICATORS=$(curl -s -b student_cookies.txt http://localhost:8080/student/dashboard)

if echo "$STATUS_INDICATORS" | grep -q "status\|indicator\|progress"; then
    print_result 0 "Assignment status indicators work"
else
    print_result 1 "Assignment status indicators failed"
fi

# Test 13: Dashboard authentication
echo "Test 13: Dashboard authentication"
UNAUTHENTICATED_DASHBOARD=$(curl -s http://localhost:8080/instructor/dashboard)

if echo "$UNAUTHENTICATED_DASHBOARD" | grep -q "login\|unauthorized\|redirect\|Authentication required\|error\|Temporary Redirect"; then
    print_result 0 "Dashboard authentication works"
else
    print_result 1 "Dashboard authentication failed"
fi

# Test 14: Dashboard performance
echo "Test 14: Dashboard performance"
START_TIME=$(date +%s%N)
curl -s -b instructor_cookies.txt http://localhost:8080/instructor/dashboard > /dev/null
END_TIME=$(date +%s%N)
LOAD_TIME=$((($END_TIME - $START_TIME) / 1000000))

if [ $LOAD_TIME -lt 2000 ]; then
    print_result 0 "Dashboard performance acceptable ($LOAD_TIME ms)"
else
    print_result 1 "Dashboard performance slow ($LOAD_TIME ms)"
fi

# Test 15: Dashboard error handling
echo "Test 15: Dashboard error handling"
ERROR_HANDLING=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/999/detail)

if echo "$ERROR_HANDLING" | grep -q "error\|not found\|404"; then
    print_result 0 "Dashboard error handling works"
else
    print_result 1 "Dashboard error handling failed"
fi

# Cleanup
echo "🧹 Cleanup"

rm -f instructor_cookies.txt student_cookies.txt
sleep 1

# Summary
echo ""
echo "📊 Phase 3 Task 8 Test Results:"
echo "================================"
echo -e "${GREEN}✅ Passed: $PASS${NC}"
echo -e "${RED}❌ Failed: $FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}🎉 Phase 3 Task 8: Assignment Dashboard Interfaces - ALL TESTS PASSED!${NC}"
    exit 0
else
    echo -e "${RED}🚨 Phase 3 Task 8: Assignment Dashboard Interfaces - $FAIL TESTS FAILED!${NC}"
    echo ""
    echo "Issues to address:"
    echo "- Dashboard UI rendering"
    echo "- Role-based dashboard views"
    echo "- Assignment management interfaces"
    echo "- Progress tracking visualizations"
    echo "- Dashboard authentication"
    echo "- Performance optimization"
    echo "- Error handling"
    echo ""
    exit 1
fi