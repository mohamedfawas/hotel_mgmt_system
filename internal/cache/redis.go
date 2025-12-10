package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	RedisURL string // used for upstash redis
	Host     string
	Port     int
	Password string
	DB       int
}

type Client struct {
	Client *redis.Client
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	var opts *redis.Options
	var err error

	if cfg.RedisURL != "" {
		// Upstash / cloud redis case
		opts, err = redis.ParseURL(cfg.RedisURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse redis url: %w", err)
		}
	} else {
		// Local Redis case
		opts = &redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Password: cfg.Password,
			DB:       cfg.DB,
		}
	}

	client := redis.NewClient(opts)

	healthCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(healthCtx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &Client{Client: client}, nil
}

func (c *Client) Set(ctx context.Context, key string, value interface{}) error {
	return c.Client.Set(ctx, key, value, 0).Err()
}

func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := c.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}
