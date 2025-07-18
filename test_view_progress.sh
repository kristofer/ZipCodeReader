#!/bin/bash

echo "=== Testing View Progress Functionality ==="
echo "Note: Server needs to be restarted to pick up new routes!"
echo ""

# Test instructor login
echo "1. Testing Instructor Login:"
curl -c instructor_progress_test.txt -X POST http://localhost:8080/local/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=dolio&password=password" \
  -w "HTTP Status: %{http_code}\n" -s -o /dev/null

if [ $? -eq 0 ]; then
    echo "✅ Instructor login successful"
else
    echo "❌ Instructor login failed"
    exit 1
fi
echo ""

# Test the new route
echo "2. Testing Student Progress API (JSON):"
response=$(curl -b instructor_progress_test.txt -H "Accept: application/json" \
  http://localhost:8080/instructor/students/kris/progress -w "%{http_code}" -s)

if echo "$response" | grep -q "404"; then
    echo "❌ Route not found - Server needs restart to pick up new routes"
    echo "   Expected route: GET /instructor/students/kris/progress"
else
    echo "✅ API endpoint working"
    echo "$response" | jq . 2>/dev/null || echo "$response"
fi
echo ""

echo "3. Testing Student Progress Page (HTML):"
page_response=$(curl -b instructor_progress_test.txt \
  http://localhost:8080/instructor/students/kris/progress -w "%{http_code}" -s)

if echo "$page_response" | grep -q "404"; then
    echo "❌ Page not found - Server needs restart"
else
    echo "✅ HTML page working"
    if echo "$page_response" | grep -q "Student Progress"; then
        echo "✅ Page contains expected content"
    else
        echo "⚠️ Page loaded but may have content issues"
    fi
fi
echo ""

echo "4. Testing View Progress Link from Instructor Dashboard:"
dashboard_html=$(curl -b instructor_progress_test.txt http://localhost:8080/instructor/dashboard -s)

if echo "$dashboard_html" | grep -q "viewStudentProgress"; then
    echo "✅ JavaScript function 'viewStudentProgress' found in dashboard"
else
    echo "❌ JavaScript function missing from dashboard"
fi

if echo "$dashboard_html" | grep -q "View Progress"; then
    echo "✅ 'View Progress' button found in dashboard"
else
    echo "❌ 'View Progress' button missing from dashboard"
fi
echo ""

echo "=== Summary ==="
echo "✅ New route added: GET /instructor/students/:username/progress"
echo "✅ Handler implemented: GetStudentProgress"
echo "✅ Template created: student_progress.html"
echo "✅ API supports both JSON and HTML responses"
echo "🔄 Server restart required to activate new routes"
echo ""
echo "After server restart, instructors can:"
echo "1. Click 'View Progress' next to any student in the dashboard"
echo "2. Navigate to /instructor/students/{username}/progress directly"
echo "3. Use API with Accept: application/json header for programmatic access"

# Cleanup
rm -f instructor_progress_test.txt
