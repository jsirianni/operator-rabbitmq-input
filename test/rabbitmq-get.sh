#!/bin/bash

ENDPOINT="http://127.0.0.1:15672"

curl -i -u dev:dev \
    -H "content-type:application/json" \
    -X POST "${ENDPOINT}/api/queues/dev/webq1/get" \
    -d'{"count":1,"encoding":"auto","truncate":50000,"ackmode":"ack_requeue_false"}'
