#! /usr/bin/bash

# auth queries
curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshi1234@gmail.com","password":"Admin@123"}' | jq
# curl -X POST http://localhost:3000/api/v1/auth/login -H "Content-Type: application/json" -d '{"email":"pranavjoshi@gmail.com","password":"Admin@123"}' | jq
TOKEN=$(curl -X POST http://localhost:3000/api/v1/auth/login -H "Content-Type: application/json" -d '{"email":"pranavjoshi1234@gmail.com","password":"Admin@123"}' | jq -r '.data.access_token')


# # user queries
curl -X GET http://localhost:3000/api/v1/users/ -H "Authorization: Bearer $TOKEN" | jq
USER_ID=$(curl -X GET http://localhost:3000/api/v1/users/ -H "Authorization: Bearer $TOKEN" | jq -r '.data[0].id')
curl -X GET http://localhost:3000/api/v1/users/$USER_ID -H "Authorization: Bearer $TOKEN" | jq
# curl -X DELETE http://localhost:3000/api/v1/users/$USER_ID -H "Authorization: Bearer $TOKEN" | jq


# # validation errors
# curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshigmail.com","password":"Admin@123"}'
# curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav123","last_name":"Joshi","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
# curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav123","last_name":"Joshi123","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
# curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi123","email":"pranavjoshi@gmail.com","password":"Admin@123"}'