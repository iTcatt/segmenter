package rest

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/iTcatt/segmenter/internal/storage"
)

type wrapperHandler func(w http.ResponseWriter, r *http.Request) error

type ErrorResponse struct {
	Message string `json:"message"`
}

func errorsMiddleware(h wrapperHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)

		switch {
		case err == nil:
			return
		case errors.Is(err, storage.ErrNotCreated):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, storage.ErrNotExist):
			w.WriteHeader(http.StatusNotFound)
		case errors.Is(err, ErrValidation):
			_ = sendJSONResponse(w, ErrorResponse{Message: err.Error()}, http.StatusBadRequest)
		default:
			_ = sendJSONResponse(w, ErrorResponse{Message: "internal error"}, http.StatusInternalServerError)
		}
	}
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, status int) error {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return err
	}
	w.WriteHeader(status)

	log.Printf("sending response: %v", data)
	return nil
}
