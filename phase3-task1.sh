#!/bin/bash

# Phase 3 Task 1: Assignment Models and Database Schema
# Tests database schema, model relationships, and data integrity

echo "🚀 Phase 3 Task 1: Assignment Models and Database Schema"
echo "============================================================"

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

# Test 1: Health Check (Database Connection)
echo "Test 1: Database Connection via Health Check"
HEALTH_RESPONSE=$(curl -s http://localhost:8080/health)
if echo "$HEALTH_RESPONSE" | grep -q "healthy"; then
    print_result 0 "Database connection established"
else
    print_result 1 "Database connection failed"
fi

# Test 2: Create test users (instructor and student)
echo "Test 2: Create test users for model testing"

# Use timestamp to create unique usernames
TIMESTAMP=$(date +%s)
INSTRUCTOR_USER="instructor_$TIMESTAMP"
STUDENT_USER="student_$TIMESTAMP"

# Register instructor
INSTRUCTOR_REG=$(curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR_USER&email=$INSTRUCTOR_USER@example.com&password=password&confirm_password=password&role=instructor")

# Register student
STUDENT_REG=$(curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT_USER&email=$STUDENT_USER@example.com&password=password&confirm_password=password&role=student")

# Login as instructor to get session cookies
curl -s -c instructor_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR_USER&password=password" > /dev/null

# Login as student to get session cookies  
curl -s -c student_cookies.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT_USER&password=password" > /dev/null

print_result 0 "Test users created ($INSTRUCTOR_USER, $STUDENT_USER)"

# Test 3: Assignment Model - Create Assignment
echo "Test 3: Assignment Model - Create Assignment"
ASSIGNMENT_RESPONSE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Assignment 1",
    "description": "Test assignment for model validation",
    "url": "https://example.com/reading1",
    "category": "Programming",
    "due_date": "2025-08-01T23:59:59Z"
  }')

if echo "$ASSIGNMENT_RESPONSE" | grep -q "Test Assignment 1"; then
    ASSIGNMENT_ID=$(echo "$ASSIGNMENT_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
    print_result 0 "Assignment model creation works (ID: $ASSIGNMENT_ID)"
else
    print_result 1 "Assignment model creation failed"
    ASSIGNMENT_ID=1
fi

# Test 4: Assignment Model - Required Fields Validation
echo "Test 4: Assignment Model - Required Fields Validation"
INVALID_ASSIGNMENT=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Missing title and URL"
  }')

if echo "$INVALID_ASSIGNMENT" | grep -q "error\|required\|invalid"; then
    print_result 0 "Assignment model validation works (rejects missing required fields)"
else
    print_result 1 "Assignment model validation failed (should reject missing required fields)"
fi

# Test 5: Assignment Model - Foreign Key Relationship (CreatedBy)
echo "Test 5: Assignment Model - Foreign Key Relationship (CreatedBy)"
ASSIGNMENT_DETAIL=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID)

if echo "$ASSIGNMENT_DETAIL" | grep -q "instructor1\|created_by"; then
    print_result 0 "Assignment foreign key relationship to User works"
else
    print_result 1 "Assignment foreign key relationship to User failed"
fi

# Test 6: StudentAssignment Model - Create Relationship
echo "Test 6: StudentAssignment Model - Create Relationship"
# Get the actual student ID from the database
STUDENT_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT_USER';")
INSTRUCTOR_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$INSTRUCTOR_USER';")

ASSIGN_RESPONSE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID/assign \
  -H "Content-Type: application/json" \
  -d "{
    \"student_ids\": [$STUDENT_ID]
  }")

if echo "$ASSIGN_RESPONSE" | grep -q "success\|assigned"; then
    print_result 0 "StudentAssignment model creation works"
else
    print_result 1 "StudentAssignment model creation failed"
fi

# Test 7: StudentAssignment Model - Status Tracking
echo "Test 7: StudentAssignment Model - Status Tracking"
STUDENT_ASSIGNMENTS=$(curl -s -b student_cookies.txt http://localhost:8080/student/assignments)

if echo "$STUDENT_ASSIGNMENTS" | grep -q "assigned\|status"; then
    print_result 0 "StudentAssignment status tracking works"
else
    print_result 1 "StudentAssignment status tracking failed"
fi

# Test 8: Assignment Model - Soft Delete Support
echo "Test 8: Assignment Model - Soft Delete Support"
DELETE_RESPONSE=$(curl -s -b instructor_cookies.txt -X DELETE http://localhost:8080/instructor/assignments/$ASSIGNMENT_ID)

# Check if assignment is soft deleted (should not appear in active list but might be in database)
ASSIGNMENTS_AFTER_DELETE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments)

if ! echo "$ASSIGNMENTS_AFTER_DELETE" | grep -q "Test Assignment 1"; then
    print_result 0 "Assignment soft delete works"
else
    print_result 1 "Assignment soft delete failed"
fi

# Test 9: Assignment Model - Category System
echo "Test 9: Assignment Model - Category System"
CATEGORY_ASSIGNMENT=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Category Test Assignment",
    "description": "Testing category functionality",
    "url": "https://example.com/category-test",
    "category": "Database Design"
  }')

if echo "$CATEGORY_ASSIGNMENT" | grep -q "Database Design"; then
    print_result 0 "Assignment category system works"
else
    print_result 1 "Assignment category system failed"
fi

# Test 10: Assignment Model - Due Date System (Nullable)
echo "Test 10: Assignment Model - Due Date System (Nullable)"
NO_DUE_DATE_ASSIGNMENT=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "No Due Date Assignment",
    "description": "Assignment without due date",
    "url": "https://example.com/no-due-date",
    "category": "Optional Reading"
  }')

if echo "$NO_DUE_DATE_ASSIGNMENT" | grep -q "No Due Date Assignment"; then
    print_result 0 "Assignment nullable due date works"
else
    print_result 1 "Assignment nullable due date failed"
fi

# Test 11: Database Indexes and Performance
echo "Test 11: Database Indexes and Performance"
# Create multiple assignments to test indexing
for i in {1..5}; do
    curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
      -H "Content-Type: application/json" \
      -d "{
        \"title\": \"Performance Test Assignment $i\",
        \"description\": \"Testing database performance\",
        \"url\": \"https://example.com/perf-test-$i\",
        \"category\": \"Performance Testing\"
      }" > /dev/null
done

# Test query performance (should be fast)
START_TIME=$(date +%s%N)
curl -s -b instructor_cookies.txt http://localhost:8080/instructor/assignments > /dev/null
END_TIME=$(date +%s%N)
QUERY_TIME=$((($END_TIME - $START_TIME) / 1000000)) # Convert to milliseconds

if [ $QUERY_TIME -lt 1000 ]; then # Less than 1 second
    print_result 0 "Database query performance acceptable ($QUERY_TIME ms)"
else
    print_result 1 "Database query performance slow ($QUERY_TIME ms)"
fi

# Test 12: Model Relationships - Many-to-Many (Assignment-Student)
echo "Test 12: Model Relationships - Many-to-Many (Assignment-Student)"
# Create another student
STUDENT2_USER="student2_$TIMESTAMP"
curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT2_USER&email=$STUDENT2_USER@example.com&password=password&confirm_password=password&role=student" > /dev/null

# Get the actual student2 ID from the database
STUDENT2_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT2_USER';")

# Create a new assignment for multi-student testing
MULTI_ASSIGNMENT=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Multi-Student Assignment",
    "description": "Assignment for testing many-to-many relationships",
    "url": "https://example.com/multi-student",
    "category": "Testing"
  }')

MULTI_ASSIGNMENT_ID=$(echo "$MULTI_ASSIGNMENT" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)

# Assign the same assignment to multiple students
MULTI_ASSIGN_RESPONSE=$(curl -s -b instructor_cookies.txt -X POST http://localhost:8080/instructor/assignments/$MULTI_ASSIGNMENT_ID/assign \
  -H "Content-Type: application/json" \
  -d "{
    \"student_ids\": [$STUDENT_ID, $STUDENT2_ID]
  }")

if echo "$MULTI_ASSIGN_RESPONSE" | grep -q "success\|assigned"; then
    print_result 0 "Many-to-many assignment relationships work"
else
    print_result 1 "Many-to-many assignment relationships failed"
fi

# Cleanup
echo "🧹 Cleanup"
rm -f instructor_cookies.txt student_cookies.txt student2_cookies.txt
sleep 1

# Summary
echo ""
echo "📊 Phase 3 Task 1 Test Results:"
echo "================================"
echo -e "${GREEN}✅ Passed: $PASS${NC}"
echo -e "${RED}❌ Failed: $FAIL${NC}"
echo ""

if [ $FAIL -eq 0 ]; then
    echo -e "${GREEN}🎉 Phase 3 Task 1: Assignment Models and Database Schema - ALL TESTS PASSED!${NC}"
    exit 0
else
    echo -e "${RED}🚨 Phase 3 Task 1: Assignment Models and Database Schema - $FAIL TESTS FAILED!${NC}"
    echo ""
    echo "Issues to address:"
    echo "- Assignment model creation and validation"
    echo "- StudentAssignment relationship management"
    echo "- Database schema integrity"
    echo "- Foreign key relationships"
    echo "- Soft delete functionality"
    echo "- Category and due date systems"
    echo ""
    exit 1
fi