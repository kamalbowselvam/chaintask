// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/kamalbowselvam/chaintask/db (interfaces: GlobalRepository)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	db "github.com/kamalbowselvam/chaintask/db"
	domain "github.com/kamalbowselvam/chaintask/domain"
)

// MockGlobalRepository is a mock of GlobalRepository interface.
type MockGlobalRepository struct {
	ctrl     *gomock.Controller
	recorder *MockGlobalRepositoryMockRecorder
}

// MockGlobalRepositoryMockRecorder is the mock recorder for MockGlobalRepository.
type MockGlobalRepositoryMockRecorder struct {
	mock *MockGlobalRepository
}

// NewMockGlobalRepository creates a new mock instance.
func NewMockGlobalRepository(ctrl *gomock.Controller) *MockGlobalRepository {
	mock := &MockGlobalRepository{ctrl: ctrl}
	mock.recorder = &MockGlobalRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGlobalRepository) EXPECT() *MockGlobalRepositoryMockRecorder {
	return m.recorder
}

// CreateTask mocks base method.
func (m *MockGlobalRepository) CreateTask(arg0 context.Context, arg1 db.CreateTaskParams) (domain.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", arg0, arg1)
	ret0, _ := ret[0].(domain.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask.
func (mr *MockGlobalRepositoryMockRecorder) CreateTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockGlobalRepository)(nil).CreateTask), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockGlobalRepository) CreateUser(arg0 context.Context, arg1 domain.User) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockGlobalRepositoryMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockGlobalRepository)(nil).CreateUser), arg0, arg1)
}

// DeleteTask mocks base method.
func (m *MockGlobalRepository) DeleteTask(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTask", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTask indicates an expected call of DeleteTask.
func (mr *MockGlobalRepositoryMockRecorder) DeleteTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTask", reflect.TypeOf((*MockGlobalRepository)(nil).DeleteTask), arg0, arg1)
}

// GetTask mocks base method.
func (m *MockGlobalRepository) GetTask(arg0 context.Context, arg1 int64) (domain.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTask", arg0, arg1)
	ret0, _ := ret[0].(domain.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTask indicates an expected call of GetTask.
func (mr *MockGlobalRepositoryMockRecorder) GetTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTask", reflect.TypeOf((*MockGlobalRepository)(nil).GetTask), arg0, arg1)
}

// GetTaskList mocks base method.
func (m *MockGlobalRepository) GetTaskList(arg0 context.Context, arg1 []int64) ([]domain.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTaskList", arg0, arg1)
	ret0, _ := ret[0].([]domain.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskList indicates an expected call of GetTaskList.
func (mr *MockGlobalRepositoryMockRecorder) GetTaskList(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTaskList", reflect.TypeOf((*MockGlobalRepository)(nil).GetTaskList), arg0, arg1)
}

// UpdateTask mocks base method.
func (m *MockGlobalRepository) UpdateTask(arg0 context.Context, arg1 domain.Task) (domain.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTask", arg0, arg1)
	ret0, _ := ret[0].(domain.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateTask indicates an expected call of UpdateTask.
func (mr *MockGlobalRepositoryMockRecorder) UpdateTask(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTask", reflect.TypeOf((*MockGlobalRepository)(nil).UpdateTask), arg0, arg1)
}
