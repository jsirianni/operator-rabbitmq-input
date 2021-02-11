package rabbitmq

import (
    "fmt"
    "net/url"
)

// setURI tests the connection string
func (r *Rabbitmq) setURI(username, password, host, port, vhost string) error {
    u := url.URL{Scheme: "amqp", Path: vhost}

    if username != "" && password != "" {
        u.User = url.UserPassword(username, password)
    } else if username != "" {
        u.User = url.User(username)
    }

    if host == "" || port == "" {
        return fmt.Errorf("host and port must be set")
    }

    u.Host = host + ":" + port

    r.uri = u.String()
    return nil
}
