#! /usr/bin/bash

# auth queries

# register
curl -s -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshi@gmail.com","password":"Admin@123"}' | jq

# login
curl -s -X POST http://localhost:3000/api/v1/auth/login -H "Content-Type: application/json" -d '{"email":"pranavjoshi@gmail.com","password":"Admin@123"}' | jq

# token and user id stored in variables for further requests
sleep 1
TOKEN=$(curl -s -X POST http://localhost:3000/api/v1/auth/login -H "Content-Type: application/json" -d '{"email":"pranavjoshi@gmail.com","password":"Admin@123"}' | jq -r '.data.access_token')
USER_ID=$(curl -s -X GET http://localhost:3000/api/v1/users/ -H "Authorization: Bearer $TOKEN" | jq -r '.data[0].id')


# user queries
# get all users
curl -s -X GET http://localhost:3000/api/v1/users/ -H "Authorization: Bearer $TOKEN" | jq
# get user by id
curl -s -X GET http://localhost:3000/api/v1/users/$USER_ID -H "Authorization: Bearer $TOKEN" | jq


# update user by id
curl -s -X PATCH http://localhost:3000/api/v1/users/$USER_ID -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"first_name":"Pranav123","last_name":"Joshi123"}' | jq

# delete user by id
curl -s -X DELETE http://localhost:3000/api/v1/users/$USER_ID -H "Authorization: Bearer $TOKEN" | jq

# # validation errors
# curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshigmail.com","password":"Admin@123"}'
# curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav123","last_name":"Joshi","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
# curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav123","last_name":"Joshi123","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
# curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi123","email":"pranavjoshi@gmail.com","password":"Admin@123"}'