package sentry

import (
	"time"

	"backend/internal/config"

	sentrygo "github.com/getsentry/sentry-go"
)

// Client wraps the Sentry server-side SDK. When no DSN is configured it is a
// no-op, so generated projects run without Sentry credentials in dev.
type Client struct {
	enabled bool
}

func New(cfg *config.Config) (*Client, error) {
	if cfg.SentryDSN == "" {
		return &Client{}, nil
	}
	if err := sentrygo.Init(sentrygo.ClientOptions{
		Dsn:              cfg.SentryDSN,
		Environment:      cfg.Env,
		EnableTracing:    true,
		TracesSampleRate: 1.0,
	}); err != nil {
		return nil, err
	}
	return &Client{enabled: true}, nil
}

func (c *Client) Close() error {
	if !c.enabled {
		return nil
	}
	sentrygo.Flush(2 * time.Second)
	return nil
}
