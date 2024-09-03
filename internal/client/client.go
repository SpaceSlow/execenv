package client

import "github.com/SpaceSlow/execenv/internal/metrics"

type Client struct {
	sender Sender
}

func NewClient() (*Client, error) {
	sender, err := newHttpSender()
	if err != nil {
		return nil, err
	}
	return &Client{sender: sender}, nil
}

func (c *Client) Send(metrics []metrics.Metric) error {
	return c.sender.Send(metrics)
}
