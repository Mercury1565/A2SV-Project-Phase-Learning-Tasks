// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	domain "Task_8-Testing_Task_Management_REST_API/domain"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// TaskRepository is an autogenerated mock type for the TaskRepository type
type TaskRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: c, task
func (_m *TaskRepository) Create(c context.Context, task *domain.Task) error {
	ret := _m.Called(c, task)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Task) error); ok {
		r0 = rf(c, task)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteTask provides a mock function with given fields: c, taskID
func (_m *TaskRepository) DeleteTask(c context.Context, taskID string) error {
	ret := _m.Called(c, taskID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(c, taskID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetTaskByID provides a mock function with given fields: c, taskID
func (_m *TaskRepository) GetTaskByID(c context.Context, taskID string) (domain.Task, error) {
	ret := _m.Called(c, taskID)

	var r0 domain.Task
	if rf, ok := ret.Get(0).(func(context.Context, string) domain.Task); ok {
		r0 = rf(c, taskID)
	} else {
		r0 = ret.Get(0).(domain.Task)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(c, taskID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTasks provides a mock function with given fields: c
func (_m *TaskRepository) GetTasks(c context.Context) ([]domain.Task, error) {
	ret := _m.Called(c)

	var r0 []domain.Task
	if rf, ok := ret.Get(0).(func(context.Context) []domain.Task); ok {
		r0 = rf(c)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(c)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateTask provides a mock function with given fields: c, taskID, updated_task
func (_m *TaskRepository) UpdateTask(c context.Context, taskID string, updated_task *domain.Task) error {
	ret := _m.Called(c, taskID, updated_task)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *domain.Task) error); ok {
		r0 = rf(c, taskID, updated_task)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTaskRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewTaskRepository creates a new instance of TaskRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTaskRepository(t mockConstructorTestingTNewTaskRepository) *TaskRepository {
	mock := &TaskRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
