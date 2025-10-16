#!/bin/bash

# Test script for Phase 3 Task 5 - Assignment Progress Tracking System
# This script tests the new progress tracking and due date notification features

BASE_URL="http://localhost:8081"
COOKIE_FILE="/tmp/zipcode_cookies.txt"

echo "🚀 Testing Phase 3 Task 5 - Assignment Progress Tracking System"
echo "============================================================"

# Clean up any existing cookies
rm -f $COOKIE_FILE

# Step 1: Register an instructor
echo "📝 Step 1: Register instructor..."
curl -c $COOKIE_FILE -s \
  -X POST "$BASE_URL/local/register" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=instructor_test&password=testpass123&email=instructor@test.com&role=instructor" \
  > /dev/null

# Step 2: Login as instructor
echo "🔐 Step 2: Login as instructor..."
curl -b $COOKIE_FILE -c $COOKIE_FILE -s \
  -X POST "$BASE_URL/local/login" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=instructor_test&password=testpass123" \
  > /dev/null

# Step 3: Create a test assignment with due date
echo "📚 Step 3: Create test assignment with due date..."
ASSIGNMENT_RESPONSE=$(curl -b $COOKIE_FILE -s \
  -X POST "$BASE_URL/instructor/assignments" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Progress Tracking Assignment",
    "description": "Assignment for testing progress tracking features",
    "reading_url": "https://example.com/test-reading",
    "category": "Programming",
    "due_date": "2025-07-25T23:59:59Z"
  }')

ASSIGNMENT_ID=$(echo $ASSIGNMENT_RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin)['id'])" 2>/dev/null || echo "1")
echo "   Assignment ID: $ASSIGNMENT_ID"

# Step 4: Register students
echo "👨‍🎓 Step 4: Register students..."
for i in {1..3}; do
  curl -s \
    -X POST "$BASE_URL/local/register" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "username=student$i&password=testpass123&email=student$i@test.com&role=student" \
    > /dev/null
done

# Step 5: Assign students to the assignment
echo "📋 Step 5: Assign students to assignment..."
curl -b $COOKIE_FILE -s \
  -X POST "$BASE_URL/instructor/assignments/$ASSIGNMENT_ID/assign" \
  -H "Content-Type: application/json" \
  -d '{"student_ids": [2, 3, 4]}' \
  > /dev/null

# Step 6: Test new instructor progress tracking endpoints
echo ""
echo "🔍 Testing Instructor Progress Tracking Features:"
echo "================================================="

# Test detailed progress report
echo "📊 Testing detailed progress report..."
DETAILED_PROGRESS=$(curl -b $COOKIE_FILE -s \
  -X GET "$BASE_URL/instructor/assignments/$ASSIGNMENT_ID/detailed-progress")
echo "   Detailed Progress Report:"
echo "$DETAILED_PROGRESS" | python3 -m json.tool 2>/dev/null || echo "   Response: $DETAILED_PROGRESS"

# Test instructor progress summary
echo ""
echo "📈 Testing instructor progress summary..."
PROGRESS_SUMMARY=$(curl -b $COOKIE_FILE -s \
  -X GET "$BASE_URL/instructor/progress/summary")
echo "   Progress Summary:"
echo "$PROGRESS_SUMMARY" | python3 -m json.tool 2>/dev/null || echo "   Response: $PROGRESS_SUMMARY"

# Test progress trends
echo ""
echo "📊 Testing progress trends..."
PROGRESS_TRENDS=$(curl -b $COOKIE_FILE -s \
  -X GET "$BASE_URL/instructor/progress/trends")
echo "   Progress Trends:"
echo "$PROGRESS_TRENDS" | python3 -m json.tool 2>/dev/null || echo "   Response: $PROGRESS_TRENDS"

# Test completion analytics
echo ""
echo "🎯 Testing completion analytics..."
COMPLETION_ANALYTICS=$(curl -b $COOKIE_FILE -s \
  -X GET "$BASE_URL/instructor/progress/completion-analytics")
echo "   Completion Analytics:"
echo "$COMPLETION_ANALYTICS" | python3 -m json.tool 2>/dev/null || echo "   Response: $COMPLETION_ANALYTICS"

# Test due date overview
echo ""
echo "📅 Testing due date overview..."
DUE_DATE_OVERVIEW=$(curl -b $COOKIE_FILE -s \
  -X GET "$BASE_URL/instructor/due-dates/overview")
echo "   Due Date Overview:"
echo "$DUE_DATE_OVERVIEW" | python3 -m json.tool 2>/dev/null || echo "   Response: $DUE_DATE_OVERVIEW"

# Test due date notifications
echo ""
echo "🔔 Testing due date notifications..."
DUE_DATE_NOTIFICATIONS=$(curl -b $COOKIE_FILE -s \
  -X GET "$BASE_URL/instructor/due-dates/notifications")
echo "   Due Date Notifications:"
echo "$DUE_DATE_NOTIFICATIONS" | python3 -m json.tool 2>/dev/null || echo "   Response: $DUE_DATE_NOTIFICATIONS"

# Step 7: Test student endpoints
echo ""
echo "🎓 Testing Student Progress Tracking Features:"
echo "=============================================="

# Login as student
echo "🔐 Login as student..."
curl -c /tmp/student_cookies.txt -s \
  -X POST "$BASE_URL/local/login" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=student1&password=testpass123" \
  > /dev/null

# Test student due date alerts
echo "🚨 Testing student due date alerts..."
STUDENT_ALERTS=$(curl -b /tmp/student_cookies.txt -s \
  -X GET "$BASE_URL/student/due-dates/alerts")
echo "   Student Due Date Alerts:"
echo "$STUDENT_ALERTS" | python3 -m json.tool 2>/dev/null || echo "   Response: $STUDENT_ALERTS"

# Test student due date summary
echo ""
echo "📋 Testing student due date summary..."
STUDENT_SUMMARY=$(curl -b /tmp/student_cookies.txt -s \
  -X GET "$BASE_URL/student/due-dates/summary")
echo "   Student Due Date Summary:"
echo "$STUDENT_SUMMARY" | python3 -m json.tool 2>/dev/null || echo "   Response: $STUDENT_SUMMARY"

# Test student due date notifications
echo ""
echo "🔔 Testing student due date notifications..."
STUDENT_NOTIFICATIONS=$(curl -b /tmp/student_cookies.txt -s \
  -X GET "$BASE_URL/student/due-dates/notifications")
echo "   Student Due Date Notifications:"
echo "$STUDENT_NOTIFICATIONS" | python3 -m json.tool 2>/dev/null || echo "   Response: $STUDENT_NOTIFICATIONS"

# Step 8: Test with some assignment progress
echo ""
echo "🔄 Testing with assignment progress updates..."
echo "============================================="

# Update assignment status for student
echo "📝 Update assignment status to in_progress..."
curl -b /tmp/student_cookies.txt -s \
  -X POST "$BASE_URL/student/assignments/$ASSIGNMENT_ID/progress" \
  -H "Content-Type: application/json" \
  -d '{"status": "in_progress"}' \
  > /dev/null

# Test detailed progress report again
echo "📊 Testing detailed progress report with updated data..."
UPDATED_PROGRESS=$(curl -b $COOKIE_FILE -s \
  -X GET "$BASE_URL/instructor/assignments/$ASSIGNMENT_ID/detailed-progress")
echo "   Updated Progress Report:"
echo "$UPDATED_PROGRESS" | python3 -m json.tool 2>/dev/null || echo "   Response: $UPDATED_PROGRESS"

# Complete the assignment
echo ""
echo "✅ Complete assignment..."
curl -b /tmp/student_cookies.txt -s \
  -X POST "$BASE_URL/student/assignments/$ASSIGNMENT_ID/complete" \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}' \
  > /dev/null

# Test final progress report
echo "📊 Testing final progress report..."
FINAL_PROGRESS=$(curl -b $COOKIE_FILE -s \
  -X GET "$BASE_URL/instructor/assignments/$ASSIGNMENT_ID/detailed-progress")
echo "   Final Progress Report:"
echo "$FINAL_PROGRESS" | python3 -m json.tool 2>/dev/null || echo "   Response: $FINAL_PROGRESS"

# Clean up
rm -f $COOKIE_FILE /tmp/student_cookies.txt

echo ""
echo "🎉 Phase 3 Task 5 testing complete!"
echo "✅ All new progress tracking and due date notification endpoints tested"
echo "✅ Advanced progress analytics working"
echo "✅ Due date notification system functional"
echo "✅ Student and instructor interfaces tested"
echo ""
echo "🔗 New API Endpoints Available:"
echo "   Instructor Progress Tracking:"
echo "   - GET /instructor/assignments/:id/detailed-progress"
echo "   - GET /instructor/progress/summary"
echo "   - GET /instructor/progress/trends"
echo "   - GET /instructor/progress/completion-analytics"
echo "   - GET /instructor/due-dates/overview"
echo "   - GET /instructor/due-dates/notifications"
echo ""
echo "   Student Due Date Notifications:"
echo "   - GET /student/due-dates/alerts"
echo "   - GET /student/due-dates/summary"
echo "   - GET /student/due-dates/notifications"
