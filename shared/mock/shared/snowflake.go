// Code generated by MockGen. DO NOT EDIT.
// Source: food-story/shared/snowflakeid (interfaces: SnowflakeInterface)
//
// Generated by this command:
//
//	mockgen -package mockshared -destination shared/mock/shared/snowflake.go food-story/shared/snowflakeid SnowflakeInterface
//

// Package mockshared is a generated GoMock package.
package mockshared

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockSnowflakeInterface is a mock of SnowflakeInterface interface.
type MockSnowflakeInterface struct {
	ctrl     *gomock.Controller
	recorder *MockSnowflakeInterfaceMockRecorder
	isgomock struct{}
}

// MockSnowflakeInterfaceMockRecorder is the mock recorder for MockSnowflakeInterface.
type MockSnowflakeInterfaceMockRecorder struct {
	mock *MockSnowflakeInterface
}

// NewMockSnowflakeInterface creates a new mock instance.
func NewMockSnowflakeInterface(ctrl *gomock.Controller) *MockSnowflakeInterface {
	mock := &MockSnowflakeInterface{ctrl: ctrl}
	mock.recorder = &MockSnowflakeInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSnowflakeInterface) EXPECT() *MockSnowflakeInterfaceMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockSnowflakeInterface) Generate() int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate")
	ret0, _ := ret[0].(int64)
	return ret0
}

// Generate indicates an expected call of Generate.
func (mr *MockSnowflakeInterfaceMockRecorder) Generate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockSnowflakeInterface)(nil).Generate))
}
