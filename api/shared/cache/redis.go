package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis *redis.Client
}

func New(addr string, db int) *Client {
	return &Client{redis: redis.NewClient(&redis.Options{Addr: addr, DB: db})}
}

func (c *Client) GetJSON(ctx context.Context, key string, dst any) (bool, error) {
	val, err := c.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(val), dst); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	body, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, key, body, ttl).Err()
}
