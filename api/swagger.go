package api

import (
	"github.com/Shelex/webhook-listener/app"
	_ "github.com/Shelex/webhook-listener/docs"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func RegisterSwagger(app *app.App) {
	app.Router.Get("/swagger/*", fiberSwagger.WrapHandler)
}
