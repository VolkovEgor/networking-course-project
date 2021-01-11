// Code generated by MockGen. DO NOT EDIT.
// Source: yak/backend/pkg/repositories (interfaces: TaskList)

// Package mock_repositories is a generated GoMock package.
package mock_repositories

import (
	models "github.com/architectv/networking-course-project/backend/pkg/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTaskList is a mock of TaskList interface.
type MockTaskList struct {
	ctrl     *gomock.Controller
	recorder *MockTaskListMockRecorder
}

// MockTaskListMockRecorder is the mock recorder for MockTaskList.
type MockTaskListMockRecorder struct {
	mock *MockTaskList
}

// NewMockTaskList creates a new mock instance.
func NewMockTaskList(ctrl *gomock.Controller) *MockTaskList {
	mock := &MockTaskList{ctrl: ctrl}
	mock.recorder = &MockTaskListMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskList) EXPECT() *MockTaskListMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTaskList) Create(arg0 *models.TaskList) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockTaskListMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTaskList)(nil).Create), arg0)
}

// Delete mocks base method.
func (m *MockTaskList) Delete(arg0 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTaskListMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTaskList)(nil).Delete), arg0)
}

// GetAll mocks base method.
func (m *MockTaskList) GetAll(arg0 int) ([]*models.TaskList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]*models.TaskList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockTaskListMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockTaskList)(nil).GetAll), arg0)
}

// GetById mocks base method.
func (m *MockTaskList) GetById(arg0 int) (*models.TaskList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", arg0)
	ret0, _ := ret[0].(*models.TaskList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockTaskListMockRecorder) GetById(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockTaskList)(nil).GetById), arg0)
}

// Update mocks base method.
func (m *MockTaskList) Update(arg0 int, arg1 *models.UpdateTaskList) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockTaskListMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockTaskList)(nil).Update), arg0, arg1)
}
