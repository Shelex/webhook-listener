package main

import (
	"context"
	"flag"
	"log"

	//_ "net/http/pprof"

	"github.com/Shelex/webhook-listener/api"
	"github.com/Shelex/webhook-listener/app"
)

var (
	httpAddr = flag.String("http", ":8080", "The address for the http subscriber")
)

// @title webhook listener API
// @version 1.0
// @description webhook listener api
// @host localhost:8080
// @BasePath /
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email ovr.shevtsov@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	flag.Parse()

	// download results for pprof profiler with "make prof"
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	app, err := app.NewApp()
	if err != nil {
		log.Fatal(err)
	}

	api.RegisterControllers(app)
	api.RegisterSwagger(app)

	// clear expired items at 0 min every third hour
	app.Cron.Schedule("0 */3 * * *", app.Repository.ClearExpired)

	httpSubscriber, err := api.RegisterHttpListener(app, *httpAddr)
	if err != nil {
		log.Fatal(err)
	}

	messageRouter, err := api.RouteMessages(app, httpSubscriber)
	if err != nil {
		log.Fatal(err)
	}

	go messageRouter.Run(context.Background())
	<-messageRouter.Running()

	app.Logger.Info("Starting HTTP server", nil)

	if err := httpSubscriber.StartHTTPServer(); err != nil {
		app.Logger.Error("Could not start HTTP server", err, nil)
	}
}
