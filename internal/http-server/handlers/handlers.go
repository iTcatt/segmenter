package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/iTcatt/avito-task/internal/storage"
	"github.com/iTcatt/avito-task/internal/types"
)

func CreateSegmentsHandler(s storage.Storage) http.HandlerFunc {
	op := "CreateSegmentsHandler:"
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("%s: ReadAll %v", op, err)
		}

		var req types.CreateSegmentsRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Printf("%s unmarshal json %v", op, err)
		}
		log.Printf("%s request %v", op, req)

		reply := make(map[string]string)

		for _, segment := range req.Segments {
			err := s.CreateSegment(segment)
			switch err {
			case storage.ErrAlreadyExist:
				if _, ok := reply[segment]; !ok {
					reply[segment] = "already exist"
				}
				log.Printf("EXIST: segment '%s' already exist", segment)
			case nil:
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
			log.Fatalf("%s: ReadAll %v", op, err)
		}

		var req types.CreateUsersRequest
		json.Unmarshal(body, &req)
		log.Printf("Create users request: %v", req)
		reply := make(map[int]string)
		for _, userID := range req.Users {
			err := s.AddUser(userID)
			switch err {
			case storage.ErrAlreadyExist:
				if _, ok := reply[userID]; !ok {
					reply[userID] = "already exist"
				}
				log.Printf("EXIST: user '%d' already exist", userID)
			case nil:
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
			log.Fatalf("%s: ReadAll %v", op, err)
		}
		var req types.UpdateUser
		json.Unmarshal(body, &req)
		log.Printf("Request: %v", req)
		// reply := make(map[int]string)
		for _, segment := range req.Add {
			err = s.AddUserToSegment(req.Id, segment)
			switch err {
			case storage.ErrAlreadyExist:
				log.Printf("EXIST: user '%d' is already in the segment '%s'", req.Id, segment)
			case nil:
				log.Printf("SUCСESS: user '%d' has been added to the segment '%s'", req.Id, segment)
			default:
				log.Printf("ERROR: user '%d' is not added to the segment: %v", req.Id, err)
			}
		}

		for _, segment := range req.Delete {
			err = s.DeleteUserFromSegment(req.Id, segment)
			switch err {
			case storage.ErrNotExist:
				log.Printf("NOTEXIST: user '%d' is not a member of the segment '%s'", req.Id, segment)
			case nil:
				log.Printf("SUCСESS: user '%d' has been removed from the segment '%s'", req.Id, segment)
			default:
				log.Printf("ERROR: user '%d' is not removed from the segment: %v", req.Id, err)
			}
		}

	}
}

func GetUserSegmentsHandler(s storage.Storage) http.HandlerFunc {
	op := "GetUserSegmentsHandler:"
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
		if err != nil {
			log.Fatalf("%s invalid query param", op)
		}
		log.Printf("Return the names of the segments that the user '%d' is a member of", userID)
		w.Header().Set("Content-Type", "application/json")

		user, err := s.GetUserSegments(userID)
		switch err {
		case storage.ErrNotExist:
			log.Printf("NOTEXIST: user '%d' is not contained in any segment", userID)
			w.WriteHeader(http.StatusNotFound)
		case nil:
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
			w.WriteHeader(http.StatusOK)
		}
	}
}

func DeleteSegmentHandler(s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		segment := r.URL.Query().Get("name")
		err := s.DeleteSegment(segment)
		switch err {
		case storage.ErrNotExist:
			log.Printf("NOTEXIST: segment '%s' has not been created", segment)
			w.WriteHeader(http.StatusNotFound)
		case nil:
			log.Printf("SUCCESS: segment '%s' successfully deleted", segment)
			w.WriteHeader(http.StatusOK)
		default:
			log.Printf("ERROR: segment '%s' has not been deleted: %v", segment, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
