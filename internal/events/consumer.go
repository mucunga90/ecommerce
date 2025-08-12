package events

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mucunga90/ecommerce/internal"
	"github.com/redis/go-redis/v9"
)

type notifierConsumer struct {
	client   *redis.Client
	stream   string
	group    string
	consumer string
	notifier notifier
}

func NewNotifierConsumer(client *redis.Client, stream, group, consumer string, notifier notifier) *notifierConsumer {
	return &notifierConsumer{client, stream, group, consumer, notifier}
}

func (n *notifierConsumer) Start(ctx context.Context) {
	for {
		res, err := n.client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    n.group,
			Consumer: n.consumer,
			Streams:  []string{n.stream, ">"},
			Count:    1,
			Block:    0,
		}).Result()

		if err != nil {
			log.Println("Redis read error:", err)
			continue
		}

		for _, stream := range res {
			for _, msg := range stream.Messages {
				eventType := msg.Values["type"].(string)
				if eventType == "order.created" {
					var payload internal.OrderCreatedEvent
					if err := json.Unmarshal([]byte(msg.Values["payload"].(string)), &payload); err != nil {
						log.Println("JSON decode error:", err)
						continue
					}

					// Send SMS & Email
					if err = n.notifier.SendOrderSMS(payload); err != nil {
						log.Printf("Failed to send SMS: %v", err)
						return
					}

					if err = n.notifier.SendOrderEmail(payload); err != nil {
						log.Printf("Failed to send Email: %v", err)
						return
					}

					// Acknowledge message
					n.client.XAck(ctx, n.stream, n.group, msg.ID)
				}
			}
		}
	}
}

type notifier interface {
	SendOrderSMS(evt internal.OrderCreatedEvent) error
	SendOrderEmail(evt internal.OrderCreatedEvent) error
}
