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

	UpdateUser(context.Context, models.UpdateUserParams) error

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

// @Summary		CreateSegments
// @Description	Create segments
// @Tags			segment
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]string
// @Failure		500	{object}	ErrorResponse
// @Router			/segment [post]
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

// @Summary		CreateUser
// @Description	Create users
// @Tags			user
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[int]string
// @Failure		500	{object}	ErrorResponse
// @Router			/user [post]
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

// @Summary		UpdateUser
// @Description	Update user segments
// @Tags			user
// @Param			id	path	int	true	"userID"
// @Accept			json
// @Produce		json
// @Success		200	{object}	models.User
// @Failure		404
// @Failure		500	{object}	ErrorResponse
// @Router			/user/{id} [patch]
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

	err = h.service.UpdateUser(r.Context(), models.UpdateUserParams{
		ID:             userID,
		AddSegments:    req.AddSegments,
		DeleteSegments: req.DeleteSegments,
	})
	if err != nil {
		return err
	}
	user, err := h.service.GetUser(r.Context(), userID)
	if err != nil {
		return err
	}
	return sendJSONResponse(w, user, http.StatusOK)
}

// @Summary		GetUser
// @Description	get user segments
// @Tags		user
// @Param		id	path	int	true	"userID"
// @Produce		json
// @Success		200	{object}	models.User
// @Failure		404
// @Failure		500	{object}	ErrorResponse
// @Router		/user/{id} [get]
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

// @Summary		DeleteSegment
// @Description	delete segment
// @Tags		segment
// @Param		name	path	string	true	"segment name"
// @Success		204
// @Failure		404
// @Failure		500	{object}	ErrorResponse
// @Router		/segment/{name} [delete]
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

// @Summary		DeleteUser
// @Description	delete user
// @Tags		user
// @Param		id	path	int	true	"userID"
// @Success		204
// @Failure		404
// @Failure		500	{object}	ErrorResponse
// @Router		/user/{id} [delete]
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
