package client

import (
	"github.com/SpaceSlow/execenv/internal/config"
	"github.com/SpaceSlow/execenv/internal/metrics"
)

type Client struct {
	sender Sender
}

func NewClient() (*Client, error) {
	cfg, err := config.GetAgentConfig()
	if err != nil {
		return nil, err
	}

	var sender Sender
	if cfg.UsedGRPCAgent {
		sender, err = newGrpcSender()
	} else {
		sender, err = newHTTPSender()
	}
	if err != nil {
		return nil, err
	}
	return &Client{sender: sender}, nil
}

func (c *Client) Send(metrics []metrics.Metric) error {
	return c.sender.Send(metrics)
}
