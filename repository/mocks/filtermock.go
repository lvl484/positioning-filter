// Code generated by MockGen. DO NOT EDIT.
// Source: filtering.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	repository "github.com/lvl484/positioning-filter/repository"
)

// MockFilters is a mock of Filters interface
type MockFilters struct {
	ctrl     *gomock.Controller
	recorder *MockFiltersMockRecorder
}

// MockFiltersMockRecorder is the mock recorder for MockFilters
type MockFiltersMockRecorder struct {
	mock *MockFilters
}

// NewMockFilters creates a new mock instance
func NewMockFilters(ctrl *gomock.Controller) *MockFilters {
	mock := &MockFilters{ctrl: ctrl}
	mock.recorder = &MockFiltersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFilters) EXPECT() *MockFiltersMockRecorder {
	return m.recorder
}

// Add mocks base method
func (m *MockFilters) Add(filter *repository.Filter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", filter)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockFiltersMockRecorder) Add(filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockFilters)(nil).Add), filter)
}

// OneByUser mocks base method
func (m *MockFilters) OneByUser(userID uuid.UUID, filterName string) (*repository.Filter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OneByUser", userID, filterName)
	ret0, _ := ret[0].(*repository.Filter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OneByUser indicates an expected call of OneByUser
func (mr *MockFiltersMockRecorder) OneByUser(userID, filterName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OneByUser", reflect.TypeOf((*MockFilters)(nil).OneByUser), userID, filterName)
}

// AllByUser mocks base method
func (m *MockFilters) AllByUser(userID uuid.UUID) ([]*repository.Filter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllByUser", userID)
	ret0, _ := ret[0].([]*repository.Filter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllByUser indicates an expected call of AllByUser
func (mr *MockFiltersMockRecorder) AllByUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllByUser", reflect.TypeOf((*MockFilters)(nil).AllByUser), userID)
}

// OffsetByUser mocks base method
func (m *MockFilters) OffsetByUser(userID uuid.UUID, offset int) ([]*repository.Filter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OffsetByUser", userID, offset)
	ret0, _ := ret[0].([]*repository.Filter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OffsetByUser indicates an expected call of OffsetByUser
func (mr *MockFiltersMockRecorder) OffsetByUser(userID, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OffsetByUser", reflect.TypeOf((*MockFilters)(nil).OffsetByUser), userID, offset)
}

// Update mocks base method
func (m *MockFilters) Update(filter *repository.Filter) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", filter)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockFiltersMockRecorder) Update(filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockFilters)(nil).Update), filter)
}

// Delete mocks base method
func (m *MockFilters) Delete(userID uuid.UUID, filterName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", userID, filterName)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockFiltersMockRecorder) Delete(userID, filterName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockFilters)(nil).Delete), userID, filterName)
}
