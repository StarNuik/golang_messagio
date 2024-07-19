#! /bin/bash
curl localhost:8080/message --request "POST" --data "{\"content\":\"$1\"}"