// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/lechat/node-upgrader/internal/interfaces (interfaces: EKSAPI)
//
// Generated by this command:
//
//	mockgen -destination=internal/mocks/eksapi_mock.go -package=mocks github.com/lechat/node-upgrader/internal/interfaces EKSAPI
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	eks "github.com/aws/aws-sdk-go-v2/service/eks"
	gomock "go.uber.org/mock/gomock"
)

// MockEKSAPI is a mock of EKSAPI interface.
type MockEKSAPI struct {
	ctrl     *gomock.Controller
	recorder *MockEKSAPIMockRecorder
}

// MockEKSAPIMockRecorder is the mock recorder for MockEKSAPI.
type MockEKSAPIMockRecorder struct {
	mock *MockEKSAPI
}

// NewMockEKSAPI creates a new mock instance.
func NewMockEKSAPI(ctrl *gomock.Controller) *MockEKSAPI {
	mock := &MockEKSAPI{ctrl: ctrl}
	mock.recorder = &MockEKSAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEKSAPI) EXPECT() *MockEKSAPIMockRecorder {
	return m.recorder
}

// ListClusters mocks base method.
func (m *MockEKSAPI) ListClusters(arg0 context.Context, arg1 *eks.ListClustersInput, arg2 ...func(*eks.Options)) (*eks.ListClustersOutput, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListClusters", varargs...)
	ret0, _ := ret[0].(*eks.ListClustersOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListClusters indicates an expected call of ListClusters.
func (mr *MockEKSAPIMockRecorder) ListClusters(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClusters", reflect.TypeOf((*MockEKSAPI)(nil).ListClusters), varargs...)
}