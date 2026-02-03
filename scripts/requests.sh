#!/usr/bin/bash

BASE_URL="http://localhost:3000/api/v1"
CREDENTIALS='{"email":"pranavjoshi@gmail.com","password":"Admin@123"}'
REGISTER_DATA='{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshi@gmail.com","password":"Admin@123"}'

# --- Register (ignores 409 if user already exists) ---
echo "=== Register ==="
curl -s -X POST "$BASE_URL/auth/register" -H "Content-Type: application/json" -d "$REGISTER_DATA" | jq

# --- Login (saves cookie for refresh/logout, captures token for user API) ---
echo "=== Login ==="
RESPONSE=$(curl -s -c cookies.txt -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d "$CREDENTIALS")

TOKEN=$(echo "$RESPONSE" | jq -r '.data.access_token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "Login failed. Response:"
  echo "$RESPONSE" | jq
  exit 1
fi

USER_ID=$(curl -s -X GET "$BASE_URL/users/" \
  -H "Authorization: Bearer $TOKEN" \
  | jq -r '.data.items[0].id')

# --- Auth (run separately if needed) ---
# curl -s -X POST "$BASE_URL/auth/register" -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshi@gmail.com","password":"Admin@123"}' | jq

# --- User queries ---
echo "=== Get all users ==="
curl -s -X GET "$BASE_URL/users/" -H "Authorization: Bearer $TOKEN" | jq

echo "=== Get user by id ==="
curl -s -X GET "$BASE_URL/users/$USER_ID" -H "Authorization: Bearer $TOKEN" | jq

echo "=== Update user ==="
curl -s -X PATCH "$BASE_URL/users/$USER_ID" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Pranav","last_name":"Joshi-Updated"}' | jq

# --- Refresh (uses cookie from login) ---
echo "=== Refresh token ==="
curl -s -X POST "$BASE_URL/auth/refresh" -b cookies.txt | jq

# --- Logout (uses cookie) ---
echo "=== Logout ==="
curl -s -X POST "$BASE_URL/auth/logout" -b cookies.txt | jq

# --- Delete (comment out to keep user for next run) ---
echo "=== Delete user ==="
curl -s -X DELETE "$BASE_URL/users/$USER_ID" -H "Authorization: Bearer $TOKEN" | jq
