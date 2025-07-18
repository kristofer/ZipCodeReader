#!/bin/bash

# Use existing cookies from the main test
echo "Test 17 debug:"
DUE_DATES_OVERVIEW_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/due-dates/overview)
echo "Response: '$DUE_DATES_OVERVIEW_RESPONSE'"
if echo "$DUE_DATES_OVERVIEW_RESPONSE" | grep -q "due\|dates\|overview"; then
    echo "✅ Test 17 passes"
else
    echo "❌ Test 17 fails"
fi

echo "Test 18 debug:"
DUE_DATES_NOTIFICATIONS_RESPONSE=$(curl -s -b instructor_cookies.txt http://localhost:8080/instructor/due-dates/notifications)
echo "Response: '$DUE_DATES_NOTIFICATIONS_RESPONSE'"
if echo "$DUE_DATES_NOTIFICATIONS_RESPONSE" | grep -q "notifications\|due\|dates"; then
    echo "✅ Test 18 passes"
else
    echo "❌ Test 18 fails"
fi