package notification

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Shelex/webhook-listener/entities"
	"github.com/Shelex/webhook-listener/repository"
	"github.com/go-chi/chi"
	"gopkg.in/olahol/melody.v1"
)

type Notification struct {
	router *melody.Melody
}

func New() *Notification {
	notification := Notification{
		router: melody.New(),
	}

	return &notification
}

func (n *Notification) Subscribe(pubSub repository.PubSub) {
	messages := pubSub.Subscribe()
	go n.Publish(messages)
}

func (n *Notification) Handle(w http.ResponseWriter, r *http.Request) {
	channel := chi.URLParam(r, "channel")
	log.Println("subscribing to channel: " + channel)
	n.router.HandleRequestWithKeys(w, r, map[string]interface{}{
		"channel": channel,
	})
}

func (n *Notification) Publish(messages <-chan *entities.Hook) {
	for hook := range messages {
		messageJson, _ := json.Marshal(map[string]interface{}{
			"payload": hook.Payload,
			"headers": hook.Headers,
			"ok":      hook.StatusOK,
		})

		go n.router.BroadcastFilter(messageJson, func(s *melody.Session) bool {
			channel, ok := s.Keys["channel"]
			if !ok {
				return false
			}
			return hook.Channel == channel
		})
		log.Printf("notification - acknowledged message %s", hook.ID)
	}
}
