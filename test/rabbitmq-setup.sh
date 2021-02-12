#!/bin/bash

ENDPOINT="http://127.0.0.1:15672"

curl -i -u dev:dev \
  -H "content-type:application/json" \
  -X PUT "${ENDPOINT}/api/exchanges/dev/webex" \
  -d'{"type":"direct","auto_delete":false,"durable":true,"internal":false,"arguments":{}}'


curl -i -u dev:dev \
    -H "content-type:application/json" \
    -X PUT "${ENDPOINT}/api/queues/dev/webq1" \
    -d'{"auto_delete":false,"durable":true,"arguments":{}}'


curl -i -u dev:dev \
    -H "content-type:application/json" \
    -X POST "${ENDPOINT}/api/bindings/dev/e/webex/q/webq1" \
    -d'{"routing_key":"webq1","arguments":{}}'
