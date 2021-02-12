#!/bin/bash

ENDPOINT="http://127.0.0.1:15672"

curl -i -u dev:dev \
    -H "content-type:application/json" \
    -X POST "${ENDPOINT}/api/exchanges/dev/webex/publish" \
    -d'{"properties":{},"routing_key":"webq1","payload":"44.4","payload_encoding":"string"}'

curl -i -u dev:dev \
    -H "content-type:application/json" \
    -X POST "${ENDPOINT}/api/exchanges/dev/webex/publish" \
    -d'{"properties":{},"routing_key":"webq1","payload":"{\"key\":\"value\",\"nested\":{\"k\":\"v\",\"id\":1,\"array\":[1,2,3]}}","payload_encoding":"string"}'

