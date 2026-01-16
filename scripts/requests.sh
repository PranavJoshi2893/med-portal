#! /usr/bin/bash

# user registration
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshi2893@gmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi","email":"pranavjoshigmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav123","last_name":"Joshi","email":"pranavjoshi2893@gmail.com","password":"Admin@123"}'
curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Pranav","last_name":"Joshi123","email":"pranavjoshi2893@gmail.com","password":"Admin@123"}'

curl -X GET http://localhost:3000/api/v1/users/
curl -X DELETE http://localhost:3000/api/v1/users/ed16c5ea-599e-4cc5-97fd-d80f9d727dc1