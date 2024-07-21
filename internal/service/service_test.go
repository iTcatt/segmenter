package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/iTcatt/segmenter/internal/models"
	"github.com/iTcatt/segmenter/internal/service/mocks"
	"github.com/iTcatt/segmenter/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetUser(t *testing.T) {
	ctx := context.Background()

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
			err:    storage.ErrNotExist,
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
			user, err := service.GetUser(ctx, test.id)
			assert.Equal(t, test.result, user)
			assert.Equal(t, test.err, err)
		})
	}
}

func TestService_CreateSegments(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name     string
		segments []string
		results  []error
		expected map[string]string
	}{
		{
			name:     "no segments",
			segments: []string{},
			results:  []error{},
			expected: map[string]string{},
		},
		{
			name:     "success",
			segments: []string{"a", "b"},
			results:  []error{nil, nil},
			expected: map[string]string{"a": "created", "b": "created"},
		},
		{
			name:     "two errors",
			segments: []string{"a", "b"},
			results:  []error{storage.ErrAlreadyExist, sql.ErrTxDone},
			expected: map[string]string{"a": "already exist", "b": "not created"},
		},
		{
			name:     "one success, one error",
			segments: []string{"a", "b"},
			results:  []error{storage.ErrAlreadyExist, nil},
			expected: map[string]string{"a": "already exist", "b": "created"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStorage := mocks.NewSegmentStorage(t)
			for i, segment := range test.segments {
				mockStorage.
					On("CreateSegment", mock.Anything, segment).
					Return(test.results[i])
			}
			service := NewService(mockStorage)
			result, err := service.CreateSegments(ctx, test.segments)
			assert.Equal(t, test.expected, result)
			assert.Nil(t, err)
		})
	}
}

func TestService_CreateUsers(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name     string
		users    []int
		results  []error
		expected map[int]string
	}{
		{
			name:     "no users",
			users:    []int{},
			results:  []error{},
			expected: map[int]string{},
		},
		{
			name:     "success",
			users:    []int{1, 2, 3},
			results:  []error{nil, nil, nil},
			expected: map[int]string{1: "created", 2: "created", 3: "created"},
		},
		{
			name:     "errors",
			users:    []int{4, 5, 6},
			results:  []error{storage.ErrAlreadyExist, sql.ErrTxDone, storage.ErrAlreadyExist},
			expected: map[int]string{4: "already exist", 5: "not created", 6: "already exist"},
		},
		{
			name:     "one success, one error",
			users:    []int{7, 8},
			results:  []error{storage.ErrAlreadyExist, nil},
			expected: map[int]string{7: "already exist", 8: "created"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStorage := mocks.NewSegmentStorage(t)
			for i, userID := range test.users {
				mockStorage.
					On("CreateUser", mock.Anything, userID).
					Return(test.results[i])
			}
			service := NewService(mockStorage)
			result, err := service.CreateUsers(ctx, test.users)
			assert.Equal(t, test.expected, result)
			assert.Nil(t, err)
		})
	}
}

func TestService_DeleteSegment(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		segment  string
		result   error
		expected error
	}{
		{
			name:     "error",
			segment:  "not exist",
			result:   storage.ErrNotExist,
			expected: storage.ErrNotExist,
		},
		{
			name:     "success delete user",
			segment:  "success",
			result:   nil,
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStorage := mocks.NewSegmentStorage(t)
			mockStorage.
				On("DeleteSegment", mock.Anything, test.segment).
				Return(test.result).
				Once()
			service := NewService(mockStorage)
			err := service.DeleteSegment(ctx, test.segment)
			assert.Equal(t, test.expected, err)
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		id       int
		result   error
		expected error
	}{
		{
			name:     "error",
			id:       0,
			result:   storage.ErrNotExist,
			expected: storage.ErrNotExist,
		},
		{
			name:     "success delete user",
			id:       1,
			result:   nil,
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStorage := mocks.NewSegmentStorage(t)
			mockStorage.
				On("DeleteUser", mock.Anything, test.id).
				Return(test.result).
				Once()
			service := NewService(mockStorage)
			err := service.DeleteUser(ctx, test.id)
			assert.Equal(t, test.expected, err)
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	ctx := context.Background()
	type isUserCreatedResults struct {
		created bool
		err     error
	}
	type resultFromDB struct {
		isUserCreated       isUserCreatedResults
		addUserToSegment    []error
		deleteUserToSegment []error
	}

	var tests = []struct {
		name     string
		params   models.UpdateUserParams
		result   resultFromDB
		expected error
	}{
		{
			name: "user not created",
			params: models.UpdateUserParams{
				ID: 0,
			},
			result: resultFromDB{
				isUserCreated: isUserCreatedResults{
					created: false,
					err:     nil,
				},
				addUserToSegment:    nil,
				deleteUserToSegment: nil,
			},
			expected: storage.ErrNotExist,
		},
		{
			name: "error find user",
			params: models.UpdateUserParams{
				ID: 1,
			},
			result: resultFromDB{
				isUserCreated: isUserCreatedResults{
					created: false,
					err:     sql.ErrNoRows,
				},
				addUserToSegment:    nil,
				deleteUserToSegment: nil,
			},
			expected: sql.ErrNoRows,
		},
		{
			name: "empty addSegment and deleteSegment",
			params: models.UpdateUserParams{
				ID:             1,
				AddSegments:    nil,
				DeleteSegments: nil,
			},
			result: resultFromDB{
				isUserCreated: isUserCreatedResults{
					created: true,
					err:     nil,
				},
			},
			expected: nil,
		},
		{
			name: "success add user to segments",
			params: models.UpdateUserParams{
				ID:             2,
				AddSegments:    []string{"a", "b"},
				DeleteSegments: nil,
			},
			result: resultFromDB{
				isUserCreated: isUserCreatedResults{
					created: true,
					err:     nil,
				},
				addUserToSegment:    []error{nil, nil},
				deleteUserToSegment: nil,
			},
			expected: nil,
		},
		{
			name: "success delete user from segments",
			params: models.UpdateUserParams{
				ID:             3,
				AddSegments:    nil,
				DeleteSegments: []string{"c", "d"},
			},
			result: resultFromDB{
				isUserCreated: isUserCreatedResults{
					created: true,
					err:     nil,
				},
				addUserToSegment:    nil,
				deleteUserToSegment: []error{nil, nil},
			},
			expected: nil,
		},
		{
			name: "unexpected error add user to segments",
			params: models.UpdateUserParams{
				ID:             4,
				AddSegments:    []string{"a"},
				DeleteSegments: nil,
			},
			result: resultFromDB{
				isUserCreated: isUserCreatedResults{
					created: true,
					err:     nil,
				},
				addUserToSegment:    []error{sql.ErrNoRows},
				deleteUserToSegment: nil,
			},
			expected: sql.ErrNoRows,
		},
		{
			name: "unexpected error: delete user to segments",
			params: models.UpdateUserParams{
				ID:             5,
				AddSegments:    []string{"c"},
				DeleteSegments: []string{"d"},
			},
			result: resultFromDB{
				isUserCreated: isUserCreatedResults{
					created: true,
					err:     nil,
				},
				addUserToSegment:    []error{storage.ErrAlreadyExist},
				deleteUserToSegment: []error{sql.ErrNoRows},
			},
			expected: sql.ErrNoRows,
		},
		{
			name: "not exist segments",
			params: models.UpdateUserParams{
				ID:             6,
				AddSegments:    []string{"a"},
				DeleteSegments: []string{"b"},
			},
			result: resultFromDB{
				isUserCreated: isUserCreatedResults{
					created: true,
					err:     nil,
				},
				addUserToSegment:    []error{storage.ErrNotExist},
				deleteUserToSegment: []error{storage.ErrNotExist},
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockStorage := mocks.NewSegmentStorage(t)
			mockStorage.
				On("IsUserCreated", mock.Anything, test.params.ID).
				Return(test.result.isUserCreated.created, test.result.isUserCreated.err).
				Once()
			for i, segment := range test.params.AddSegments {
				mockStorage.
					On("AddUserToSegment", mock.Anything, test.params.ID, segment).
					Return(test.result.addUserToSegment[i])
			}
			for i, segment := range test.params.DeleteSegments {
				mockStorage.
					On("DeleteUserFromSegment", mock.Anything, test.params.ID, segment).
					Return(test.result.deleteUserToSegment[i])
			}

			service := NewService(mockStorage)
			err := service.UpdateUser(ctx, test.params)
			assert.Equal(t, test.expected, err)
		})

	}
}
