#!/bin/bash


token=$(curl -X POST http://localhost:8080/auth/login -H "Content-Type:application/json" --data '{"email": "victorrbravo@gmail.com","password": "pass"}')


auth=$(echo "Authorization:$token" | tr -d '"')
echo "$auth"
AUTH=$auth
comment=$1
#params='{"comment": "'$comment'"}'
params='{"name":"'$comment'","token_link":"","pattern":"or"}'
echo $params
data=$(curl -X POST http://localhost:8080/social/add-step/a1 -H "Content-Type:application/json" --data $params -H $AUTH)

echo $data | jq .

