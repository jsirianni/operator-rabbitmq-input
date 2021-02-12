package rabbitmq

import (
    "time"

    "github.com/observiq/stanza/entry"

    "github.com/streadway/amqp"
)

// parseMessage will parse a rabbitmq message as an Entry
func parseMessage(msg amqp.Delivery) *entry.Entry {
    m := make(map[string]interface{})
    m["acknowledger"]     = msg.Acknowledger
    m["content_type"]     = msg.ContentType
    m["content_encoding"] = msg.ContentEncoding
    m["delivery_mode"]    = msg.DeliveryMode
    m["priority"]         = msg.Priority
    m["correlation_id"]   = msg.CorrelationId
    m["reply_to"]         = msg.ReplyTo
    m["expiration"]       = msg.Expiration
    m["message_id"]       = msg.MessageId
    m["type"]             = msg.Type
    m["user_id"]          = msg.UserId
    m["app_id"]           = msg.AppId
    m["consumer_tag"]     = msg.ConsumerTag
    m["redelivered"]      = msg.Redelivered
    m["exchange"]         = msg.Exchange

    // msg.Headers is type amqp.Table{} which is an alias for map[string]interface{}
    // https://github.com/streadway/amqp/blob/e6b33f460591b0acb2f13b04ef9cf493720ffe17/types.go#L225
    m["headers"] = amqp.Table{}
    if len(msg.Headers) > 0 {
        m["headers"] = msg.Headers
    }

    // Convert body to a string. If the body is json, the user should parse
    // it using the json_parser operator. The overhead of checking if body is
    // json and then parsing it is too expensive for high throughput workloads
    // that do not consistently pass json in the body
    m["body"] = ""
    if msg.Body != nil {
        m["body"] = string(msg.Body)
    }

    e := entry.New()
    e.Record = m

    // Rabbitmq messages will have a zero time (0001-01-01T00:00:00Z)
    // if the publisher does not set the timestamp. Promote the timestamp
    // if it is set, otherwise drop it completely.
    if validTimestamp(msg.Timestamp) {
        e.Timestamp = msg.Timestamp
    }

    return e
}

// validTimestamp is true if a timestamp is valid
func validTimestamp(t time.Time) bool {
    return ! t.IsZero()
}
