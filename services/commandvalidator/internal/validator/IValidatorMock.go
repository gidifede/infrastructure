// Code generated by MockGen. DO NOT EDIT.
// Source: ./IValidator.go

// Package validator is a generated GoMock package.
package validator

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockIValidator is a mock of IValidator interface.
type MockIValidator struct {
	ctrl     *gomock.Controller
	recorder *MockIValidatorMockRecorder
}

// MockIValidatorMockRecorder is the mock recorder for MockIValidator.
type MockIValidatorMockRecorder struct {
	mock *MockIValidator
}

// NewMockIValidator creates a new mock instance.
func NewMockIValidator(ctrl *gomock.Controller) *MockIValidator {
	mock := &MockIValidator{ctrl: ctrl}
	mock.recorder = &MockIValidatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIValidator) EXPECT() *MockIValidatorMockRecorder {
	return m.recorder
}

// CommandMatchAPIPath mocks base method.
func (m *MockIValidator) CommandMatchAPIPath(command, apiPath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommandMatchAPIPath", command, apiPath)
	ret0, _ := ret[0].(error)
	return ret0
}

// CommandMatchAPIPath indicates an expected call of CommandMatchAPIPath.
func (mr *MockIValidatorMockRecorder) CommandMatchAPIPath(command, apiPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommandMatchAPIPath", reflect.TypeOf((*MockIValidator)(nil).CommandMatchAPIPath), command, apiPath)
}

// GetJSONSchema mocks base method.
func (m *MockIValidator) GetJSONSchema(ctx context.Context, filename string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJSONSchema", ctx, filename)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJSONSchema indicates an expected call of GetJSONSchema.
func (mr *MockIValidatorMockRecorder) GetJSONSchema(ctx, filename interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJSONSchema", reflect.TypeOf((*MockIValidator)(nil).GetJSONSchema), ctx, filename)
}

// GetMainSchemaPath mocks base method.
func (m *MockIValidator) GetMainSchemaPath(commandType, commandOrigin string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMainSchemaPath", commandType, commandOrigin)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMainSchemaPath indicates an expected call of GetMainSchemaPath.
func (mr *MockIValidatorMockRecorder) GetMainSchemaPath(commandType, commandOrigin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMainSchemaPath", reflect.TypeOf((*MockIValidator)(nil).GetMainSchemaPath), commandType, commandOrigin)
}

// GetReferenceSchemaPath mocks base method.
func (m *MockIValidator) GetReferenceSchemaPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReferenceSchemaPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetReferenceSchemaPath indicates an expected call of GetReferenceSchemaPath.
func (mr *MockIValidatorMockRecorder) GetReferenceSchemaPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReferenceSchemaPath", reflect.TypeOf((*MockIValidator)(nil).GetReferenceSchemaPath))
}

// ValidateCommand mocks base method.
func (m *MockIValidator) ValidateCommand(body, bodyJSONSchema, bodyJSONSchemaRef []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateCommand", body, bodyJSONSchema, bodyJSONSchemaRef)
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateCommand indicates an expected call of ValidateCommand.
func (mr *MockIValidatorMockRecorder) ValidateCommand(body, bodyJSONSchema, bodyJSONSchemaRef interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateCommand", reflect.TypeOf((*MockIValidator)(nil).ValidateCommand), body, bodyJSONSchema, bodyJSONSchemaRef)
}
