package services

import (
	"time"

	"github.com/quipo/statsd"
)

// StatsdClient represents the interface for a set of operations that can be performed
// to track performance of queues and topics in adfilteringservice
type StatsdClient interface {
	Incr(id string, value int64) error
	Decr(id string, value int64) error
	Timing(id string, value int64) error
	PrecisionTiming(id string, value time.Duration) error
	IncrGauge(id string, value int64) error
	DecrGauge(id string, value int64) error
	SetGauge(id string, value int64) error
}

// RealStatsdClient will report stats to a StatsD compatible service
type RealStatsdClient struct {
	prefix   string
	address  string
	client   *statsd.StatsdClient
	interval time.Duration
}

// StatsdConfig is
type StatsdConfig struct {
	Address       string
	Prefix        string
	FlushInterval int
	Type          string
	Enabled       bool
}

// NewStatsdClient returns a stats client based on the monitoring config.
// for production where type is statsd, it returns the result of NewStatsdClient
// for dev or where type is not set, it returns the result of NewNOOPClient
func NewStatsdClient(config *StatsdConfig) StatsdClient {
	if !config.Enabled {
		return NewNOOPStatsdClient()
	}

	switch config.Type {
	case "statsd":
		return NewRealStatsdClient(config.Address, config.Prefix, time.Second*time.Duration(config.FlushInterval))
	default:
		return NewNOOPStatsdClient()
	}
}

// NewRealStatsdClient will create a new StatsdClient to be used for reporting metrics
func NewRealStatsdClient(address string, prefix string, interval time.Duration) StatsdClient {
	client := &RealStatsdClient{
		prefix:   prefix,
		interval: interval,
		address:  address,
	}
	client.client = statsd.NewStatsdClient(address, prefix)
	client.client.CreateSocket()
	return client
}

// Incr increases the value of a given counter
func (c *RealStatsdClient) Incr(id string, value int64) error {
	return c.client.Incr(id, value)
}

// Decr decreases the value of a given counter
func (c *RealStatsdClient) Decr(id string, value int64) error {
	return c.client.Decr(id, value)
}

// Timing sets a timing value
func (c *RealStatsdClient) Timing(id string, value int64) error {
	return c.client.Timing(id, value)
}

// PrecisionTiming sets a timing value using a duration
func (c *RealStatsdClient) PrecisionTiming(id string, duration time.Duration) error {
	return c.client.PrecisionTiming(id, duration)
}

// IncrGauge increases the value of a given gauge delta
func (c *RealStatsdClient) IncrGauge(id string, value int64) error {
	return c.client.GaugeDelta(id, value)
}

// DecrGauge decreases the value of a given gauge delta
func (c *RealStatsdClient) DecrGauge(id string, value int64) error {
	return c.client.GaugeDelta(id, -value)
}

// SetGauge sets the level of the given gauge
func (c *RealStatsdClient) SetGauge(id string, value int64) error {
	return c.client.Gauge(id, value)
}

// NOOPStatsdClient is to sub in when we don't want to write stats
type NOOPStatsdClient struct {
}

// NewNOOPStatsdClient returns a new NOOPClient
func NewNOOPStatsdClient() *NOOPStatsdClient {
	return &NOOPStatsdClient{}
}

// Incr does nothing
func (c *NOOPStatsdClient) Incr(id string, value int64) error {
	return nil
}

// Decr does nothing
func (c *NOOPStatsdClient) Decr(id string, value int64) error {
	return nil
}

// Timing does nothing
func (c *NOOPStatsdClient) Timing(id string, value int64) error {
	return nil
}

// PrecisionTiming does nothing
func (c *NOOPStatsdClient) PrecisionTiming(id string, duration time.Duration) error {
	return nil
}

// IncrGauge does nothing
func (c *NOOPStatsdClient) IncrGauge(id string, value int64) error {
	return nil
}

// DecrGauge does nothing
func (c *NOOPStatsdClient) DecrGauge(id string, value int64) error {
	return nil
}

// SetGauge does nothing
func (c *NOOPStatsdClient) SetGauge(id string, value int64) error {
	return nil
}
