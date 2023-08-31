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
	return func(w http.ResponseWriter, r *http.Request){
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("ReadAll")
		}

		var req types.CreateRequest
		err = json.Unmarshal(body, &req)
		if err != nil {
			log.Printf("%s unmarshal json %v", op, err)
		}
		log.Printf("Request: %v", req)

		create := make(map[string]string)
		
		for _, segment := range req.Segments {
			err := s.CreateSegment(segment)
			if err != nil {
				create[segment] = "not created"
				log.Printf("Create segment '%s' failed: %v\n", segment, err)
			} else {
				create[segment] = "created"
				log.Printf("Segment '%v' was created", segment)
			}
		}
		resp, err := json.Marshal(create)
		if err != nil {
			log.Printf("%s Responce not json %v", op, err)
		}
		_, err = w.Write(resp)
		if err != nil {
			log.Printf("%s Responce not write %v",op, err)
		}
	}
} 

