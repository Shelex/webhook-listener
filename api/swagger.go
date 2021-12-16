package api

import (
	"github.com/Shelex/webhook-listener/app"
	_ "github.com/Shelex/webhook-listener/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterSwagger(app *app.App) {
	app.Router.Get("/swagger/*", httpSwagger.WrapHandler)
}
