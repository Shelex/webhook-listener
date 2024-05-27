package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Shelex/webhook-listener/app"
	"github.com/Shelex/webhook-listener/entities"
	"github.com/Shelex/webhook-listener/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

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
// @Success 200 {array} entities.HooksByChannel
// @Router /api/{channel} [get]
func (c *Controller) getMessages(ctx *fiber.Ctx) error {
	channel := ctx.Params("channel")
	limit, _ := strconv.Atoi(ctx.Query("limit"))
	offset, _ := strconv.Atoi(ctx.Query("offset"))

	hooks, count, err := c.app.Repository.Get(channel, repository.Pagination{
		Limit:  int64(limit),
		Offset: int64(offset),
	})

	if err != nil {
		ctx.SendStatus(400)
		return err
	} else {
		ctx.JSON(
			entities.HooksByChannel{
				Data:  hooks,
				Count: count,
			},
		)
	}
	return nil
}

// deleteChannel godoc
// @Summary delete channel
// @Description delete channel
// @Accept  json
// @Param  channel path string true "name"
// @Success 200
// @Router /api/{channel} [delete]
func (c *Controller) deleteMessages(ctx *fiber.Ctx) error {
	channel := ctx.Params("channel")
	err := c.app.Repository.Delete(channel)
	if err != nil {
		ctx.SendStatus(400)
		return err
	}
	return nil
}

// addMessage godoc
// @Summary post message
// @Description post message
// @Accept  json
// @Param  channel path string true "name"
// @Param  message body map[string]interface{} true "message" Example(entities.HookExample)
// @Param  failUntil query int false  "fail requests until timestamp"
// @Param  justreply query boolean false  "do not handle message by service and just reply with status code"
// @Success 200
// @Router /api/{channel} [post]
func (c *Controller) addMessage(ctx *fiber.Ctx) error {
	justReply := ctx.Query("justreply")
	if justReply != "" {
		if justReply == "true" {
			ctx.SendStatus(http.StatusOK)
		} else {
			ctx.SendStatus(http.StatusServiceUnavailable)
		}
		return nil
	}

	body := ctx.Body()

	headers, _ := json.Marshal(ctx.GetReqHeaders())

	channel := ctx.Params("channel")

	failUntil := ctx.Query("failUntil")

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
		ctx.SendStatus(http.StatusServiceUnavailable)
	}

	c.app.PubSub.Publish(hook)

	return nil
}

func RegisterControllers(app *app.App) {
	controller := Controller{app}
	api := app.Router.Group("/api")
	channel := "/:channel"
	api.Post(channel, controller.addMessage)
	api.Get(channel, controller.getMessages)
	api.Delete(channel, controller.deleteMessages)
	app.Router.Use("/listen", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Router.Get("/listen/:channel", websocket.New(app.Notification.Handler))
}
