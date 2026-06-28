package posthog

import (
	"backend/internal/config"

	ph "github.com/posthog/posthog-go"
)

// Client wraps the PostHog server-side client. When no API key is configured it
// is a no-op, so generated projects run without PostHog credentials in dev.
type Client struct {
	client ph.Client
}

func New(cfg *config.Config) (*Client, error) {
	if cfg.PostHogAPIKey == "" {
		return &Client{}, nil
	}
	host := cfg.PostHogHost
	if host == "" {
		host = "https://us.i.posthog.com"
	}
	c, err := ph.NewWithConfig(cfg.PostHogAPIKey, ph.Config{Endpoint: host})
	if err != nil {
		return nil, err
	}
	return &Client{client: c}, nil
}

func (c *Client) Capture(distinctID, event string, props map[string]any) {
	if c.client == nil {
		return
	}
	properties := ph.NewProperties()
	for k, v := range props {
		properties.Set(k, v)
	}
	_ = c.client.Enqueue(ph.Capture{
		DistinctId: distinctID,
		Event:      event,
		Properties: properties,
	})
}

func (c *Client) Close() error {
	if c.client == nil {
		return nil
	}
	return c.client.Close()
}
