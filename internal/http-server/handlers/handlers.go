package handlers

import (
	"encoding/json"
	"errors"
	"github.com/iTcatt/avito-task/internal/http-server/requests"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/iTcatt/avito-task/internal/http-server/replies"
	"github.com/iTcatt/avito-task/internal/storage"
)

func CreateSegmentsHandler(s storage.Storage) http.HandlerFunc {
	op := "CreateSegmentsHandler:"
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("%s: ReadAll %v", op, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var req requests.CreateSegments
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Printf("%s unmarshal json %v", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("%s request %v", op, req)

		reply := make(map[string]string)
		for _, segment := range req.Segments {
			err := s.CreateSegment(segment)
			switch {
			case errors.Is(err, storage.ErrAlreadyExist):
				if _, ok := reply[segment]; !ok {
					reply[segment] = "already exist"
				}
				log.Printf("EXIST: segment '%s' already exist", segment)
			case err == nil:
				reply[segment] = "created"
				log.Printf("SUCCESS: segment '%s' was created", segment)
			default:
				reply[segment] = "not created"
				log.Printf("ERROR: create segment '%s' failed: %v\n", segment, err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(reply)
		if err != nil {
			log.Printf("%s encode error: %v", op, err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Printf("%s sending reply: %v", op, reply)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func CreateUsersHandler(s storage.Storage) http.HandlerFunc {
	op := "CreateUsersHandle:"
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("%s: ReadAll %v", op, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var req requests.CreateUsers
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Printf("%s unmarshal json %v", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Create users request: %v", req)

		reply := make(map[int]string)
		for _, userID := range req.Users {
			err := s.AddUser(userID)
			switch {
			case errors.Is(err, storage.ErrAlreadyExist):
				if _, ok := reply[userID]; !ok {
					reply[userID] = "already exist"
				}
				log.Printf("EXIST: user '%d' already exist", userID)
			case err == nil:
				reply[userID] = "created"
				log.Printf("SUCCESS: user '%d' was created", userID)
			default:
				reply[userID] = "not created"
				log.Printf("ERROR: create user '%d' failed: %v", userID, err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(reply)
		if err != nil {
			log.Printf("%s encode error: %v", op, err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Printf("%s sending reply %v", op, reply)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func UpdateUserHandler(s storage.Storage) http.HandlerFunc {
	op := "UpdateUserHandler:"
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("%s: ReadAll %v", op, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var req requests.UpdateUser
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Printf("%s unmarshal json %v", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Request: %v", req)
		w.Header().Set("Content-Type", "application/json")

		isUserCreate := true
		reply := replies.UpdateUser{
			ID:     req.Id,
			Add:    make(map[string]string, len(req.Add)),
			Delete: make(map[string]string, len(req.Delete)),
		}

		for _, segment := range req.Add {
			err = s.AddUserToSegment(req.Id, segment)
			switch {
			case errors.Is(err, storage.ErrAlreadyExist):
				log.Printf("EXIST: user '%d' is already in the segment '%s'", req.Id, segment)
				reply.Add[segment] = "already exist"
			case errors.Is(err, storage.ErrNotCreated):
				log.Printf("NOTCREATED: segment '%s' is not created", segment)
				reply.Add[segment] = "not created"
			case err == nil:
				log.Printf("SUCСESS: user '%d' has been added to the segment '%s'", req.Id, segment)
				reply.Add[segment] = "added"
			case errors.Is(err, storage.ErrNotExist):
				isUserCreate = false
				break
			default:
				log.Printf("ERROR: user '%d' is not added to the segment: %v", req.Id, err)
				reply.Add[segment] = "error"
			}
		}

		if !isUserCreate {
			reply.ID = -1
			err := json.NewEncoder(w).Encode(reply)
			log.Printf("%s user '%d' not created", op, req.Id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Printf("%s encode error %v", op, err)
			} else {
				w.WriteHeader(http.StatusBadRequest)
				log.Printf("%s sending reply %v", op, reply)
			}
			return
		}

		for _, segment := range req.Delete {
			err = s.DeleteUserFromSegment(req.Id, segment)
			switch {
			case errors.Is(err, storage.ErrNotExist):
				log.Printf("NOTEXIST: user '%d' is not a member of the segment '%s'", req.Id, segment)
				reply.Delete[segment] = "not exist"
			case errors.Is(err, storage.ErrNotCreated):
				log.Printf("NOTCREATED: segment '%s' is not created", segment)
				reply.Delete[segment] = "not created"
			case err == nil:
				log.Printf("SUCСESS: user '%d' has been removed from the segment '%s'", req.Id, segment)
				reply.Delete[segment] = "removed"
			default:
				log.Printf("ERROR: user '%d' is not removed from the segment: %v", req.Id, err)
				reply.Delete[segment] = "error"
			}
		}

		err = json.NewEncoder(w).Encode(reply)
		if err != nil {
			log.Printf("%s encode error: %v", op, err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Printf("%s sending reply %v", op, reply)
			w.WriteHeader(http.StatusOK)
		}

	}
}

func GetUserSegmentsHandler(s storage.Storage) http.HandlerFunc {
	op := "GetUserSegmentsHandler:"
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			log.Printf("%s invalid query param %v", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Printf("%s received user_id '%d'", op, userID)
		w.Header().Set("Content-Type", "application/json")

		user, err := s.GetUserSegments(userID)
		switch {
		case errors.Is(err, storage.ErrNotExist):
			log.Printf("NOTEXIST: user '%d' is not contained in any segment", userID)
			w.WriteHeader(http.StatusNotFound)
		case err == nil:
			log.Printf("SUCСESS: user '%d' is in the segments: '%v'", userID, user.Segments)
			w.WriteHeader(http.StatusOK)
		default:
			log.Printf("ERROR: user '%d' segment data cannot be retrieved: %v", userID, err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			log.Printf("%s encode error: %v", op, err)
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			log.Printf("%s sending reply %v", op, user)
		}
	}
}

func DeleteSegmentHandler(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		segment := r.URL.Query().Get("name")
		err := s.DeleteSegment(segment)
		switch {
		case errors.Is(err, storage.ErrNotExist):
			log.Printf("NOTEXIST: segment '%s' has not been created", segment)
			w.WriteHeader(http.StatusNotFound)
		case err == nil:
			log.Printf("SUCCESS: segment '%s' successfully deleted", segment)
			w.WriteHeader(http.StatusOK)
		default:
			log.Printf("ERROR: segment '%s' has not been deleted: %v", segment, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func DeleteUserHandler(s storage.Storage) http.HandlerFunc {
	op := "DeleteUserHandler:"
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Printf("%s invalid query param %v", op, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = s.DeleteUser(userID)
		switch {
		case errors.Is(err, storage.ErrNotExist):
			log.Printf("NOTEXIST: user '%d' has not been created", userID)
			w.WriteHeader(http.StatusNotFound)
		case err == nil:
			log.Printf("SUCCESS: user '%d' successfully deleted", userID)
			w.WriteHeader(http.StatusOK)
		default:
			log.Printf("ERROR: user '%d' has not been deleted: %v", userID, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
