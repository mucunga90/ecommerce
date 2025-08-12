package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	client *redis.Client
	stream string
}

func NewPublisher(client *redis.Client, stream string) *Publisher {
	return &Publisher{client: client, stream: stream}
}

func (p *Publisher) Publish(ctx context.Context, eventType string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	args := &redis.XAddArgs{
		Stream: p.stream,
		Values: map[string]interface{}{
			"type":    eventType,
			"payload": string(data),
		},
	}

	if err := p.client.XAdd(ctx, args).Err(); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}
