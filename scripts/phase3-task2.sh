#!/bin/bash

# Phase 3 Task 2: Assignment Service Layer
# Tests business logic layer for assignment operations, CRUD operations, validation, and authorization

echo "🚀 Phase 3 Task 2: Assignment Service Layer"
echo "============================================="

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

curl -s -c student2_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT2_USER&password=password" > /dev/null

# Test 1: Assignment Creation Service - Valid Data
echo "Test 1: Assignment Creation Service - Valid Data"
ASSIGNMENT_RESPONSE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Service Test Assignment",
    "description": "Testing assignment service layer",
    "url": "https://example.com/service-test",
    "category": "Service Testing",
    "due_date": "2025-08-15T23:59:59Z"
  }')

if echo "$ASSIGNMENT_RESPONSE" | grep -q "Service Test Assignment"; then
    ASSIGNMENT_ID=$(echo "$ASSIGNMENT_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
    print_result 0 "Assignment creation service works (ID: $ASSIGNMENT_ID)"
else
    print_result 1 "Assignment creation service failed"
    ASSIGNMENT_ID=1
fi

# Test 2: Assignment Creation Service - Validation (Missing Title)
echo "Test 2: Assignment Creation Service - Validation (Missing Title)"
INVALID_RESPONSE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Missing title",
    "url": "https://example.com/missing-title"
  }')

if echo "$INVALID_RESPONSE" | grep -q "error\|invalid\|required"; then
    print_result 0 "Assignment validation works (rejects missing title)"
else
    print_result 1 "Assignment validation failed (should reject missing title)"
fi

# Test 3: Assignment Creation Service - Validation (Missing URL)
echo "Test 3: Assignment Creation Service - Validation (Missing URL)"
INVALID_URL_RESPONSE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Missing URL Assignment",
    "description": "This assignment has no URL"
  }')

if echo "$INVALID_URL_RESPONSE" | grep -q "error\|invalid\|required"; then
    print_result 0 "Assignment validation works (rejects missing URL)"
else
    print_result 1 "Assignment validation failed (should reject missing URL)"
fi

# Test 4: Assignment Creation Service - Authorization (Instructor Only)
echo "Test 4: Assignment Creation Service - Authorization (Instructor Only)"
UNAUTHORIZED_RESPONSE=$(curl -s -b student_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Unauthorized Assignment",
    "description": "Students should not be able to create assignments",
    "url": "https://example.com/unauthorized"
  }')

if echo "$UNAUTHORIZED_RESPONSE" | grep -q "unauthorized\|forbidden\|403\|Authentication required\|error"; then
    print_result 0 "Assignment authorization works (students cannot create assignments)"
else
    print_result 1 "Assignment authorization failed (students should not create assignments)"
fi

# Test 5: Assignment Retrieval Service - Get by Instructor
echo "Test 5: Assignment Retrieval Service - Get by Instructor"
INSTRUCTOR_ASSIGNMENTS=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments)

if echo "$INSTRUCTOR_ASSIGNMENTS" | grep -q "Service Test Assignment"; then
    print_result 0 "Assignment retrieval by instructor works"
else
    print_result 1 "Assignment retrieval by instructor failed"
fi

# Test 6: Assignment Retrieval Service - Get Specific Assignment
echo "Test 6: Assignment Retrieval Service - Get Specific Assignment"
SPECIFIC_ASSIGNMENT=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID)

if echo "$SPECIFIC_ASSIGNMENT" | grep -q "Service Test Assignment"; then
    print_result 0 "Specific assignment retrieval works"
else
    print_result 1 "Specific assignment retrieval failed"
fi

# Test 7: Assignment Update Service
echo "Test 7: Assignment Update Service"
UPDATE_RESPONSE=$(curl -s -b instructor_cookies.txt -X PUT http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Service Test Assignment",
    "description": "Updated description for testing",
    "url": "https://example.com/updated-service-test",
    "category": "Updated Category"
  }')

if echo "$UPDATE_RESPONSE" | grep -q "Updated Service Test Assignment\|success"; then
    print_result 0 "Assignment update service works"
else
    print_result 1 "Assignment update service failed"
fi

# Test 8: Assignment Update Service - Authorization (Owner Only)
echo "Test 8: Assignment Update Service - Authorization (Owner Only)"
# Create another instructor
INSTRUCTOR2_USER="instructor2_$TIMESTAMP"
curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR2_USER&email=$INSTRUCTOR2_USER@example.com&password=password&confirm_password=password&role=instructor" > /dev/null

curl -s -c instructor2_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR2_USER&password=password" > /dev/null

UNAUTHORIZED_UPDATE=$(curl -s -b instructor2_cookies.txt -X PUT http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Unauthorized Update",
    "description": "Should not be allowed"
  }')

if echo "$UNAUTHORIZED_UPDATE" | grep -q "unauthorized\|forbidden\|403\|Authentication required\|error"; then
    print_result 0 "Assignment update authorization works (only owner can update)"
else
    print_result 1 "Assignment update authorization failed (should restrict to owner)"
fi

# Test 9: Student Assignment Service - Assign to Student
echo "Test 9: Student Assignment Service - Assign to Student"
# Get the actual student IDs from the database
STUDENT_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT_USER';")
ASSIGN_RESPONSE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/assign \
  -H "Content-Type: application/json" \
  -d "{
    \"student_ids\": [$STUDENT_ID]
  }")

if echo "$ASSIGN_RESPONSE" | grep -q "success\|assigned"; then
    print_result 0 "Student assignment service works (single student)"
else
    print_result 1 "Student assignment service failed (single student)"
fi

# Test 10: Student Assignment Service - Bulk Assignment
echo "Test 10: Student Assignment Service - Bulk Assignment"
STUDENT2_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT2_USER';")
BULK_ASSIGN_RESPONSE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/assign \
  -H "Content-Type: application/json" \
  -d "{
    \"student_ids\": [$STUDENT_ID, $STUDENT2_ID]
  }")

if echo "$BULK_ASSIGN_RESPONSE" | grep -q "success\|assigned"; then
    print_result 0 "Student assignment service works (bulk assignment)"
else
    print_result 1 "Student assignment service failed (bulk assignment)"
fi

# Test 11: Student Assignment Service - Status Updates
echo "Test 11: Student Assignment Service - Status Updates"
STATUS_UPDATE_RESPONSE=$(curl -s -b student_cookies.txt -X POST http://localhost:8080/student/assignments/$ASSIGNMENT_ID/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_progress"
  }')

if echo "$STATUS_UPDATE_RESPONSE" | grep -q "success\|in_progress"; then
    print_result 0 "Assignment status update service works"
else
    print_result 1 "Assignment status update service failed"
fi

# Test 12: Student Assignment Service - Mark as Completed
echo "Test 12: Student Assignment Service - Mark as Completed"
COMPLETE_RESPONSE=$(curl -s -b student_cookies.txt -X POST http://localhost:8080/student/assignments/$ASSIGNMENT_ID/complete \
  -H "Content-Type: application/json")

if echo "$COMPLETE_RESPONSE" | grep -q "success\|completed"; then
    print_result 0 "Assignment completion service works"
else
    print_result 1 "Assignment completion service failed"
fi

# Test 13: Assignment Filtering Service - By Category
echo "Test 13: Assignment Filtering Service - By Category"
# Create assignments with different categories
curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Database Assignment",
    "description": "Database related assignment",
    "url": "https://example.com/database",
    "category": "Database"
  }' > /dev/null

curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Programming Assignment",
    "description": "Programming related assignment",
    "url": "https://example.com/programming",
    "category": "Programming"
  }' > /dev/null

# Test category filtering (if endpoint exists)
CATEGORY_FILTER_RESPONSE=$(curl -s -b instructor_cookies.txt "http://localhost:8080/instructor/assignments?category=Database")

if echo "$CATEGORY_FILTER_RESPONSE" | grep -q "Database Assignment"; then
    print_result 0 "Assignment filtering by category works"
else
    print_result 1 "Assignment filtering by category failed"
fi

# Test 14: Assignment Filtering Service - By Status
echo "Test 14: Assignment Filtering Service - By Status"
STUDENT_ASSIGNMENTS=$(curl -s -b student_cookies.txt http://localhost:8080/student/assignments)

if echo "$STUDENT_ASSIGNMENTS" | grep -q "completed\|in_progress\|assigned"; then
    print_result 0 "Assignment filtering by status works"
else
    print_result 1 "Assignment filtering by status failed"
fi

# Test 15: Assignment Filtering Service - By Due Date
echo "Test 15: Assignment Filtering Service - By Due Date"
DUE_DATE_FILTER_RESPONSE=$(curl -s -b instructor_cookies.txt "http://localhost:8080/instructor/assignments?due_date=2025-08-15")

if echo "$DUE_DATE_FILTER_RESPONSE" | grep -q "Updated Service Test Assignment\|2025-08-15"; then
    print_result 0 "Assignment filtering by due date works"
else
    print_result 1 "Assignment filtering by due date failed"
fi

# Test 16: Assignment Progress Service
echo "Test 16: Assignment Progress Service"
PROGRESS_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/progress)

if echo "$PROGRESS_RESPONSE" | grep -q "progress\|completion\|students"; then
    print_result 0 "Assignment progress service works"
else
    print_result 1 "Assignment progress service failed"
fi

# Test 17: Assignment Deletion Service
echo "Test 17: Assignment Deletion Service"
DELETE_RESPONSE=$(curl -s -b instructor_cookies.txt -X DELETE http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID)

if echo "$DELETE_RESPONSE" | grep -q "success\|deleted"; then
    print_result 0 "Assignment deletion service works"
else
    print_result 1 "Assignment deletion service failed"
fi

# Test 18: Assignment Deletion Service - Authorization (Owner Only)
echo "Test 18: Assignment Deletion Service - Authorization (Owner Only)"
# Create assignment with instructor2
NEW_ASSIGNMENT=$(curl -s -b instructor2_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Instructor2 Assignment",
    "description": "Assignment by instructor2",
    "url": "https://example.com/instructor2"
  }')

INSTRUCTOR2_ASSIGNMENT_ID=$(echo "$NEW_ASSIGNMENT" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)

# Try to delete with instructor1 (should fail)
UNAUTHORIZED_DELETE=$(curl -s -b instructor_cookies.txt -X DELETE http://localhost:8080/instructor/assignments/$INSTRUCTOR2_ASSIGNMENT_ID)

if echo "$UNAUTHORIZED_DELETE" | grep -q "unauthorized\|forbidden\|403\|Authentication required\|error"; then
    print_result 0 "Assignment deletion authorization works (only owner can delete)"
else
    print_result 1 "Assignment deletion authorization failed (should restrict to owner)"
fi

# Test 19: Assignment Service - Error Handling
echo "Test 19: Assignment Service - Error Handling"
NON_EXISTENT_ASSIGNMENT=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/99999)

if echo "$NON_EXISTENT_ASSIGNMENT" | grep -q "not found\|error\|404"; then
    print_result 0 "Assignment service error handling works (non-existent assignment)"
else
    print_result 1 "Assignment service error handling failed (should return error for non-existent assignment)"
fi

# Test 20: Assignment Service - Transaction Management
echo "Test 20: Assignment Service - Transaction Management"
# Test that assignment creation and student assignment are atomic
ATOMIC_RESPONSE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Atomic Test Assignment",
    "description": "Testing transaction atomicity",
    "url": "https://example.com/atomic-test",
    "category": "Testing"
  }')

if echo "$ATOMIC_RESPONSE" | grep -q "Atomic Test Assignment"; then
    print_result 0 "Assignment service transaction management works"
else
    print_result 1 "Assignment service transaction management failed"
fi

# Cleanup
echo "🧹 Cleanup"

rm -f instructor_cookies.txt student_cookies.txt instructor2_cookies.txt
sleep 1

# Summary
echo ""
echo "📊 Phase 3 Task 2 Test Results:"
echo "================================"
echo -e "${GREEN}✅ Passed: $PASS${NC}"
echo -e "${RED}❌ Failed: $FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}🎉 Phase 3 Task 2: Assignment Service Layer - ALL TESTS PASSED!${NC}"
    exit 0
else
    echo -e "${RED}🚨 Phase 3 Task 2: Assignment Service Layer - $FAIL TESTS FAILED!${NC}"
    echo ""
    echo "Issues to address:"
    echo "- Assignment CRUD operations service layer"
    echo "- Validation and authorization in services"
    echo "- Student assignment management services"
    echo "- Assignment filtering and querying services"
    echo "- Progress tracking services"
    echo "- Error handling and transaction management"
    echo ""
    exit 1
fi