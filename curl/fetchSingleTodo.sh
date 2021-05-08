#!/bin/zsh

curl http://localhost:8080/api/v1/todos/2 -v | jq # need jq installed to pretty print json or remove the | jq

