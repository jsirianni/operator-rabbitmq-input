package rabbitmq

import (
    "fmt"
	"context"

	"github.com/open-telemetry/opentelemetry-log-collection/operator"
	"github.com/open-telemetry/opentelemetry-log-collection/operator/helper"

    "go.uber.org/zap"
)

func init() {
	operator.Register("rabbitmq_input", func() operator.Builder { return NewQueueInputConfig("") })
}

// NewQueueInputConfig creates a new message queue input config with default values
func NewQueueInputConfig(operatorID string) *QueueInputConfig {
	return &QueueInputConfig{
		InputConfig: helper.NewInputConfig(operatorID, "rabbitmq_input"),
	}
}

// QueueInputConfig is the configuration of a message queue input operator.
type QueueInputConfig struct {
	helper.InputConfig `yaml:",inline"`

    // required
    Username   string `json:"username,omitempty" yaml:"username,omitempty"`
    Password   string `json:"password,omitempty" yaml:"password,omitempty"`
    RoutingKey string `json:"routing_key,omitempty" yaml:"routing_key,omitempty"`

    // optional
    Address string `json:"address,omitempty" yaml:"address,omitempty"`
    Port    string `json:"port,omitempty" yaml:"port,omitempty"`
    VHost   string `json:"vhost,omitempty" yaml:"vhost,omitempty"`
    Tag     string `json:"tag,omitempty" yaml:"tag,omitempty"`
}

// Build will build a message queue input operator.
func (c *QueueInputConfig) Build(context operator.BuildContext) ([]operator.Operator, error) {
	inputOperator, err := c.InputConfig.Build(context)
	if err != nil {
		return nil, err
	}

    if c.Username == "" {
        return nil, fmt.Errorf("missing required rabbitmq_input parameter 'username'")
    }

    if c.Password == "" {
        return nil, fmt.Errorf("missing required rabbitmq_input parameter 'password'")
    }

    if c.RoutingKey == "" {
        return nil, fmt.Errorf("missing required rabbitmq_input parameter 'routing_key'")
    }

    if c.Address == "" {
        c.Address = "localhost"
    }

    if c.Port == "" {
        c.Port = "5672"
    }

    // Go routine count is 2. This could be made available as an option in the
    // future, however, testing has shown that anything more than four threads
    // is not significantly faster. One go routine is capable of 12.5k messages
    // per second on an AMD Ryzen 3700X (16c32t). Similar performance was observed
    // with four go routines. The Rabbitmq environment was a single node on the
    // same system.
    workerCount := 2

    q, err := newRabbitmq(
        c.Username, c.Password,
        c.Address, c.Port,
        c.VHost, c.RoutingKey, c.Tag, workerCount,
    )
    if err != nil {
        return nil, err
    }

	queueInput := &QueueInput{
		InputOperator: inputOperator,
        rabbit: q,
	}
	return []operator.Operator{queueInput}, nil
}

// QueueInput is an operator that reads input from a message queue.
type QueueInput struct {
	helper.InputOperator
	cancel context.CancelFunc
    rabbit Rabbitmq
}

// Start will start generating log entries.
func (g *QueueInput) Start() error {
    ctx, cancel := context.WithCancel(context.Background())
    g.cancel = cancel

    go g.connectionManager(ctx)
    go g.workerManager(ctx)

    return nil
}

// Stop will stop generating logs.
func (g *QueueInput) Stop() error {
    // canceling the context will notify consumer go routines to stop
    // after finishing their in flight request
	g.cancel()
    g.rabbit.wg.Wait()
    if err := g.rabbit.close(); err != nil {
        g.Errorw("Error while closing connection to Rabbitmq", zap.Error(err))
    } else {
        g.Infow("Closed connection to Rabbitmq")
    }
	return nil
}
