package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Shelex/webhook-listener/entities"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
)

var DB Storage

type Storage interface {
	Add(hook entities.Hook) error
	Get(channel string, pagination Pagination) ([]entities.Hook, int64, error)
	Delete(channel string) error
}

func Subscribe(pubSub *gochannel.GoChannel) error {
	messages, err := pubSub.Subscribe(context.Background(), "webhooks")

	if err != nil {
		return errors.New("failed to subscribe for messages")
	}

	go Persist(messages)

	return nil
}

func Persist(messages <-chan *message.Message) {
	for msg := range messages {
		hook := entities.Hook{
			ID:         msg.UUID,
			Channel:    msg.Metadata.Get("channel"),
			Created_at: time.Now().UTC().Unix(),
			Payload:    string(msg.Payload),
		}

		go DB.Add(hook)
		log.Printf("storage - acknowledged message %s", msg.UUID)
		msg.Ack()
	}
}
