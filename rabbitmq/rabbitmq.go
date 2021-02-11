package rabbitmq

import (
    "sync"

    "github.com/streadway/amqp"
)

// newRabbitmq will return a new Rabbitmq configuration.
func newRabbitmq(user, pass, addr, port, vhost, routingKey, tag string, workerCount int) (Rabbitmq, error) {
    r := Rabbitmq{
        routingKey: routingKey,
        conn: nil,
        channel: nil,
        tag: tag,
        workerCount: workerCount,
    }
    err := r.setURI(user, pass, addr, port, vhost)
    return r, err
}

// Rabbitmq represents a Rabbitmq configuration.
type Rabbitmq struct {
    // Rabbitmq connection URI
    uri string

    // Conn is the connection to Rabbitmq
    conn    *amqp.Connection

    // Channel will return messages from Rabbitmq
    channel *amqp.Channel

    // Rabbitmq consumer parameters passed to Rabbitmq.channel.Consume()
    routingKey string
    tag        string

    // Delivery is returned by Rabbitmq.channel.Consume()
    deliveries <-chan amqp.Delivery

    // WaitGroup used for Rabbitmq consumer worker go routines
    wg     sync.WaitGroup

    // Mutex is locked when connectionManager is attempting to
    // establish a connection to Rabbitmq
    connecting sync.WaitGroup

    // number of go routines for reading messages from the queue
    workerCount int
}
