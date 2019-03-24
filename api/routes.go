package main

import (
	"github.com/go-chi/chi"
)

func routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get("/user/{userID}", GetUser)
	router.Delete("/user/{rowID}", DeleteUser)
	router.Post("/user", UpdateUser)

	router.Get("/account/{userID}", GetAccount)
	router.Delete("/account/{rowID}", DeleteAccount)
	router.Post("/account", UpdateAccount)

	router.Get("/device/{userID}", GetDevice)
	router.Delete("/device/{rowID}", DeleteDevice)
	router.Post("/device", UpdateDevice)

	return router
}
