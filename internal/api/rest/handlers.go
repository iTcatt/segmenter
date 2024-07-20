package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/iTcatt/segmenter/internal/models"
)

var ErrValidation = errors.New("validation error")

type SegmentService interface {
	CreateSegments(context.Context, []string) (map[string]string, error)
	CreateUsers(context.Context, []int) (map[int]string, error)

	GetUser(context.Context, int) (models.User, error)

	UpdateUser(context.Context, models.UpdateUserParams) (models.User, error)

	DeleteSegment(context.Context, string) error
	DeleteUser(context.Context, int) error
}

type Handler struct {
	service SegmentService
}

func NewHandler(s SegmentService) *Handler {
	return &Handler{
		service: s,
	}
}

func (h *Handler) CreateSegments(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var req struct {
		Segments []string `json:"segments"`
	}
	if err = json.Unmarshal(body, &req); err != nil {
		return err
	}
	log.Printf("CreateSegments request: %v", req)
	reply, err := h.service.CreateSegments(r.Context(), req.Segments)
	if err != nil {
		return err
	}
	return sendJSONResponse(w, reply, http.StatusCreated)
}

func (h *Handler) CreateUsers(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var req struct {
		Users []int `json:"users"`
	}
	if err = json.Unmarshal(body, &req); err != nil {
		return err
	}

	log.Printf("Create users request: %v", req)

	reply, err := h.service.CreateUsers(r.Context(), req.Users)
	if err != nil {
		return err
	}
	return sendJSONResponse(w, reply, http.StatusCreated)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	op := "UpdateUser:"

	id := chi.URLParam(r, "id")
	log.Printf("%s received user_id '%s'", op, id)

	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("%v: %v", ErrValidation, err)
		return ErrValidation
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	var req struct {
		AddSegments    []string `json:"add_segments"`
		DeleteSegments []string `json:"delete_segments"`
	}
	if err = json.Unmarshal(body, &req); err != nil {
		return err
	}
	log.Printf("%s received '%v'", op, req)

	user, err := h.service.UpdateUser(r.Context(), models.UpdateUserParams{
		ID:             userID,
		AddSegments:    req.AddSegments,
		DeleteSegments: req.DeleteSegments,
	})
	if err != nil {
		return err
	}

	return sendJSONResponse(w, user, http.StatusOK)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	op := "GetUserSegmentsHandler:"

	id := chi.URLParam(r, "id")
	log.Printf("%s received user_id '%s'", op, id)

	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("%v: %v", ErrValidation, err)
		return ErrValidation
	}

	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		return err
	}

	return sendJSONResponse(w, user, http.StatusOK)
}

func (h *Handler) DeleteSegment(w http.ResponseWriter, r *http.Request) error {
	op := "DeleteSegment"

	segment := chi.URLParam(r, "name")
	log.Printf("%s received segment name '%s'", op, segment)

	if err := h.service.DeleteSegment(r.Context(), segment); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	op := "DeleteUserHandler:"

	id := chi.URLParam(r, "id")
	log.Printf("%s received user_id '%s'", op, id)

	userID, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("%v: %v", ErrValidation, err)
		return ErrValidation
	}

	if err := h.service.DeleteUser(r.Context(), userID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
