package main

import (
	"fmt"
	"log"
	"os"

	//_ "net/http/pprof"

	"github.com/Shelex/webhook-listener/api"
	"github.com/Shelex/webhook-listener/app"
	"github.com/joho/godotenv"
)

// @title webhook listener API
// @version 1.0
// @description webhook listener api
// @schemes https
// @host webhook-api.shelex.dev
// @BasePath /
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email ovr.shevtsov@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %s", err)
	}

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

	log.Println("Starting HTTP server")

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if err := app.Router.Listen(fmt.Sprintf("%s:%s", host, port)); err != nil {
		log.Printf("Could not start HTTP server %s:\n", err)
	}
}
