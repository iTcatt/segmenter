// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/iTcatt/segmenter/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// SegmentStorage is an autogenerated mock type for the SegmentStorage type
type SegmentStorage struct {
	mock.Mock
}

// AddUserToSegment provides a mock function with given fields: ctx, userID, segment
func (_m *SegmentStorage) AddUserToSegment(ctx context.Context, userID int, segment string) error {
	ret := _m.Called(ctx, userID, segment)

	if len(ret) == 0 {
		panic("no return value specified for AddUserToSegment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) error); ok {
		r0 = rf(ctx, userID, segment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateSegment provides a mock function with given fields: ctx, name
func (_m *SegmentStorage) CreateSegment(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for CreateSegment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateUser provides a mock function with given fields: ctx, id
func (_m *SegmentStorage) CreateUser(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSegment provides a mock function with given fields: ctx, name
func (_m *SegmentStorage) DeleteSegment(ctx context.Context, name string) error {
	ret := _m.Called(ctx, name)

	if len(ret) == 0 {
		panic("no return value specified for DeleteSegment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, name)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUser provides a mock function with given fields: ctx, id
func (_m *SegmentStorage) DeleteUser(ctx context.Context, id int) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUserFromSegment provides a mock function with given fields: ctx, userID, segment
func (_m *SegmentStorage) DeleteUserFromSegment(ctx context.Context, userID int, segment string) error {
	ret := _m.Called(ctx, userID, segment)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUserFromSegment")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int, string) error); ok {
		r0 = rf(ctx, userID, segment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUser provides a mock function with given fields: ctx, id
func (_m *SegmentStorage) GetUser(ctx context.Context, id int) (models.User, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetUser")
	}

	var r0 models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (models.User, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) models.User); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.User)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsUserCreated provides a mock function with given fields: ctx, userID
func (_m *SegmentStorage) IsUserCreated(ctx context.Context, userID int) (bool, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for IsUserCreated")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) (bool, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) bool); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSegmentStorage creates a new instance of SegmentStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSegmentStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *SegmentStorage {
	mock := &SegmentStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
