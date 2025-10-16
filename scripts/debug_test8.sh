#!/bin/bash

# Debug Test 8: Progress visualization data
echo "Debug Test 8: Progress visualization data"

# Setup like the main test
TIMESTAMP=$(date +%s)
INSTRUCTOR_USER="instructor_$TIMESTAMP"

# Register and login
curl -s -X POST http://localhost:8080/local/register \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR_USER&email=$INSTRUCTOR_USER@example.com&password=password&confirm_password=password&role=instructor" > /dev/null

curl -s -c instructor_debug.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$INSTRUCTOR_USER&password=password" > /dev/null

# Test the endpoint
echo "Testing /instructor/progress/trends endpoint:"
VISUALIZATION_DATA=$(curl -s -b instructor_debug.txt http://localhost:8080/instructor/progress/trends)
echo "Response: '$VISUALIZATION_DATA'"

# Test the search patterns
echo "Checking for 'data':"
echo "$VISUALIZATION_DATA" | grep -q "data" && echo "Found 'data'" || echo "Not found 'data'"

echo "Checking for 'visualization':"
echo "$VISUALIZATION_DATA" | grep -q "visualization" && echo "Found 'visualization'" || echo "Not found 'visualization'"

echo "Checking for 'chart':"
echo "$VISUALIZATION_DATA" | grep -q "chart" && echo "Found 'chart'" || echo "Not found 'chart'"

# Cleanup
rm -f instructor_debug.txt