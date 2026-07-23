package queue

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
	name   string
}

var _ Queue = (*Redis)(nil)

func OpenRedis(ctx context.Context, address string, name string) (*Redis, error) {
	if strings.TrimSpace(address) == "" || strings.TrimSpace(name) == "" {
		return nil, errors.New("Redis address and queue name must not be empty")
	}

	client := redis.NewClient(&redis.Options{Addr: address})
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("ping Redis: %w", err)
	}

	return &Redis{client: client, name: name}, nil
}

func (q *Redis) Enqueue(ctx context.Context, jobID string) error {
	if strings.TrimSpace(jobID) == "" {
		return errors.New("queue job ID must not be empty")
	}
	if err := q.client.LPush(ctx, q.name, jobID).Err(); err != nil {
		return fmt.Errorf("enqueue job: %w", err)
	}
	return nil
}

func (q *Redis) Dequeue(ctx context.Context, timeout time.Duration) (string, error) {
	result, err := q.client.BRPop(ctx, timeout, q.name).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrEmpty
	}
	if err != nil {
		return "", fmt.Errorf("dequeue job: %w", err)
	}
	if len(result) != 2 || result[1] == "" {
		return "", ErrEmpty
	}
	return result[1], nil
}

func (q *Redis) Close() error {
	return q.client.Close()
}

func (q *Redis) Ping(ctx context.Context) error {
	if err := q.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping Redis: %w", err)
	}
	return nil
}
