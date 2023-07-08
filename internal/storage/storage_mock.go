// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package storage is a generated GoMock package.
package storage

import (
	context "context"
	reflect "reflect"

	entities "github.com/akrillis/k8storage/internal/entities"
	gomock "github.com/golang/mock/gomock"
)

// MockFrontender is a mock of Frontender interface.
type MockFrontender struct {
	ctrl     *gomock.Controller
	recorder *MockFrontenderMockRecorder
}

// MockFrontenderMockRecorder is the mock recorder for MockFrontender.
type MockFrontenderMockRecorder struct {
	mock *MockFrontender
}

// NewMockFrontender creates a new mock instance.
func NewMockFrontender(ctrl *gomock.Controller) *MockFrontender {
	mock := &MockFrontender{ctrl: ctrl}
	mock.recorder = &MockFrontenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFrontender) EXPECT() *MockFrontenderMockRecorder {
	return m.recorder
}

// Put mocks base method.
func (m *MockFrontender) Put(ctx context.Context, req *entities.PutFileRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockFrontenderMockRecorder) Put(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockFrontender)(nil).Put), ctx, req)
}

// MockBackender is a mock of Backender interface.
type MockBackender struct {
	ctrl     *gomock.Controller
	recorder *MockBackenderMockRecorder
}

// MockBackenderMockRecorder is the mock recorder for MockBackender.
type MockBackenderMockRecorder struct {
	mock *MockBackender
}

// NewMockBackender creates a new mock instance.
func NewMockBackender(ctrl *gomock.Controller) *MockBackender {
	mock := &MockBackender{ctrl: ctrl}
	mock.recorder = &MockBackenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBackender) EXPECT() *MockBackenderMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockBackender) Get(ctx context.Context, req *entities.GetFileRequest) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, req)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockBackenderMockRecorder) Get(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockBackender)(nil).Get), ctx, req)
}