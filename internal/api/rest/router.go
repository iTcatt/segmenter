package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/iTcatt/segmenter/docs"
	"github.com/swaggo/http-swagger"
)

func NewRouter(h *Handler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/api/user/{id}", errorsMiddleware(h.GetUser))
	router.Post("/api/user", errorsMiddleware(h.CreateUsers))
	router.Post("/api/segment", errorsMiddleware(h.CreateSegments))
	router.Patch("/api/user/{id}", errorsMiddleware(h.UpdateUser))
	router.Delete("/api/user/{id}", errorsMiddleware(h.DeleteUser))
	router.Delete("/api/segment/{name}", errorsMiddleware(h.DeleteSegment))

	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("swagger/doc.json")))

	return router
}
