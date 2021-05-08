#!/bin/zsh

curl POST http://localhost:8080/api/v1/todos/ \
  -F "title=new todo" \
  -F "completed=false" \
  -v | jq # need jq installed to pretty print json or remove the | jq

