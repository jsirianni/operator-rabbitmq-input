# operator-rabbitmq-input

An operator for [Stanza](https://github.com/observIQ/stanza) and [OpenTelemetry](https://github.com/open-telemetry/opentelemetry-log-collection)

## `rabbitmq_input` operator

The `tcp_input` operator listens for logs on one or more TCP connections. The operator assumes that logs are newline separated.

### Configuration Fields

| Field             | Default          | Description                                                  |
| ---               | ---              | ---                                                          |
| `username`        |                  | Rabbitmq Username (Required)                                 |
| `password`        |                  | Rabbitmq Password (Required)                                 |
| `routing_key`     |                  | Rabbitmq Routing Key (Required)                              |
| `address`         | `localhost`      | Rabbitmq Address (IP address, hostname, or fqdn) (optional)  |
| `port`            | `5672`           | Rabbitmq client connection port (optional)                   |
| `vhost`           | `/`              | Rabbitmq virtual host (optional)                             |
| `tag`             |                  | Rabbitmq client connection tag (optional)                    |

### Output Fields

The `rabbitmq_input` operator returns a [Delivery](https://github.com/streadway/amqp/blob/1c71cc93ed716f9a6f4c2ae8955c25f9176d9f19/delivery.go#L28)
from the [amqp](github.com/streadway/amqp) package. Please see the source for in depth
details of each field.

The only field modified by `rabbitmq_input` is the `body` field. Converted from `[]byte` to `string`.
This field could be passed to `json_parser` for further parsing, or sent directly to an output operator.
Body will contain any message sent to Rabbitmq by the message producers.

| Field  | Type     | Example                | Description                      |
| ---    | ---      | ---                    | ---                              |
| body   | `string` | `"{\"key\":\"value\"}` | The body of the Rabbitmq message |

### Example Configurations

#### Simple

Configuration:
```yaml
pipeline:
- type: rabbitmq_input
  username: dev
  password: dev
  address: 127.0.0.1
  vhost: devhost
  routing_key: devkey
  tag: devel
- type: stdout
```

Send a log to the Rabbitmq REST endpoint:
```bash
curl -i -u dev:dev \
    -H "content-type:application/json" \
    -X POST "127.0.0.1:15672/api/exchanges/dev/devhost/publish" \
    -d'{"properties":{},"routing_key":"devkey","payload":"{\"key\":\"value\"}","payload_encoding":"string"}'
```

Generated entries:
```json
{
  "timestamp": "2021-02-11T19:04:05.717506625-05:00",
  "severity": 0,
  "record": {
    "acknowledger": {},
    "app_id": "",
    "body": "{\"key\":\"value\"}",
    "consumer_tag": "devel",
    "content_encoding": "",
    "content_type": "",
    "correlation_id": "",
    "delivery_mode": 0,
    "exchange": "devhost",
    "expiration": "",
    "headers": {},
    "message_id": "",
    "priority": 0,
    "redelivered": false,
    "reply_to": "",
    "type": "",
    "user_id": ""
  }
}
```

### Testing

See `test/README.md` for details
