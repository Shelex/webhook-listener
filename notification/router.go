package notification

import (
	"log"

	"github.com/Shelex/webhook-listener/repository"
	"github.com/gofiber/websocket/v2"
)

type Notification struct {
	pubSub repository.PubSub
}

func New(webhooks repository.PubSub) *Notification {
	return &Notification{
		pubSub: webhooks,
	}
}

func (n *Notification) Handler(c *websocket.Conn) {
	for hook := range n.pubSub.Subscribe() {
		if (c.Params("channel")) == hook.Channel {
			if err := c.WriteJSON(map[string]interface{}{
				"payload": hook.Payload,
				"headers": hook.Headers,
				"ok":      hook.StatusOK,
			}); err != nil {
				log.Printf("failed to send ws message: %s", err.Error())
			}

		}
		log.Printf("notification - acknowledged message %s", hook.ID)
	}
}
