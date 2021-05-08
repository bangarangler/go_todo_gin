#!/bin/zsh

curl -X DELETE http://localhost:8080/api/v1/todos/7 -v | jq # need jq installed to pretty print json or remove the | jq

