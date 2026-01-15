#! /usr/bin/bash

curl -X POST http://localhost:3000/api/v1/users/register -H "Content-Type: application/json" -d '{"first_name":"Anant","last_name":"Joshi","email":"anantjoshi1753@gmail.com","password":"Admin@123"}'
curl -X GET http://localhost:3000/api/v1/users/