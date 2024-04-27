package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handler) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/api/user/{id}", errorsMiddleware(h.GetUserSegments))
	router.Post("/api/user", errorsMiddleware(h.CreateUsers))
	router.Post("/api/segment", errorsMiddleware(h.CreateSegments))
	router.Patch("/api/user/{id}", errorsMiddleware(h.UpdateUser))
	router.Delete("/api/user/{id}", errorsMiddleware(h.DeleteUser))
	router.Delete("/api/segment/{name}", errorsMiddleware(h.DeleteSegment))

	return router
}
