package app

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func ProvideRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	//router.Use(middleware.RequestID)
	//router.Use(middleware.RealIP)
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Access-Control-Allow-Origin", "Content-Type", "Authorization"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}).Handler)

	return router
}
