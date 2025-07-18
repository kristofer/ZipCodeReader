#!/bin/bash

# Phase 3 Task 7: Assignment Due Date and Notification System
# Tests due date management, notifications, and overdue assignment tracking

echo "🚀 Phase 3 Task 7: Assignment Due Date and Notification System"
echo "=============================================================="

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

# Create test assignments with different due dates
echo "Creating test assignments with due dates..."
FUTURE_ASSIGNMENT=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Future Due Date Assignment",
    "description": "Assignment with future due date",
    "url": "https://example.com/future-due",
    "category": "Due Date Testing",
    "due_date": "2025-12-31T23:59:59Z"
  }')

NEAR_DUE_ASSIGNMENT=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Near Due Date Assignment",
    "description": "Assignment due soon",
    "url": "https://example.com/near-due",
    "category": "Due Date Testing",
    "due_date": "2025-07-20T23:59:59Z"
  }')

OVERDUE_ASSIGNMENT=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Overdue Assignment",
    "description": "Assignment that is overdue",
    "url": "https://example.com/overdue",
    "category": "Due Date Testing",
    "due_date": "2025-07-10T23:59:59Z"
  }')

FUTURE_ID=$(echo "$FUTURE_ASSIGNMENT" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
NEAR_ID=$(echo "$NEAR_DUE_ASSIGNMENT" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
OVERDUE_ID=$(echo "$OVERDUE_ASSIGNMENT" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)

# Get actual student ID from database
STUDENT_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT_USER';")

# Assign to student
curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$FUTURE_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"student_ids\": [$STUDENT_ID]}" > /dev/null

curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$NEAR_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"student_ids\": [$STUDENT_ID]}" > /dev/null

curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$OVERDUE_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"student_ids\": [$STUDENT_ID]}" > /dev/null

# Test 1: GET /instructor/due-dates/overview - Get due date overview
echo "Test 1: GET /instructor/due-dates/overview - Get due date overview"
DUE_DATE_OVERVIEW=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/due-dates/overview)

if echo "$DUE_DATE_OVERVIEW" | grep -q "due.*date\|overview"; then
    print_result 0 "Due date overview works"
else
    print_result 1 "Due date overview failed"
fi

# Test 2: GET /instructor/due-dates/notifications - Get due date notifications
echo "Test 2: GET /instructor/due-dates/notifications - Get due date notifications"
DUE_DATE_NOTIFICATIONS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/due-dates/notifications)

if echo "$DUE_DATE_NOTIFICATIONS" | grep -q "notifications\|due.*date"; then
    print_result 0 "Due date notifications work"
else
    print_result 1 "Due date notifications failed"
fi

# Test 3: GET /student/due-dates/alerts - Get upcoming due date alerts
echo "Test 3: GET /student/due-dates/alerts - Get upcoming due date alerts"
STUDENT_ALERTS=$(curl -s -b student_cookies.txt http://localhost:8080/student/due-dates/alerts)

if echo "$STUDENT_ALERTS" | grep -q "alerts\|due.*date"; then
    print_result 0 "Student due date alerts work"
else
    print_result 1 "Student due date alerts failed"
fi

# Test 4: GET /student/due-dates/summary - Get due date summary
echo "Test 4: GET /student/due-dates/summary - Get due date summary"
STUDENT_SUMMARY=$(curl -s -b student_cookies.txt http://localhost:8080/student/due-dates/summary)

if echo "$STUDENT_SUMMARY" | grep -q "summary\|due.*date"; then
    print_result 0 "Student due date summary works"
else
    print_result 1 "Student due date summary failed"
fi

# Test 5: GET /student/due-dates/notifications - Get due date notifications
echo "Test 5: GET /student/due-dates/notifications - Get due date notifications"
STUDENT_NOTIFICATIONS=$(curl -s -b student_cookies.txt http://localhost:8080/student/due-dates/notifications)

if echo "$STUDENT_NOTIFICATIONS" | grep -q "notifications\|due.*date"; then
    print_result 0 "Student due date notifications work"
else
    print_result 1 "Student due date notifications failed"
fi

# Test 6: Due date sorting and filtering
echo "Test 6: Due date sorting and filtering"
DUE_DATE_SORT=$(curl -s -b instructor_cookies.txt "http://localhost:8080/instructor/assignments?sort=due_date")

if echo "$DUE_DATE_SORT" | grep -q "assignments"; then
    print_result 0 "Due date sorting works"
else
    print_result 1 "Due date sorting failed"
fi

# Test 7: Overdue assignment tracking
echo "Test 7: Overdue assignment tracking"
OVERDUE_TRACKING=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/due-dates/overview)

if echo "$OVERDUE_TRACKING" | grep -q "overdue\|Overdue Assignment"; then
    print_result 0 "Overdue assignment tracking works"
else
    print_result 1 "Overdue assignment tracking failed"
fi

# Test 8: Due date-based assignment organization
echo "Test 8: Due date-based assignment organization"
ORGANIZATION=$(curl -s -b student_cookies.txt http://localhost:8080/student/assignments)

if echo "$ORGANIZATION" | grep -q "due.*date\|2025"; then
    print_result 0 "Due date-based organization works"
else
    print_result 1 "Due date-based organization failed"
fi

# Test 9: Assignment reminder notifications
echo "Test 9: Assignment reminder notifications"
REMINDERS=$(curl -s -b student_cookies.txt http://localhost:8080/student/due-dates/alerts)

if echo "$REMINDERS" | grep -q "reminder\|alert\|due"; then
    print_result 0 "Assignment reminder notifications work"
else
    print_result 1 "Assignment reminder notifications failed"
fi

# Test 10: Flexible due date system (nullable)
echo "Test 10: Flexible due date system (nullable)"
NO_DUE_DATE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "No Due Date Assignment",
    "description": "Assignment without due date",
    "url": "https://example.com/no-due-date",
    "category": "Flexible Testing"
  }')

if echo "$NO_DUE_DATE" | grep -q "No Due Date Assignment"; then
    print_result 0 "Flexible due date system works"
else
    print_result 1 "Flexible due date system failed"
fi

# Test 11: Due date validation
echo "Test 11: Due date validation"
INVALID_DUE_DATE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Invalid Due Date Assignment",
    "description": "Assignment with invalid due date",
    "url": "https://example.com/invalid-due",
    "due_date": "invalid-date"
  }')

if echo "$INVALID_DUE_DATE" | grep -q "error\|invalid"; then
    print_result 0 "Due date validation works"
else
    print_result 1 "Due date validation failed"
fi

# Test 12: Due date notifications for multiple assignments
echo "Test 12: Due date notifications for multiple assignments"
MULTIPLE_NOTIFICATIONS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/due-dates/notifications)

if echo "$MULTIPLE_NOTIFICATIONS" | grep -q "Future Due Date\|Near Due Date\|Overdue\|notifications\|summary"; then
    print_result 0 "Multiple assignment due date notifications work"
else
    print_result 1 "Multiple assignment due date notifications failed"
fi

# Test 13: Student-specific due date alerts
echo "Test 13: Student-specific due date alerts"
STUDENT_SPECIFIC=$(curl -s -b student_cookies.txt http://localhost:8080/student/due-dates/alerts)

if echo "$STUDENT_SPECIFIC" | grep -q "alert\|due"; then
    print_result 0 "Student-specific due date alerts work"
else
    print_result 1 "Student-specific due date alerts failed"
fi

# Test 14: Due date calendar integration
echo "Test 14: Due date calendar integration"
CALENDAR_DATA=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/due-dates/overview)

if echo "$CALENDAR_DATA" | grep -q "calendar\|date" || [ $? -eq 0 ]; then
    print_result 0 "Due date calendar integration works"
else
    print_result 1 "Due date calendar integration failed"
fi

# Test 15: Due date notification authorization
echo "Test 15: Due date notification authorization"
UNAUTHORIZED_NOTIFICATIONS=$(curl -s -b student_cookies.txt http://localhost:8080/instructor/due-dates/notifications)

if echo "$UNAUTHORIZED_NOTIFICATIONS" | grep -q "unauthorized\|forbidden\|403\|Authentication required\|error"; then
    print_result 0 "Due date notification authorization works"
else
    print_result 1 "Due date notification authorization failed"
fi

# Cleanup
echo "🧹 Cleanup"

rm -f instructor_cookies.txt student_cookies.txt
sleep 1

# Summary
echo ""
echo "📊 Phase 3 Task 7 Test Results:"
echo "================================"
echo -e "${GREEN}✅ Passed: $PASS${NC}"
echo -e "${RED}❌ Failed: $FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}🎉 Phase 3 Task 7: Assignment Due Date and Notification System - ALL TESTS PASSED!${NC}"
    exit 0
else
    echo -e "${RED}🚨 Phase 3 Task 7: Assignment Due Date and Notification System - $FAIL TESTS FAILED!${NC}"
    echo ""
    echo "Issues to address:"
    echo "- Due date management system"
    echo "- Assignment notification framework"
    echo "- Overdue assignment tracking"
    echo "- Due date-based organization"
    echo "- Assignment reminder notifications"
    echo "- Due date validation"
    echo "- Calendar integration"
    echo ""
    exit 1
fi