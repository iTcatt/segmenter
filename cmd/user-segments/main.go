package main

import (
	"fmt"
	"github.com/iTcatt/avito-task/internal/config"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iTcatt/avito-task/internal/http-server/handlers"
	"github.com/iTcatt/avito-task/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad()
	dbPath := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.User, cfg.Storage.Password, cfg.Storage.Name)

	db, err := postgres.NewPostgresStorage(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	err = db.StartUp()
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/api/segments", handlers.CreateSegmentsHandler(db))
	router.Delete("/api/segments", handlers.DeleteSegmentHandler(db))
	router.Get("/api/segments", handlers.GetUserSegmentsHandler(db))
	router.Post("/api/users", handlers.CreateUsersHandler(db))
	router.Delete("/api/users", handlers.DeleteUserHandler(db))
	router.Post("/api/update", handlers.UpdateUserHandler(db))

	address := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Fatal(http.ListenAndServe(address, router))
}
