package api

import (
	"bytes"
	"fmt"
	"log"
	stdHttp "net/http"
	"strconv"
	"time"

	"github.com/Shelex/webhook-listener/app"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-http/pkg/http"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-chi/chi"
)

func RegisterHttpListener(app *app.App, addr string) (*http.Subscriber, error) {
	buf := &bytes.Buffer{}

	httpSubscriber, err := http.NewSubscriber(
		addr,
		http.SubscriberConfig{
			Router: app.Router,
			UnmarshalMessageFunc: func(topic string, request *stdHttp.Request) (*message.Message, error) {
				buf.Reset()
				buf.ReadFrom(request.Body)
				body := buf.Bytes()

				channel := chi.URLParam(request, "channel")

				message := message.NewMessage(
					watermill.NewUUID(),
					body,
				)

				failUntil := request.URL.Query().Get("failUntil")
				statusOk := "true"
				if failUntil != "" {
					currentTime := time.Now().UTC()
					expiryTimestamp, _ := strconv.Atoi(failUntil)

					expirationTime := time.Unix(int64(expiryTimestamp), 0).UTC()

					if expirationTime.After(currentTime) {
						statusOk = "false"
						message.Nack()
					}
				}

				message.Metadata.Set("channel", channel)
				message.Metadata.Set("statusOk", statusOk)

				log.Println("have message for channel " + channel)

				return message, nil
			},
		},
		app.Logger,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create http subscriber: %s", err)
	}
	return httpSubscriber, nil
}
