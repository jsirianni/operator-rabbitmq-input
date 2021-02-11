package rabbitmq

import (
    "time"
    "strconv"
    "testing"
    "encoding/json"

    "github.com/open-telemetry/opentelemetry-log-collection/entry"

    "github.com/streadway/amqp"
    "github.com/stretchr/testify/require"
)

// see init()
var timeNow time.Time
var timeZero time.Time
var testBenchParseMessage amqp.Delivery


func TestParseMessage(t *testing.T) {
	cases := []struct {
		name           string
		inputRecord    amqp.Delivery
		expectedRecord *entry.Entry
	}{
		{
			"timestamp",
			amqp.Delivery{
                Timestamp: timeNow,
            },
			&entry.Entry{
                Timestamp: timeNow,
                Record: map[string]interface{}{
                    "acknowledger": nil,
                    "content_type": "",
                    "content_encoding": "",
                    "delivery_mode": uint8(0),
                    "priority": uint8(0),
                    "correlation_id": "",
                    "reply_to": "",
                    "message_id": "",
                    "expiration": "",
                    "type": "",
                    "user_id": "",
                    "app_id": "",
                    "consumer_tag": "",
                    "redelivered": false,
                    "exchange": "",
                    "headers": amqp.Table{},
                    "body": "",
                },
            },
		},
        {
            "body",
            amqp.Delivery{
                Timestamp: timeNow,
                Body: []byte(`{"key":"value","int":1,"bool":true}`),
            },
            &entry.Entry{
                Timestamp: timeNow,
                Record: map[string]interface{}{
                    "acknowledger": nil,
                    "content_type": "",
                    "content_encoding": "",
                    "delivery_mode": uint8(0),
                    "priority": uint8(0),
                    "correlation_id": "",
                    "reply_to": "",
                    "message_id": "",
                    "expiration": "",
                    "type": "",
                    "user_id": "",
                    "app_id": "",
                    "consumer_tag": "",
                    "redelivered": false,
                    "exchange": "",
                    "headers": amqp.Table{},
                    "body": `{"key":"value","int":1,"bool":true}`,
                },
            },
        },
        {
            "headers",
            amqp.Delivery{
                Timestamp: timeNow,
                Body: []byte(`{"key":"value","int":1,"bool":true}`),
                Headers: amqp.Table{
                    "X-Forwarded-For": "10.0.0.1",
                    "X-Token": []string{
                        "00000",
                        "11111",
                    },
                },
            },
            &entry.Entry{
                Timestamp: timeNow,
                Record: map[string]interface{}{
                    "acknowledger": nil,
                    "content_type": "",
                    "content_encoding": "",
                    "delivery_mode": uint8(0),
                    "priority": uint8(0),
                    "correlation_id": "",
                    "reply_to": "",
                    "message_id": "",
                    "expiration": "",
                    "type": "",
                    "user_id": "",
                    "app_id": "",
                    "consumer_tag": "",
                    "redelivered": false,
                    "exchange": "",
                    "headers": amqp.Table{
                        "X-Forwarded-For": "10.0.0.1",
                        "X-Token": []string{
                            "00000",
                            "11111",
                        },
                    },
                    "body": `{"key":"value","int":1,"bool":true}`,
                },
            },
        },
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := parseMessage(tc.inputRecord)
			require.Equal(t, tc.expectedRecord, e)
		})
	}
}

func TestParseMessageZeroTime(t *testing.T) {
    cases := []struct {
        name           string
        inputRecord    amqp.Delivery
    }{
        {
            "zero-timestamp",
            amqp.Delivery{
                Timestamp: timeZero,
            },
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            e := parseMessage(tc.inputRecord)
            require.NotEqual(t, e.Timestamp, tc.inputRecord.Timestamp)
        })
    }
}

func TestValidTimestamp(t *testing.T) {
    loc, _ := time.LoadLocation("Asia/Shanghai")

    cases := []struct{
        name string
        inputRecord    time.Time
        expectedRecord bool
    }{
        {
            "zero",
            timeZero,
            false,
        },
        {
            "new-time",
            time.Time{},
            false,
        },
        {
            "time-now",
            timeNow,
            true,
        },
        {
            "shanghai",
            time.Now().In(loc),
            true,
        },
    }

    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) {
            require.Equal(t, tc.expectedRecord, validTimestamp(tc.inputRecord))
        })
    }
}

func BenchmarkParseMessage(b *testing.B) {
	for n := 0; n < b.N; n++ {
	    parseMessage(testBenchParseMessage)
	}
}

func initTestTime() (err error) {
    const zeroUTCSTR = "0001-01-01T00:00:00Z"
    const zeroUTCSTRLayout = "2006-01-02T15:04:05Z0700"
    timeNow = time.Now()
    timeZero, err = time.Parse(zeroUTCSTRLayout, zeroUTCSTR)
    return
}

func initTestBenchParseMessage() error {
    body := make(map[string]interface{})
    for i := 0; i < 100; i++ {
        key := strconv.Itoa(i)
        body[key] = []interface{}{
            "value",
            1,
            true,
        }
    }

    b, err := json.Marshal(body)
    if err != nil {
        return err
    }

    testBenchParseMessage = amqp.Delivery{
        Acknowledger: nil,
        Headers: map[string]interface{}{
            "X-Forwarded-For": "10.0.0.1",
            "X-Token": "38069b4c-b0c2-4276-9743-0760a6648459",
        },
        ContentType: "encoding/json",
        ContentEncoding: "gzip",
        DeliveryMode: 1,
        Priority: 1,
        CorrelationId: "2637168f-b07e-436a-b090-13f319e6363a",
        ReplyTo: "10.0.0.1",
        Expiration: "",
        MessageId: "001",
        Timestamp: time.Now(),
        Type: "app",
        UserId: "001",
        AppId: "001",
        ConsumerTag: "devel",
        MessageCount: 900,
        DeliveryTag: 2,
        Redelivered: false,
        Exchange: "webex",
        RoutingKey: "webq1",
        Body: b,
    }

    return nil
}

func init() {
    if err := initTestTime(); err != nil {
        panic(err)
    }

    if err := initTestBenchParseMessage(); err != nil {
        panic(err)
    }
}
