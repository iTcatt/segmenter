package service

import (
	"context"
	"errors"
	"log"

	"github.com/iTcatt/avito-task/internal/models"
	"github.com/iTcatt/avito-task/internal/storage"
)

type SegmentStorage interface {
	CreateSegment(ctx context.Context, name string) error
	CreateUser(ctx context.Context, id int) error
	AddUserToSegment(ctx context.Context, userID, segmentID int) error

	GetUser(ctx context.Context, id int) (models.User, error)
	GetSegmentIDByName(ctx context.Context, name string) (int, error)

	DeleteSegment(ctx context.Context, name string) error
	DeleteUser(ctx context.Context, id int) error
	DeleteUserFromSegment(ctx context.Context, userID, segmentID int) error
}

type Service struct {
	repo SegmentStorage
}

func NewService(repo SegmentStorage) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSegments(ctx context.Context, segments []string) (map[string]string, error) {
	reply := make(map[string]string)
	for _, segment := range segments {
		err := s.repo.CreateSegment(ctx, segment)
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
	return reply, nil
}

func (s *Service) CreateUsers(ctx context.Context, users []int) (map[int]string, error) {
	result := make(map[int]string)
	for _, userID := range users {
		err := s.repo.CreateUser(ctx, userID)
		switch {
		case errors.Is(err, storage.ErrAlreadyExist):
			if _, ok := result[userID]; !ok {
				result[userID] = "already exist"
			}
			log.Printf("EXIST: user '%d' already exist", userID)
		case err == nil:
			result[userID] = "created"
			log.Printf("SUCCESS: user '%d' was created", userID)
		default:
			result[userID] = "not created"
			log.Printf("ERROR: create user '%d' failed: %v", userID, err)
		}
	}
	return result, nil
}

func (s *Service) GetUserSegments(ctx context.Context, id int) ([]string, error) {
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		log.Printf("ERROR: get user '%d': %v", id, err)
		return nil, err
	}
	return user.Segments, nil
}

func (s *Service) UpdateUser(ctx context.Context, params models.UpdateUserParams) (models.User, error) {
	for _, segment := range params.AddSegments {
		segmentID, err := s.repo.GetSegmentIDByName(ctx, segment)
		if err != nil {
			log.Printf("ERROR: get segment id failed: %v", err)
			continue
		}
		if err := s.repo.AddUserToSegment(ctx, params.ID, segmentID); err != nil {
			log.Printf("ERROR: add user segment to segment failed: %v", err)
			return models.User{}, err
		}
	}

	for _, segment := range params.DeleteSegments {
		segmentID, err := s.repo.GetSegmentIDByName(ctx, segment)
		if err != nil {
			log.Printf("ERROR: get segment id failed: %v", err)
			continue
		}
		if err := s.repo.DeleteUserFromSegment(ctx, params.ID, segmentID); err != nil {
			log.Printf("ERROR: delete user segment from segment failed: %v", err)
			return models.User{}, err
		}
	}

	user, err := s.repo.GetUser(ctx, params.ID)
	if err != nil {
		log.Printf("ERROR: get user '%d': %v", params.ID, err)
		return models.User{}, err
	}
	return user, nil
}

func (s *Service) DeleteSegment(ctx context.Context, name string) error {
	err := s.repo.DeleteSegment(ctx, name)
	if err != nil {
		log.Printf("ERROR: delete segment '%v': %v", name, err)
		return err
	}
	return nil
}

func (s *Service) DeleteUser(ctx context.Context, id int) error {
	err := s.repo.DeleteUser(ctx, id)
	if err != nil {
		log.Printf("ERROR: delete user '%d': %v", id, err)
		return err
	}
	return nil
}
