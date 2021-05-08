#!/bin/zsh

curl http://localhost:8080/api/v1/todos/ -v | jq # need jq installed to pretty print json or remove the | jq

