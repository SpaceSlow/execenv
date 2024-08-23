package config

import (
	"encoding"
	"flag"
	"fmt"
	"net"
)

var (
	_ flag.Value               = (*CIDR)(nil)
	_ encoding.TextUnmarshaler = (*CIDR)(nil)
)

type CIDR struct {
	net.IPNet
}

func NewCIDR(s string) CIDR {
	var c CIDR
	c.Set(s)
	return c
}

func (c CIDR) String() string {
	return c.IPNet.String()
}

func (c *CIDR) Set(s string) error {
	if s == "" {
		return nil
	}
	_, parsedCIDR, err := net.ParseCIDR(s)
	if err != nil {
		return fmt.Errorf("parse cidr error: %w", err)
	}
	*c = CIDR{*parsedCIDR}
	return nil
}

func (c *CIDR) UnmarshalText(text []byte) error {
	return c.Set(string(text))
}
