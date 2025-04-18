// Code generated by MockGen. DO NOT EDIT.
// Source: infrastructure.go
//
// Generated by this command:
//
//      mockgen -source=infrastructure.go
//

// Package mock_infrastructure is a generated GoMock package.
package mock_infrastructure

import (
        reflect "reflect"

        gomock "go.uber.org/mock/gomock"
)

// MockAlertProducing is a mock of AlertProducing interface.
type MockAlertProducing struct {
        ctrl     *gomock.Controller
        recorder *MockAlertProducingMockRecorder
        isgomock struct{}
}

// MockAlertProducingMockRecorder is the mock recorder for MockAlertProducing.
type MockAlertProducingMockRecorder struct {
        mock *MockAlertProducing
}

// NewMockAlertProducing creates a new mock instance.
func NewMockAlertProducing(ctrl *gomock.Controller) *MockAlertProducing {
        mock := &MockAlertProducing{ctrl: ctrl}
        mock.recorder = &MockAlertProducingMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAlertProducing) EXPECT() *MockAlertProducingMockRecorder {
        return m.recorder
}

// PublicMessage mocks base method.
func (m *MockAlertProducing) PublicMessage(body []byte) error {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "PublicMessage", body)
        ret0, _ := ret[0].(error)
        return ret0
}

// PublicMessage indicates an expected call of PublicMessage.
func (mr *MockAlertProducingMockRecorder) PublicMessage(body any) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublicMessage", reflect.TypeOf((*MockAlertProducing)(nil).PublicMessage), body)
}