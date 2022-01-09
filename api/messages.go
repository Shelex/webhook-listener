package api

import (
	"fmt"

	"github.com/Shelex/webhook-listener/app"
	"github.com/ThreeDotsLabs/watermill-http/pkg/http"
	"github.com/ThreeDotsLabs/watermill/message"
)

// postMessage godoc
// @Summary post message
// @Description post message
// @Accept  json
// @Param  channel path string true "name"
// @Param  message body object true "message"
// @Param  failUntil query int  false  "fail requests until timestamp"
// @Success 200
// @Router /api/{channel} [post]
func RouteMessages(app *app.App, subscriber *http.Subscriber) (*message.Router, error) {
	messageHandler, err := message.NewRouter(
		message.RouterConfig{},
		app.Logger,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create message handler: %s", err)
	}

	// POST route from http subscriber to publisher (channel)
	messageHandler.AddHandler(
		"http_to_channel", // name for debug
		"/api/{channel}",  // route
		subscriber,
		"webhooks", // topic
		app.PubSub,
		func(msg *message.Message) ([]*message.Message, error) {
			return []*message.Message{msg}, nil
		},
	)

	return messageHandler, nil
}
