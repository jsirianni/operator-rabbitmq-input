package rabbitmq

import (
    "testing"

    "github.com/stretchr/testify/require"
)

func TestURI(t *testing.T) {
    cases := []struct{
        name string
        user string
        pass string
        host string
        port string
        vhost     string
        expected  string
        expectErr bool
    }{
        {
            "simple",
            "dev",
            "dev",
            "localhost",
            "5672",
            "",
            "amqp://dev:dev@localhost:5672",
            false,
        },
        {
            "vhost-port",
            "dev",
            "dev",
            "10.0.0.4",
            "5673",
            "stage",
            "amqp://dev:dev@10.0.0.4:5673/stage",
            false,
        },
        {
            "encoded-password",
            "dev",
            "p@ssw0rd1",
            "10.0.0.4",
            "5673",
            "stage",
            "amqp://dev:p%40ssw0rd1@10.0.0.4:5673/stage",
            false,
        },
        {
            "no-user-vhost",
            "",
            "p@ssw0rd1",
            "10.0.0.4",
            "5673",
            "",
            "amqp://10.0.0.4:5673",
            false,
        },
        {
            "bad-vhost",
            "dev",
            "dev",
            "10.0.0.4",
            "5673",
            "/stage",
            "amqp://dev:dev@10.0.0.4:5673/stage",
            false,
        },
        {
            "no-host",
            "dev",
            "dev",
            "",
            "5673",
            "dev",
            "",
            true,
        },
        {
            "no-port",
            "dev",
            "dev",
            "localhost",
            "",
            "dev",
            "",
            true,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            r := Rabbitmq{}
            err := r.setURI(tc.user, tc.pass, tc.host, tc.port, tc.vhost)
            if tc.expectErr {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            require.Equal(t, tc.expected, r.uri)
        })
    }
}
