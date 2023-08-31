package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iTcatt/avito-task/internal/http-server/handlers"
	"github.com/iTcatt/avito-task/internal/storage/postgres"
)

func main() {
	db, err := postgres.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}
	err = db.StartUp()
	if err != nil {
		log.Println(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/api/create", handlers.CreateSegmentsHandler(db))

	log.Fatal(http.ListenAndServe(":3000", router))
}
