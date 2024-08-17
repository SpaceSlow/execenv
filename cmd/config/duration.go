package config

import (
	"encoding"
	"encoding/json"
	"flag"
	"time"
)

var (
	_ flag.Value               = (*Duration)(nil)
	_ encoding.TextUnmarshaler = (*Duration)(nil)
	_ json.Marshaler           = (*Duration)(nil)
	_ json.Unmarshaler         = (*Duration)(nil)
)

type Duration struct {
	time.Duration
}

func (d *Duration) Set(s string) error {
	var err error
	d.Duration, err = time.ParseDuration(s)
	return err
}

func (d *Duration) UnmarshalText(text []byte) error {
	return d.Set(string(text))
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	var err error
	d.Duration, err = time.ParseDuration(value)
	return err
}
