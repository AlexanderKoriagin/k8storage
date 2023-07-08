// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	entities "github.com/akrillis/k8storage/internal/entities"
	gomock "github.com/golang/mock/gomock"
)

// MockReceiver is a mock of Receiver interface.
type MockReceiver struct {
	ctrl     *gomock.Controller
	recorder *MockReceiverMockRecorder
}

// MockReceiverMockRecorder is the mock recorder for MockReceiver.
type MockReceiverMockRecorder struct {
	mock *MockReceiver
}

// NewMockReceiver creates a new mock instance.
func NewMockReceiver(ctrl *gomock.Controller) *MockReceiver {
	mock := &MockReceiver{ctrl: ctrl}
	mock.recorder = &MockReceiverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReceiver) EXPECT() *MockReceiverMockRecorder {
	return m.recorder
}

// Put mocks base method.
func (m *MockReceiver) Put(ctx context.Context, req *entities.PutFileRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Put indicates an expected call of Put.
func (mr *MockReceiverMockRecorder) Put(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockReceiver)(nil).Put), ctx, req)
}

// MockRestorer is a mock of Restorer interface.
type MockRestorer struct {
	ctrl     *gomock.Controller
	recorder *MockRestorerMockRecorder
}

// MockRestorerMockRecorder is the mock recorder for MockRestorer.
type MockRestorerMockRecorder struct {
	mock *MockRestorer
}

// NewMockRestorer creates a new mock instance.
func NewMockRestorer(ctrl *gomock.Controller) *MockRestorer {
	mock := &MockRestorer{ctrl: ctrl}
	mock.recorder = &MockRestorerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRestorer) EXPECT() *MockRestorerMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockRestorer) Get(ctx context.Context, req *entities.GetFileRequest) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, req)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRestorerMockRecorder) Get(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRestorer)(nil).Get), ctx, req)
}