package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

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
		log.Printf("Request: %v", req)

		reply := make(map[string]string)

		for _, segment := range req.Segments {
			err := s.CreateSegment(segment)
			switch err {
			case storage.ErrAlreadyExist:
				if _, ok := reply[segment]; !ok {
					reply[segment] = "already exist"
				}
				log.Printf("Segment '%v' already exist", segment)
			case nil:
				reply[segment] = "created"
				log.Printf("Segment '%v' was created", segment)
			default:
				reply[segment] = "not created"
				log.Printf("Create segment '%s' failed: %v\n", segment, err)
			}
		}
		resp, err := json.Marshal(reply)
		if err != nil {
			log.Printf("%s Responce not json %v", op, err)
		}
		_, err = w.Write(resp)
		if err != nil {
			log.Printf("%s Responce not write %v", op, err)
		}
	}
}

func CreateUsersHandler(s storage.Storage) http.HandlerFunc {
	op := "CreateUsersHandle:"
	_ = op
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatalf("%s: ReadAll %v", op, err)
		}
		var req types.CreateUsersRequest
		json.Unmarshal(body, &req)

		reply := make(map[int]string)
		for _, user_id := range req.Users {
			err := s.AddUser(user_id)
			switch err {
			case storage.ErrAlreadyExist:
				if _, ok := reply[user_id]; !ok {
					reply[user_id] = "already exist"
				}
				log.Printf("User '%d' already exist", user_id)
			case nil:
				reply[user_id] = "created"
				log.Printf("User '%d' was created", user_id)
			default:
				reply[user_id] = "not created"
				log.Printf("Create user '%d' failed: %v", user_id, err)
			}
		}

		resp, err := json.Marshal(reply)
		if err != nil {
			log.Printf("%s Responce not json %v", op, err)
		}
		_, err = w.Write(resp)
		if err != nil {
			log.Printf("%s Responce not write %v", op, err)
		}
	}
}
