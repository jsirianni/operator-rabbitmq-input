package rabbitmq

import (
    "time"
    "context"

    "github.com/streadway/amqp"
    "github.com/jpillora/backoff"
    "go.uber.org/zap"
)

// connectionManager will keep the connection to Rabbitmq open.
func (g *QueueInput) connectionManager(ctx context.Context) {
    retry := &backoff.Backoff{
        Min:    500 * time.Millisecond,
        Max:    10 * time.Second,
        Factor: 2,
    }

    for {
        select {
        case <-ctx.Done():
            return
        default:
        }

        if g.rabbit.conn == nil || g.rabbit.conn.IsClosed() {
            g.rabbit.connecting.Add(1)
            for {
                select {
                case <-ctx.Done():
                    return
                default:
                }

                if err := g.rabbit.connect(); err != nil {
                    g.Errorf("Error connecting to Rabbitmq", zap.Error(err))
                    time.Sleep(retry.Duration())
                    continue
                }
                break
            }
            retry.Reset()
            g.Infow("Rabbitmq connection successful")
            g.rabbit.connecting.Done()
        }

        time.Sleep(time.Second * 2)
    }
}

// connect will connect to Rabbitmq.
func (r *Rabbitmq) connect() error {
    var err error

    r.conn, err = amqp.Dial(r.uri)
    if err != nil {
        return err
    }

    r.channel, err = r.conn.Channel()
    if err != nil {
        return err
    }

    count := r.workerCount * 4
    if err := r.channel.Qos(count, 0, false); err != nil {
        return err
    }

    r.deliveries, err = r.channel.Consume(
        // https://github.com/streadway/amqp/blob/v1.0.0/channel.go#L1052
        r.routingKey,  // name
        r.tag,         // consumerTag,
        false,         // auto ack
        false,         // exclusive
        false,         // noLocal
        false,         // noWait
        nil,           // arguments
    )
    return err
}

// closeConnection will close the connection to Rabbitmq.
func (r *Rabbitmq) close() error {
    if r.conn == nil || r.conn.IsClosed() {
        return nil
    }
    return r.conn.Close()
}
