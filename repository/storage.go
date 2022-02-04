package repository

import (
	"log"
	"time"

	"github.com/Shelex/webhook-listener/entities"
)

var DB Storage

type Pagination struct {
	Offset int64
	Limit  int64
}

type Storage interface {
	Add(hook entities.Hook) error
	Get(channel string, pagination Pagination) ([]entities.Hook, int64, error)
	Delete(channel string) error
	ClearExpired() error
}

func Subscribe(pubSub PubSub) {
	messages := pubSub.Subscribe()
	go Persist(messages)
}

func Persist(messages <-chan *entities.Hook) {
	for hook := range messages {
		go DB.Add(*hook)
		log.Printf("storage - acknowledged message %s", hook.ID)
	}
}

func GetExpiryDate() int64 {
	return time.Now().Add(time.Duration(-3*24) * time.Hour).UTC().Unix()
}
