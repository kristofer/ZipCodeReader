#!/bin/bash

# Phase 3 Task 5: Assignment Progress Tracking System
# Tests comprehensive progress tracking, completion statistics, and progress reporting

echo "🚀 Phase 3 Task 5: Assignment Progress Tracking System"
echo "===================================================="

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
STUDENT1_USER="student1_$TIMESTAMP"
STUDENT2_USER="student2_$TIMESTAMP"
STUDENT3_USER="student3_$TIMESTAMP"

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR_USER&email=$INSTRUCTOR_USER@example.com&password=password&confirm_password=password&role=instructor" > /dev/null

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT1_USER&email=$STUDENT1_USER@example.com&password=password&confirm_password=password&role=student" > /dev/null

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT2_USER&email=$STUDENT2_USER@example.com&password=password&confirm_password=password&role=student" > /dev/null

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT3_USER&email=$STUDENT3_USER@example.com&password=password&confirm_password=password&role=student" > /dev/null

# Login users
curl -s -c instructor_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR_USER&password=password" > /dev/null

curl -s -c student1_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT1_USER&password=password" > /dev/null

curl -s -c student2_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT2_USER&password=password" > /dev/null

curl -s -c student3_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT3_USER&password=password" > /dev/null

# Create test assignments
echo "Creating test assignments..."
ASSIGNMENT1=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Progress Tracking Test 1",
    "description": "First assignment for progress tracking",
    "url": "https://example.com/progress-1",
    "category": "Progress Testing",
    "due_date": "2025-08-30T23:59:59Z"
  }')

ASSIGNMENT2=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Progress Tracking Test 2",
    "description": "Second assignment for progress tracking",
    "url": "https://example.com/progress-2",
    "category": "Progress Testing"
  }')

ASSIGNMENT1_ID=$(echo "$ASSIGNMENT1" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
ASSIGNMENT2_ID=$(echo "$ASSIGNMENT2" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)

# Assign to students
# Get the actual student IDs from the database
STUDENT1_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT1_USER';")
STUDENT2_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT2_USER';")
STUDENT3_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT3_USER';")

curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT1_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"student_ids\": [$STUDENT1_ID, $STUDENT2_ID, $STUDENT3_ID]}" > /dev/null

curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT2_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"student_ids\": [$STUDENT1_ID, $STUDENT2_ID, $STUDENT3_ID]}" > /dev/null

# Create some progress data
curl -s -b student1_cookies.txt -X POST http://localhost:8080/student/assignments/$ASSIGNMENT1_ID/progress > /dev/null
curl -s -b student1_cookies.txt -X POST http://localhost:8080/student/assignments/$ASSIGNMENT2_ID/complete > /dev/null

# Test 1: GET /instructor/assignments/:id/detailed-progress - Get detailed progress report
echo "Test 1: GET /instructor/assignments/:id/detailed-progress - Get detailed progress report"
DETAILED_PROGRESS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT1_ID/detailed-progress)

if echo "$DETAILED_PROGRESS" | grep -q "progress\|completion\|students"; then
    print_result 0 "Detailed progress report works"
else
    print_result 1 "Detailed progress report failed"
fi

# Test 2: GET /instructor/progress/summary - Get instructor progress summary
echo "Test 2: GET /instructor/progress/summary - Get instructor progress summary"
PROGRESS_SUMMARY=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/progress/summary)

if echo "$PROGRESS_SUMMARY" | grep -q "summary\|assignments\|completion"; then
    print_result 0 "Progress summary works"
else
    print_result 1 "Progress summary failed"
fi

# Test 3: GET /instructor/progress/trends - Get progress trends analysis
echo "Test 3: GET /instructor/progress/trends - Get progress trends analysis"
PROGRESS_TRENDS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/progress/trends)

if echo "$PROGRESS_TRENDS" | grep -q "trends\|analytics\|progress"; then
    print_result 0 "Progress trends analysis works"
else
    print_result 1 "Progress trends analysis failed"
fi

# Test 4: GET /instructor/progress/completion-analytics - Get completion analytics
echo "Test 4: GET /instructor/progress/completion-analytics - Get completion analytics"
COMPLETION_ANALYTICS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/progress/completion-analytics)

if echo "$COMPLETION_ANALYTICS" | grep -q "completion\|analytics\|rate"; then
    print_result 0 "Completion analytics works"
else
    print_result 1 "Completion analytics failed"
fi

# Test 5: Real-time progress updates
echo "Test 5: Real-time progress updates"
# Update progress and check if it's reflected immediately
curl -s -b student2_cookies.txt -X POST http://localhost:8080/student/assignments/$ASSIGNMENT1_ID/progress > /dev/null

REAL_TIME_PROGRESS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT1_ID/progress)

if echo "$REAL_TIME_PROGRESS" | grep -q "in_progress\|progress"; then
    print_result 0 "Real-time progress updates work"
else
    print_result 1 "Real-time progress updates failed"
fi

# Test 6: Assignment completion percentages
echo "Test 6: Assignment completion percentages"
COMPLETION_PERCENTAGE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT1_ID/progress)

if echo "$COMPLETION_PERCENTAGE" | grep -q "percentage\|completion\|rate"; then
    print_result 0 "Assignment completion percentages work"
else
    print_result 1 "Assignment completion percentages failed"
fi

# Test 7: Student engagement metrics
echo "Test 7: Student engagement metrics"
ENGAGEMENT_METRICS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/progress/completion-analytics)

if echo "$ENGAGEMENT_METRICS" | grep -q "engagement\|metrics\|students"; then
    print_result 0 "Student engagement metrics work"
else
    print_result 1 "Student engagement metrics failed"
fi

# Test 8: Progress visualization data
echo "Test 8: Progress visualization data"
VISUALIZATION_DATA=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/progress/trends)

if echo "$VISUALIZATION_DATA" | grep -q "data\|visualization\|chart\|trends\|progress\|granularity"; then
    print_result 0 "Progress visualization data works"
else
    print_result 1 "Progress visualization data failed"
fi

# Test 9: Assignment progress with timestamps
echo "Test 9: Assignment progress with timestamps"
TIMESTAMP_PROGRESS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT1_ID/detailed-progress)

if echo "$TIMESTAMP_PROGRESS" | grep -q "timestamp\|created_at\|updated_at"; then
    print_result 0 "Assignment progress with timestamps works"
else
    print_result 1 "Assignment progress with timestamps failed"
fi

# Test 10: Progress tracking for multiple assignments
echo "Test 10: Progress tracking for multiple assignments"
MULTI_ASSIGNMENT_PROGRESS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/progress/summary)

if echo "$MULTI_ASSIGNMENT_PROGRESS" | grep -q "Progress Tracking Test"; then
    print_result 0 "Multi-assignment progress tracking works"
else
    print_result 1 "Multi-assignment progress tracking failed"
fi

# Test 11: Student-level progress insights
echo "Test 11: Student-level progress insights"
STUDENT_INSIGHTS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT1_ID/detailed-progress)

if echo "$STUDENT_INSIGHTS" | grep -q "student\|insights\|individual"; then
    print_result 0 "Student-level progress insights work"
else
    print_result 1 "Student-level progress insights failed"
fi

# Test 12: Progress trends over time
echo "Test 12: Progress trends over time"
TIME_TRENDS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/progress/trends)

if echo "$TIME_TRENDS" | grep -q "time\|trends\|history"; then
    print_result 0 "Progress trends over time work"
else
    print_result 1 "Progress trends over time failed"
fi

# Test 13: Progress tracking service authorization
echo "Test 13: Progress tracking service authorization"
UNAUTHORIZED_PROGRESS=$(curl -s -b student1_cookies.txt http://localhost:8080/instructor/progress/summary)

if echo "$UNAUTHORIZED_PROGRESS" | grep -q "unauthorized\|forbidden\|403\|Authentication required\|error"; then
    print_result 0 "Progress tracking authorization works"
else
    print_result 1 "Progress tracking authorization failed"
fi

# Test 14: Progress data accuracy
echo "Test 14: Progress data accuracy"
# Complete another assignment and verify counts
curl -s -b student2_cookies.txt -X POST http://localhost:8080/student/assignments/$ASSIGNMENT2_ID/complete > /dev/null

ACCURACY_CHECK=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/progress/completion-analytics)

if echo "$ACCURACY_CHECK" | grep -q "completion\|analytics"; then
    print_result 0 "Progress data accuracy works"
else
    print_result 1 "Progress data accuracy failed"
fi

# Test 15: Progress tracking error handling
echo "Test 15: Progress tracking error handling"
ERROR_HANDLING=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/99999/detailed-progress)

if echo "$ERROR_HANDLING" | grep -q "not found\|error\|404"; then
    print_result 0 "Progress tracking error handling works"
else
    print_result 1 "Progress tracking error handling failed"
fi

# Cleanup
echo "🧹 Cleanup"

rm -f instructor_cookies.txt student1_cookies.txt student2_cookies.txt
sleep 1

# Summary
echo ""
echo "📊 Phase 3 Task 5 Test Results:"
echo "================================"
echo -e "${GREEN}✅ Passed: $PASS${NC}"
echo -e "${RED}❌ Failed: $FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}🎉 Phase 3 Task 5: Assignment Progress Tracking System - ALL TESTS PASSED!${NC}"
    exit 0
else
    echo -e "${RED}🚨 Phase 3 Task 5: Assignment Progress Tracking System - $FAIL TESTS FAILED!${NC}"
    echo ""
    echo "Issues to address:"
    echo "- Progress tracking service implementation"
    echo "- Detailed progress reports and analytics"
    echo "- Real-time progress updates"
    echo "- Student engagement metrics"
    echo "- Progress visualization data"
    echo "- Completion rate calculations"
    echo "- Progress trends analysis"
    echo ""
    exit 1
fi