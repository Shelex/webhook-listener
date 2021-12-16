package app

import (
	"fmt"

	"github.com/Shelex/webhook-listener/notification"
	"github.com/Shelex/webhook-listener/repository"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/go-chi/chi"
)

type App struct {
	Queue        *gochannel.GoChannel
	Repository   repository.Storage
	Router       *chi.Mux
	Notification *notification.Notification
	Logger       watermill.LoggerAdapter
}

func NewApp() (*App, error) {

	logger := watermill.NewStdLogger(true, true).With(
		watermill.LogFields{},
	)

	pubSub := gochannel.NewGoChannel(
		gochannel.Config{
			Persistent: true,
		},
		logger,
	)

	storage, err := repository.NewRedis()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %s", err)
	}

	if err := repository.Subscribe(pubSub); err != nil {
		return nil, fmt.Errorf("repository failed to subscribe: %s", err)
	}

	notification, err := notification.New(pubSub)
	if err != nil {
		return nil, fmt.Errorf("failed to register notifications: %s", err)
	}

	router := ProvideRouter()

	return &App{
		Queue:        pubSub,
		Repository:   storage,
		Router:       router,
		Notification: notification,
		Logger:       logger,
	}, nil
}
