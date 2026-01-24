#! /usr/bin/bash

# user registration
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshigmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav123","last_name":"Joshi","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav123","last_name":"Joshi123","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi123","email":"pranavjoshi@gmail.com","password":"Admin@123"}'

curl -X GET http://localhost:3000/api/v1/users/
curl -X GET http://localhost:3000/api/v1/users/8497f8e9-98cf-4408-a46d-393a220bf2d4
curl -X DELETE http://localhost:3000/api/v1/users/8497f8e9-98cf-4408-a46d-393a220bf2d4
curl -X POST http://localhost:3000/api/v1/users/login -H "Content-Type: application/json" -d '{"email":"pranavjoshi@gmail.com","password":"Admin@123"}'