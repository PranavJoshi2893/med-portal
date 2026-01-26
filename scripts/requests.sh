#! /usr/bin/bash

# auth queries
curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/auth/login -H "Content-Type: application/json" -d '{"email":"pranavjoshi@gmail.com","password":"Admin@123"}'


# user queries
curl -X GET http://localhost:3000/api/v1/users/
curl -X GET http://localhost:3000/api/v1/users/019bef37-b9e5-72fe-80c9-49fed77a21d7
curl -X DELETE http://localhost:3000/api/v1/users/8497f8e9-98cf-4408-a46d-393a220bf2d4

# validation errors
curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshigmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav123","last_name":"Joshi","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav123","last_name":"Joshi123","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/auth/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi123","email":"pranavjoshi@gmail.com","password":"Admin@123"}'