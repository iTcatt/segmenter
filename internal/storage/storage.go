package storage

import (
	"errors"

	"github.com/iTcatt/avito-task/internal/types"
)

type Storage interface {
	CreateSegment(name string) error
	DeleteSegment(name string) error
	AddUser(id int) error
	DeleteUser(id int) error
	AddUserToSegment(id int, segment string) error
	DeleteUserFromSegment(id int, segment string) error
	GetUserSegments(id int) (types.User, error)
}

var ErrAlreadyExist = errors.New("already exist")
var ErrNotExist = errors.New("not exist")
