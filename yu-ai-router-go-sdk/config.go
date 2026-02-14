package yuairoutersdk

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL        = "http://localhost:8123/api"
	defaultConnectTimeout = 10 * time.Second
	defaultReadTimeout    = 30 * time.Second
	defaultMaxRetries     = 3
	defaultRetryDelay     = 1 * time.Second
)

type Config struct {
	APIKey         string
	BaseURL        string
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	MaxRetries     int
	RetryDelay     time.Duration
	HTTPClient     *http.Client
}

func DefaultConfig(apiKey string) Config {
	return Config{
		APIKey:         apiKey,
		BaseURL:        defaultBaseURL,
		ConnectTimeout: defaultConnectTimeout,
		ReadTimeout:    defaultReadTimeout,
		MaxRetries:     defaultMaxRetries,
		RetryDelay:     defaultRetryDelay,
	}
}

func (c *Config) normalize() error {
	c.APIKey = strings.TrimSpace(c.APIKey)
	c.BaseURL = strings.TrimRight(strings.TrimSpace(c.BaseURL), "/")

	if c.APIKey == "" {
		return errors.New("api key cannot be empty")
	}
	if c.BaseURL == "" {
		c.BaseURL = defaultBaseURL
	}
	if c.ConnectTimeout <= 0 {
		c.ConnectTimeout = defaultConnectTimeout
	}
	if c.ReadTimeout <= 0 {
		c.ReadTimeout = defaultReadTimeout
	}
	if c.MaxRetries < 0 {
		c.MaxRetries = defaultMaxRetries
	}
	if c.RetryDelay < 0 {
		c.RetryDelay = defaultRetryDelay
	}
	return nil
}
