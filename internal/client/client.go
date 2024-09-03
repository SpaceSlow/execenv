package client

import "github.com/SpaceSlow/execenv/internal/metrics"

type Client struct {
	strategy Sender
}

func NewClient() (*Client, error) {
	strategy, err := newHttpStrategy()
	if err != nil {
		return nil, err
	}
	return &Client{strategy: strategy}, nil
}

func (c *Client) Send(metrics []metrics.Metric) error {
	return c.strategy.Send(metrics)
}
