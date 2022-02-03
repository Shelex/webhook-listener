package api

import (
	"net/http"
	"strconv"

	"github.com/Shelex/webhook-listener/app"
	"github.com/Shelex/webhook-listener/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
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

func RegisterControllers(app *app.App) {
	controller := Controller{app}
	app.Router.Get("/api/{channel}", controller.getMessages)
	app.Router.Delete("/api/{channel}", controller.deleteMessages)
	app.Router.Get("/listen/{channel}", app.Notification.Handle)
}
