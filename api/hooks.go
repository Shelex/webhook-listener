package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Shelex/webhook-listener/app"
	"github.com/Shelex/webhook-listener/entities"
	"github.com/Shelex/webhook-listener/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var buf = &bytes.Buffer{}

type Controller struct {
	app *app.App
}

// getMessages godoc
// @Summary get messages for channel
// @Description get messages for channel
// @Accept  json
// @Produce  json
// @Param  channel path string true "name"
// @Param  offset query int  false  "pagination offset"
// @Param  limit query int  false  "pagination limit"
// @Success 200 {array} []entities.Hook
// @Router /api/{channel} [get]
func (c *Controller) getMessages(w http.ResponseWriter, r *http.Request) {
	channel := chi.URLParam(r, "channel")
	query := r.URL.Query()
	limit, _ := strconv.Atoi(query.Get("limit"))
	offset, _ := strconv.Atoi(query.Get("offset"))

	hooks, count, err := c.app.Repository.Get(channel, repository.Pagination{
		Limit:  int64(limit),
		Offset: int64(offset),
	})

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		render.JSON(w, r, map[string]interface{}{
			"data":  hooks,
			"count": count,
		})
	}
}

// deleteChannel godoc
// @Summary delete channel
// @Description delete channel
// @Accept  json
// @Param  channel path string true "name"
// @Success 200
// @Router /api/{channel} [delete]
func (c *Controller) deleteMessages(w http.ResponseWriter, r *http.Request) {
	channel := chi.URLParam(r, "channel")
	err := c.app.Repository.Delete(channel)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

// addMessage godoc
// @Summary post message
// @Description post message
// @Accept  json
// @Param  channel path string true "name"
// @Param  message body object true "message"
// @Param  failUntil query int false  "fail requests until timestamp"
// @Param  justreply query boolean false  "do not handle message by service and just reply with status code"
// @Success 200
// @Router /api/{channel} [post]
func (c *Controller) addMessage(w http.ResponseWriter, r *http.Request) {
	justReply := r.URL.Query().Get("justreply")
	if justReply != "" {
		if justReply == "true" {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
		return
	}

	buf.Reset()
	buf.ReadFrom(r.Body)
	body := buf.Bytes()

	headers, _ := json.Marshal(r.Header)

	channel := chi.URLParam(r, "channel")

	failUntil := r.URL.Query().Get("failUntil")

	hook := entities.Hook{
		ID:         uuid.NewString(),
		Channel:    channel,
		Created_at: time.Now().UTC().Unix(),
		Payload:    string(body),
		Headers:    string(headers),
		StatusOK:   true,
	}

	if failUntil != "" {
		currentTime := time.Now().UTC()
		expiryTimestamp, _ := strconv.Atoi(failUntil)

		expirationTime := time.Unix(int64(expiryTimestamp), 0).UTC()

		if expirationTime.After(currentTime) {
			hook.StatusOK = false
		}
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	c.app.PubSub.Publish(hook)

	log.Println("have message for channel " + channel)
}

func RegisterControllers(app *app.App) {
	controller := Controller{app}
	app.Router.Post("/api/{channel}", controller.addMessage)
	app.Router.Get("/api/{channel}", controller.getMessages)
	app.Router.Delete("/api/{channel}", controller.deleteMessages)
	app.Router.Get("/listen/{channel}", app.Notification.Handle)
}
