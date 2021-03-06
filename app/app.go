package app

import (
	"fmt"

	"github.com/Shelex/webhook-listener/notification"
	"github.com/Shelex/webhook-listener/repository"
	"github.com/Shelex/webhook-listener/scheduler"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	PubSub       repository.PubSub
	Repository   repository.Storage
	Router       *fiber.App
	Notification *notification.Notification
	Cron         *scheduler.Scheduler
}

func NewApp() (*App, error) {
	storage, err := repository.NewStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %s", err)
	}

	webhooks := repository.NewPubSub()

	repository.Subscribe(webhooks)

	notifications := notification.New(webhooks)

	router := ProvideRouter()

	cron := scheduler.NewScheduler()

	return &App{
		PubSub:       webhooks,
		Repository:   storage,
		Router:       router,
		Notification: notifications,
		Cron:         cron,
	}, nil
}
