#! /bin/bash
curl localhost/message \
    --request "POST" \
    --data "{\"content\":\"$1\"}"