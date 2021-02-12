package main

import (
    "fmt"
    "time"
    "sync"
    "encoding/json"
    "github.com/streadway/amqp"
)


type Conn struct {
	Channel *amqp.Channel
    wg sync.WaitGroup
}

func GetConn(rabbitURL string) (Conn, error) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return Conn{}, err
	}

	ch, err := conn.Channel()
	return Conn{
		Channel: ch,
	}, err
}

func (conn Conn) Publish(data []byte) error {
	return conn.Channel.Publish(
		"webex",
		"webq1",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		})
}

func (c *Conn) Loader() {
    c.wg.Add(1)
    defer c.wg.Done()

    fmt.Println("starting worker. . . ")
    time.Sleep(time.Second * 2)

    raw := map[string]interface{}{
        "string": "value",
        "int": 100000000000000,
        "bool": true,
        "nested": map[string]string{
            "a":"b",
            "x":"Y",
        },
        "array": []string{
            "value-a",
            "value-b",
            "value-c",
        },
        "deep": map[string]interface{}{
            "top_level": map[string]interface{}{
                "name": "second_level",
                "num": 2,
                "payload": map[string]interface{}{
                    "a":"b",
                    "x": true,
                },
            },
        },
    }

    data, err := json.Marshal(raw)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    i := 0
    for {
        if i == 100000 {
            return
        }
        if err := c.Publish(data); err != nil {
            fmt.Println(err.Error())
            return
        }
        i++
    }
}

func main() {
    url := "amqp://dev:dev@localhost:5672/dev"
    c, err := GetConn(url)
    if err != nil {
        panic(err)
    }

    fmt.Println("configured connection")

    go c.Loader()
/*    go c.Loader()
    go c.Loader()
    go c.Loader()
    go c.Loader()
    go c.Loader()*/

    time.Sleep(time.Second * 3)
    c.wg.Wait()
}
