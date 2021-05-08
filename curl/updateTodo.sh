#!/bin/zsh

curl -X PUT http://localhost:8080/api/v1/todos/2 \
  -F "title=update from curl sh" \
  -F "completed=0" \
  -v | jq # need jq installed to pretty print json or remove the | jq

