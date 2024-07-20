package service

import (
	"context"
	"testing"

	"github.com/iTcatt/segmenter/internal/models"
	"github.com/iTcatt/segmenter/internal/service/mocks"
	"github.com/iTcatt/segmenter/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateSegments(t *testing.T) {

}

func TestService_GetUser(t *testing.T) {

	tests := []struct {
		name   string
		id     int
		result models.User
		err    error
	}{
		{
			name:   "no id",
			id:     0,
			result: models.User{},
			err:    storage.ErrNotCreated,
		},
		{
			name: "success",
			id:   1,
			result: models.User{
				ID:       1,
				Segments: []string{"a", "b", "c"},
			},
			err: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStorage := mocks.NewSegmentStorage(t)
			mockStorage.
				On("GetUser", mock.Anything, test.id).
				Return(test.result, test.err).
				Once()

			service := NewService(mockStorage)
			user, err := service.GetUser(context.Background(), test.id)
			assert.Equal(t, test.result, user)
			assert.Equal(t, test.err, err)
		})
	}
}
