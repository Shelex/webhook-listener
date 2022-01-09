package notification

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-chi/chi"
	"gopkg.in/olahol/melody.v1"
)

type Notification struct {
	router *melody.Melody
}

func New(channel *gochannel.GoChannel) (*Notification, error) {
	messages, err := channel.Subscribe(context.Background(), "webhooks")
	if err != nil {
		return nil, errors.New("failed to subscribe for messages")
	}

	notification := Notification{
		router: melody.New(),
	}

	go notification.Publish(messages)

	return &notification, nil
}

func (n *Notification) Handle(w http.ResponseWriter, r *http.Request) {
	channel := chi.URLParam(r, "channel")
	log.Println("subscribing to channel: " + channel)
	n.router.HandleRequestWithKeys(w, r, map[string]interface{}{
		"channel": channel,
	})
}

func (n *Notification) Publish(messages <-chan *message.Message) {
	for msg := range messages {
		messageJson, _ := json.Marshal(map[string]interface{}{
			"payload": string(msg.Payload),
			"ok":      msg.Metadata.Get("statusOk") == "true",
		})

		go n.router.BroadcastFilter(messageJson, func(s *melody.Session) bool {
			channel, ok := s.Keys["channel"]
			if !ok {
				return false
			}
			return msg.Metadata.Get("channel") == channel
		})
		log.Printf("notification - acknowledged message %s", msg.UUID)
		msg.Ack()
	}
}
