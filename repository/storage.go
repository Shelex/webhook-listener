package repository

import (
	"log"
	"sync"
	"time"

	"github.com/Shelex/webhook-listener/entities"
)

var DB Storage

type Pagination struct {
	Offset int64
	Limit  int64
}

type Storage interface {
	Add(hooks ...*entities.Hook) error
	Get(channel string, pagination Pagination) ([]entities.Hook, int64, error)
	Delete(channel string) error
	ClearExpired() error
}

func Subscribe(pubSub PubSub) {
	messages := pubSub.Subscribe()
	go Persist(messages)
}

func Persist(messages <-chan *entities.Hook) {
	batch := Batch{}

	go batch.Process()

	for hook := range messages {
		batch.Add(hook)
	}
}

type Batch struct {
	messages []*entities.Hook
	mux      sync.Mutex
	cancel   chan struct{}
}

func (batch *Batch) Add(message *entities.Hook) {
	batch.mux.Lock()
	batch.messages = append(batch.messages, message)
	batch.mux.Unlock()
}

func (batch *Batch) Process() {
	ticker := time.NewTicker(500 * time.Millisecond)

	go func() {
		for {
			select {
			case <-ticker.C:
				batch.mux.Lock()
				if len(batch.messages) > 0 {
					if err := DB.Add(batch.messages...); err != nil {
						log.Printf("storage - failed to save message: %s", err.Error())
					}
					log.Printf("redis: saved %d hooks", len(batch.messages))
					batch.messages = nil
				}
				batch.mux.Unlock()
			case <-batch.cancel:
				ticker.Stop()
				return
			}
		}
	}()
}

func GetExpiryDate() int64 {
	return time.Now().Add(time.Duration(-3*24) * time.Hour).UTC().Unix()
}
