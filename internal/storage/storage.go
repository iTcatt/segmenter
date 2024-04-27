package storage

import (
	"errors"
)

var ErrAlreadyExist = errors.New("already exist")
var ErrNotExist = errors.New("not exist")
var ErrNotCreated = errors.New("not created")
