package main

import (
	"log"
	"net/http"

	"github.com/iTcatt/segmenter/internal/service"

	"github.com/iTcatt/segmenter/internal/api/rest"
	"github.com/iTcatt/segmenter/internal/config"
	"github.com/iTcatt/segmenter/internal/storage/postgres"
)

// @title			segmenter
// @version		1.0
// @description	REST API server for saving users and their segments
// @host			localhost:3000
// @BasePath		/api
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
	handler := rest.NewHandler(serv)
	server := http.Server{
		Addr:    cfg.Server.Endpoint,
		Handler: rest.NewRouter(handler),
	}

	log.Fatal(server.ListenAndServe())
}
