package storage

import "github.com/iTcatt/avito-task/internal/types"

type Storage interface {
	CreateSegment(name string) error
	DeleteSegment(name string) error
	AddUser(id int, addedSegments []string, removedSegments []string) error
	GetSegments(id int) (types.User, error)
}
