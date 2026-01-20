#! /usr/bin/bash

# user registration
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshigmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav123","last_name":"Joshi","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav123","last_name":"Joshi123","email":"pranavjoshi@gmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi123","email":"pranavjoshi@gmail.com","password":"Admin@123"}'

curl -X GET http://localhost:3000/api/v1/users/
curl -X GET http://localhost:3000/api/v1/users/dbb72380-b10b-4592-a630-09c5a055c8c7
curl -X DELETE http://localhost:3000/api/v1/users/dbb72380-b10b-4592-a630-09c5a055c8c7