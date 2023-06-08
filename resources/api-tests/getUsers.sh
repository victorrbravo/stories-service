#!/bin/bash


token=$(curl -X POST http://localhost:8080/auth/login -H "Content-Type:application/json" --data '{"email": "victorrbravo@gmail.com","password": "pass"}')


auth=$(echo "Authorization:$token" | tr -d '"')
echo "$auth"
export AUTH=$auth
data=$(curl -X GET -H "Content-Type:application/json" -H $auth http://localhost:8080/users/)

echo $data | jq .

