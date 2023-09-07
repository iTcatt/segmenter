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

	router.Post("/api/segments", handlers.CreateSegmentsHandler(db))
	router.Post("/api/users", handlers.CreateUsersHandler(db))
	router.Post("/api/update", handlers.UpdateUserHandler(db))
	log.Fatal(http.ListenAndServe(":3000", router))
}

