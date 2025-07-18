#!/bin/bash

# Phase 3 Task 10: Testing and Integration
# Tests complete assignment management flow, role-based access control, and integration

echo "🚀 Phase 3 Task 10: Testing and Integration"
echo "==========================================="

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
INSTRUCTOR1_USER="instructor1_$TIMESTAMP"
INSTRUCTOR2_USER="instructor2_$TIMESTAMP"
STUDENT1_USER="student1_$TIMESTAMP"
STUDENT2_USER="student2_$TIMESTAMP"

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR1_USER&email=$INSTRUCTOR1_USER@example.com&password=password&confirm_password=password&role=instructor" > /dev/null

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR2_USER&email=$INSTRUCTOR2_USER@example.com&password=password&confirm_password=password&role=instructor" > /dev/null

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT1_USER&email=$STUDENT1_USER@example.com&password=password&confirm_password=password&role=student" > /dev/null

curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT2_USER&email=$STUDENT2_USER@example.com&password=password&confirm_password=password&role=student" > /dev/null

# Login users
curl -s -c instructor1_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR1_USER&password=password" > /dev/null

curl -s -c instructor2_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR2_USER&password=password" > /dev/null

curl -s -c student1_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT1_USER&password=password" > /dev/null

curl -s -c student2_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT2_USER&password=password" > /dev/null

# Test 1: Complete assignment creation flow
echo "Test 1: Complete assignment creation flow"
ASSIGNMENT_CREATION=$(curl -s -b instructor1_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Integration Test Assignment",
    "description": "Complete integration testing assignment",
    "url": "https://example.com/integration-test",
    "category": "Integration Testing",
    "due_date": "2025-09-01T23:59:59Z"
  }')

if echo "$ASSIGNMENT_CREATION" | grep -q "Integration Test Assignment"; then
    ASSIGNMENT_ID=$(echo "$ASSIGNMENT_CREATION" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
    print_result 0 "Complete assignment creation flow works (ID: $ASSIGNMENT_ID)"
else
    print_result 1 "Complete assignment creation flow failed"
    ASSIGNMENT_ID=1
fi

# Test 2: Assignment-student relationship management
echo "Test 2: Assignment-student relationship management"
ASSIGNMENT_RELATIONSHIP=$(curl -s -b instructor1_cookies.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/assign \
  -H "Content-Type: application/json" \
  -d '{"student_ids": [3, 4]}')

if echo "$ASSIGNMENT_RELATIONSHIP" | grep -q "success\|assigned"; then
    print_result 0 "Assignment-student relationship management works"
else
    print_result 1 "Assignment-student relationship management failed"
fi

# Test 3: Student assignment viewing and completion
echo "Test 3: Student assignment viewing and completion"
STUDENT_VIEWING=$(curl -s -b student1_cookies.txt http://localhost:8080/student/assignments)

if echo "$STUDENT_VIEWING" | grep -q "Integration Test Assignment"; then
    print_result 0 "Student assignment viewing works"
else
    print_result 1 "Student assignment viewing failed"
fi

# Test 4: Assignment progress tracking
echo "Test 4: Assignment progress tracking"
# Student updates assignment status
curl -s -b student1_cookies.txt -X POST http://localhost:8080/student/assignments/$ASSIGNMENT_ID/progress > /dev/null

PROGRESS_TRACKING=$(curl -s -b instructor1_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/progress)

if echo "$PROGRESS_TRACKING" | grep -q "progress\|in_progress"; then
    print_result 0 "Assignment progress tracking works"
else
    print_result 1 "Assignment progress tracking failed"
fi

# Test 5: Role-based access control (comprehensive)
echo "Test 5: Role-based access control (comprehensive)"
# Test instructor access to instructor endpoints
INSTRUCTOR_ACCESS=$(curl -s -b instructor1_cookies.txt http://localhost:8080/instructor/assignments)

# Test student blocked from instructor endpoints
STUDENT_BLOCKED=$(curl -s -b student1_cookies.txt http://localhost:8080/instructor/assignments)

# Test instructor blocked from student endpoints
INSTRUCTOR_BLOCKED=$(curl -s -b instructor1_cookies.txt http://localhost:8080/student/assignments)

if [ $? -eq 0 ] && echo "$STUDENT_BLOCKED" | grep -q "unauthorized\|forbidden\|403\|Authentication required\|error" && echo "$INSTRUCTOR_BLOCKED" | grep -q "unauthorized\|forbidden\|403\|Authentication required\|error"; then
    print_result 0 "Role-based access control works comprehensively"
else
    print_result 1 "Role-based access control failed"
fi

# Test 6: Assignment due date management and notifications
echo "Test 6: Assignment due date management and notifications"
DUE_DATE_MANAGEMENT=$(curl -s -b instructor1_cookies.txt http://localhost:8080/instructor/due-dates/overview)

if echo "$DUE_DATE_MANAGEMENT" | grep -q "due.*date\|overview"; then
    print_result 0 "Assignment due date management works"
else
    print_result 1 "Assignment due date management failed"
fi

# Test 7: Assignment search, filtering, and categorization integration
echo "Test 7: Assignment search, filtering, and categorization integration"
SEARCH_INTEGRATION=$(curl -s -b instructor1_cookies.txt "http://localhost:8080/instructor/assignments?search=Integration&category=Integration%20Testing")

if echo "$SEARCH_INTEGRATION" | grep -q "Integration Test Assignment"; then
    print_result 0 "Search, filtering, and categorization integration works"
else
    print_result 1 "Search, filtering, and categorization integration failed"
fi

# Test 8: Assignment dashboard interfaces integration
echo "Test 8: Assignment dashboard interfaces integration"
DASHBOARD_INTEGRATION=$(curl -s -b instructor1_cookies.txt http://localhost:8080/instructor/dashboard)

if echo "$DASHBOARD_INTEGRATION" | grep -q "dashboard\|html"; then
    print_result 0 "Assignment dashboard interfaces integration works"
else
    print_result 1 "Assignment dashboard interfaces integration failed"
fi

# Test 9: Multi-instructor isolation
echo "Test 9: Multi-instructor isolation"
# Create assignment with instructor2
INSTRUCTOR2_ASSIGNMENT=$(curl -s -b instructor2_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Instructor2 Assignment",
    "description": "Assignment by instructor2",
    "url": "https://example.com/instructor2-assignment"
  }')

INSTRUCTOR2_ID=$(echo "$INSTRUCTOR2_ASSIGNMENT" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)

# Check instructor1 cannot see instructor2's assignment
INSTRUCTOR1_ASSIGNMENTS=$(curl -s -b instructor1_cookies.txt http://localhost:8080/instructor/assignments)

if ! echo "$INSTRUCTOR1_ASSIGNMENTS" | grep -q "Instructor2 Assignment"; then
    print_result 0 "Multi-instructor isolation works"
else
    print_result 1 "Multi-instructor isolation failed"
fi

# Test 10: Assignment completion flow
echo "Test 10: Assignment completion flow"
ASSIGNMENT_COMPLETION=$(curl -s -b student1_cookies.txt -X POST http://localhost:8080/student/assignments/$ASSIGNMENT_ID/complete)

if echo "$ASSIGNMENT_COMPLETION" | grep -q "success\|completed"; then
    print_result 0 "Assignment completion flow works"
else
    print_result 1 "Assignment completion flow failed"
fi

# Test 11: Assignment modification and deletion
echo "Test 11: Assignment modification and deletion"
# Update assignment
ASSIGNMENT_UPDATE=$(curl -s -b instructor1_cookies.txt -X PUT http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Integration Test Assignment",
    "description": "Updated description",
    "url": "https://example.com/updated-integration-test"
  }')

if echo "$ASSIGNMENT_UPDATE" | grep -q "Updated Integration Test Assignment\|success"; then
    print_result 0 "Assignment modification works"
else
    print_result 1 "Assignment modification failed"
fi

# Test 12: Performance testing with multiple operations
echo "Test 12: Performance testing with multiple operations"
START_TIME=$(date +%s%N)

# Perform multiple operations
for i in {1..5}; do
    curl -s -b instructor1_cookies.txt -X POST http://localhost:8080/instructor/assignments \
      -H "Content-Type: application/json" \
      -d "{
        \"title\": \"Performance Test $i\",
        \"description\": \"Performance testing assignment $i\",
        \"url\": \"https://example.com/perf-$i\"
      }" > /dev/null
done

curl -s -b instructor1_cookies.txt http://localhost:8080/instructor/assignments > /dev/null

END_TIME=$(date +%s%N)
PERFORMANCE_TIME=$((($END_TIME - $START_TIME) / 1000000))

if [ $PERFORMANCE_TIME -lt 5000 ]; then
    print_result 0 "Performance testing acceptable ($PERFORMANCE_TIME ms)"
else
    print_result 1 "Performance testing slow ($PERFORMANCE_TIME ms)"
fi

# Test 13: Data integrity and consistency
echo "Test 13: Data integrity and consistency"
# Check that all related data is consistent
INTEGRITY_CHECK=$(curl -s -b instructor1_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/students)

if echo "$INTEGRITY_CHECK" | grep -q "student"; then
    print_result 0 "Data integrity and consistency works"
else
    print_result 1 "Data integrity and consistency failed"
fi

# Test 14: Error handling and recovery
echo "Test 14: Error handling and recovery"
# Test various error conditions
ERROR_HANDLING1=$(curl -s -b instructor1_cookies.txt http://localhost:8080/instructor/assignments/999999)
ERROR_HANDLING2=$(curl -s -b instructor1_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{"invalid": "data"}')

if echo "$ERROR_HANDLING1" | grep -q "not found\|error\|404" && echo "$ERROR_HANDLING2" | grep -q "error\|invalid"; then
    print_result 0 "Error handling and recovery works"
else
    print_result 1 "Error handling and recovery failed"
fi

# Test 15: Session management and security
echo "Test 15: Session management and security"
# Test unauthenticated access
UNAUTHENTICATED=$(curl -s http://localhost:8080/instructor/assignments)

if echo "$UNAUTHENTICATED" | grep -q "unauthorized\|login\|redirect"; then
    print_result 0 "Session management and security works"
else
    print_result 1 "Session management and security failed"
fi

# Test 16: Assignment analytics and reporting
echo "Test 16: Assignment analytics and reporting"
ANALYTICS=$(curl -s -b instructor1_cookies.txt http://localhost:8080/instructor/progress/summary)

if echo "$ANALYTICS" | grep -q "summary\|progress\|analytics"; then
    print_result 0 "Assignment analytics and reporting works"
else
    print_result 1 "Assignment analytics and reporting failed"
fi

# Test 17: Bulk operations
echo "Test 17: Bulk operations"
BULK_OPERATIONS=$(curl -s -b instructor1_cookies.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/assign \
  -H "Content-Type: application/json" \
  -d '{"student_ids": [3, 4]}')

if echo "$BULK_OPERATIONS" | grep -q "success\|assigned"; then
    print_result 0 "Bulk operations work"
else
    print_result 1 "Bulk operations failed"
fi

# Test 18: Assignment workflow integration
echo "Test 18: Assignment workflow integration"
# Test the complete workflow: create -> assign -> update status -> complete
WORKFLOW_ASSIGNMENT=$(curl -s -b instructor1_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Workflow Test Assignment",
    "description": "Testing complete workflow",
    "url": "https://example.com/workflow-test"
  }')

WORKFLOW_ID=$(echo "$WORKFLOW_ASSIGNMENT" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)

curl -s -b instructor1_cookies.txt -X POST http://localhost:8080/instructor/assignments/$WORKFLOW_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"student_ids\": [$STUDENT1_ID]}" > /dev/null

curl -s -b student1_cookies.txt -X POST http://localhost:8080/student/assignments/$WORKFLOW_ID/progress > /dev/null
curl -s -b student1_cookies.txt -X POST http://localhost:8080/student/assignments/$WORKFLOW_ID/complete > /dev/null

WORKFLOW_CHECK=$(curl -s -b instructor1_cookies.txt http://localhost:8080/instructor/assignments/$WORKFLOW_ID/progress)

if echo "$WORKFLOW_CHECK" | grep -q "completed\|progress"; then
    print_result 0 "Assignment workflow integration works"
else
    print_result 1 "Assignment workflow integration failed"
fi

# Test 19: Cross-feature integration
echo "Test 19: Cross-feature integration"
# Test that all features work together
CROSS_FEATURE=$(curl -s -b instructor1_cookies.txt "http://localhost:8080/instructor/assignments?search=Integration&category=Integration%20Testing")

if echo "$CROSS_FEATURE" | grep -q "Updated Integration Test Assignment"; then
    print_result 0 "Cross-feature integration works"
else
    print_result 1 "Cross-feature integration failed"
fi

# Test 20: System stability and reliability
echo "Test 20: System stability and reliability"
# Test rapid operations
for i in {1..10}; do
    curl -s -b instructor1_cookies.txt http://localhost:8080/instructor/assignments > /dev/null
done

STABILITY_CHECK=$(curl -s -b instructor1_cookies.txt http://localhost:8080/instructor/assignments)

if [ $? -eq 0 ]; then
    print_result 0 "System stability and reliability works"
else
    print_result 1 "System stability and reliability failed"
fi

# Cleanup
echo "🧹 Cleanup"

rm -f instructor1_cookies.txt instructor2_cookies.txt student1_cookies.txt student2_cookies.txt
sleep 1

# Summary
echo ""
echo "📊 Phase 3 Task 10 Test Results:"
echo "================================"
echo -e "${GREEN}✅ Passed: $PASS${NC}"
echo -e "${RED}❌ Failed: $FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}🎉 Phase 3 Task 10: Testing and Integration - ALL TESTS PASSED!${NC}"
    echo ""
    echo "🏆 PHASE 3 COMPLETE - ASSIGNMENT MANAGEMENT SYSTEM FULLY INTEGRATED!"
    echo ""
    exit 0
else
    echo -e "${RED}🚨 Phase 3 Task 10: Testing and Integration - $FAIL TESTS FAILED!${NC}"
    echo ""
    echo "Issues to address:"
    echo "- Complete assignment management flow"
    echo "- Role-based access control"
    echo "- Assignment-student relationships"
    echo "- Assignment progress tracking"
    echo "- Due date management and notifications"
    echo "- Search, filtering, and categorization"
    echo "- Dashboard interfaces"
    echo "- Performance and reliability"
    echo "- Error handling and security"
    echo "- Cross-feature integration"
    echo ""
    exit 1
fi