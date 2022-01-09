package app

import (
	"net/http"
	"os"
)

// ServeWeb is serving static folder built from web page sources
func (app *App) ServeWeb() {
	root := "./web/build"
	fs := http.FileServer(http.Dir(root))

	app.Router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
