#!/bin/bash

echo "Debug: Testing student isolation after phase3-task4.sh creates users and assignments"

# Create test users like the main script does
TIMESTAMP=$(date +%s)
INSTRUCTOR_USER="instructor_$TIMESTAMP"
STUDENT_USER="student_$TIMESTAMP"
STUDENT2_USER="student2_$TIMESTAMP"

echo "Creating users: $INSTRUCTOR_USER, $STUDENT_USER, $STUDENT2_USER"

# Register users
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
curl -s -c instructor_debug.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR_USER&password=password" > /dev/null

curl -s -c student1_debug.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT_USER&password=password" > /dev/null

curl -s -c student2_debug.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$STUDENT2_USER&password=password" > /dev/null

# Create assignments
echo "Creating assignments..."
ASSIGNMENT1=$(curl -s -b instructor_debug.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Student Test Assignment 1",
    "description": "Assignment for student testing",
    "url": "https://example.com/assignment1",
    "category": "Student Testing"
  }')

ASSIGNMENT2=$(curl -s -b instructor_debug.txt -X POST http://localhost:8080/instructor/assignments \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Student Test Assignment 2",
    "description": "Assignment for student testing",
    "url": "https://example.com/assignment2",
    "category": "Student Testing"
  }')

ASSIGNMENT1_ID=$(echo "$ASSIGNMENT1" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
ASSIGNMENT2_ID=$(echo "$ASSIGNMENT2" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)

echo "Created assignments: $ASSIGNMENT1_ID, $ASSIGNMENT2_ID"

# Get student IDs
STUDENT_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT_USER';")
STUDENT2_ID=$(sqlite3 zipcodereader.db "SELECT id FROM users WHERE username = '$STUDENT2_USER';")

echo "Student IDs: $STUDENT_ID, $STUDENT2_ID"

# Assign assignments
echo "Assigning Assignment1 to both students..."
curl -s -b instructor_debug.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT1_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"student_ids\": [$STUDENT_ID, $STUDENT2_ID]}" > /dev/null

echo "Assigning Assignment2 to only Student1..."
curl -s -b instructor_debug.txt -X POST http://localhost:8080/instructor/assignments/$ASSIGNMENT2_ID/assign \
  -H "Content-Type: application/json" \
  -d "{\"student_ids\": [$STUDENT_ID]}" > /dev/null

# Check what each student sees
echo "Student1 assignments:"
STUDENT1_ASSIGNMENTS=$(curl -s -b student1_debug.txt http://localhost:8080/student/assignments)
echo "$STUDENT1_ASSIGNMENTS" | head -5
STUDENT1_COUNT=$(echo "$STUDENT1_ASSIGNMENTS" | grep -c "Student Test Assignment" 2>/dev/null || echo "0")
echo "Student1 count: $STUDENT1_COUNT"

echo "Student2 assignments:"
STUDENT2_ASSIGNMENTS=$(curl -s -b student2_debug.txt http://localhost:8080/student/assignments)
echo "$STUDENT2_ASSIGNMENTS" | head -5
STUDENT2_COUNT=$(echo "$STUDENT2_ASSIGNMENTS" | grep -c "Student Test Assignment" 2>/dev/null || echo "0")
echo "Student2 count: $STUDENT2_COUNT"

# Cleanup
rm -f instructor_debug.txt student1_debug.txt student2_debug.txt