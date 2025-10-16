#!/bin/bash

echo "Testing student assignment detail functionality..."
echo "User: kris (ID: 1)"
echo ""

echo "1. Testing student assignment ID 2 detail:"
curl -b kris_cookies.txt -s http://localhost:8080/student/assignments/2/detail

echo ""
echo ""
echo "2. Testing student assignment ID 6 detail:"
curl -b kris_cookies.txt -s http://localhost:8080/student/assignments/6/detail

echo ""
echo ""
echo "3. Testing mark assignment 2 as in progress:"
curl -b kris_cookies.txt -s -X POST -H "Content-Type: application/json" -d '{"status": "in_progress"}' http://localhost:8080/student/assignments/2/progress

echo ""
echo ""
echo "4. Testing mark assignment 6 as in progress:"
curl -b kris_cookies.txt -s -X POST -H "Content-Type: application/json" -d '{"status": "in_progress"}' http://localhost:8080/student/assignments/6/progress

echo ""
echo ""
echo "5. Testing assignment completion for assignment 2:"
curl -b kris_cookies.txt -s -X POST -H "Content-Type: application/json" -d '{"status": "completed"}' http://localhost:8080/student/assignments/2/complete

echo ""
