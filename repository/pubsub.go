package repository

import (
	"encoding/json"

	"github.com/Shelex/webhook-listener/entities"
	"github.com/redis/go-redis/v9"
)

type PubSub interface {
	Subscribe() <-chan *entities.Hook
	Publish(payload entities.Hook)
}

type PubSubClient struct {
	client  *redis.Client
	channel string
}

func NewPubSub() *PubSubClient {
	return &PubSubClient{
		client:  NewRedisClient(),
		channel: "webhooks",
	}
}

func (p *PubSubClient) Subscribe() <-chan *entities.Hook {
	subscriber := p.client.Subscribe(ctx, p.channel)
	return redisMessageToHook(subscriber.Channel())
}

func (p *PubSubClient) Publish(hook entities.Hook) {
	p.client.Publish(ctx, p.channel, hook)
}

func redisMessageToHook(redisChannel <-chan *redis.Message) <-chan *entities.Hook {
	channel := make(chan *entities.Hook)

	go func() {
		for msg := range redisChannel {
			var hook entities.Hook
			json.Unmarshal([]byte(msg.Payload), &hook)
			channel <- &hook
		}
	}()

	return channel
}
