package main

import (
	"github.com/iTcatt/avito-task/internal/service"
	"log"
	"net/http"

	api "github.com/iTcatt/avito-task/internal/api/http"
	"github.com/iTcatt/avito-task/internal/config"
	"github.com/iTcatt/avito-task/internal/storage/postgres"
)

func main() {
	cfg := config.MustLoad()
	db, err := postgres.NewStorage(cfg.Storage)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.StartUp(); err != nil {
		log.Fatal(err)
	}

	serv := service.NewService(db)
	handler := api.NewHandler(serv)
	server := http.Server{
		Addr:    cfg.Server.Endpoint,
		Handler: api.NewRouter(handler),
	}

	log.Fatal(server.ListenAndServe())
}
