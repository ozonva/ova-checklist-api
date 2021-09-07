// Code generated by MockGen. DO NOT EDIT.
// Source: ../repo.go

// Package mock_repo is a generated GoMock package.
package mock_repo

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/ozonva/ova-checklist-api/internal/types"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// AddChecklists mocks base method.
func (m *MockRepo) AddChecklists(ctx context.Context, checklists []types.Checklist) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddChecklists", ctx, checklists)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddChecklists indicates an expected call of AddChecklists.
func (mr *MockRepoMockRecorder) AddChecklists(ctx, checklists interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddChecklists", reflect.TypeOf((*MockRepo)(nil).AddChecklists), ctx, checklists)
}

// DescribeChecklist mocks base method.
func (m *MockRepo) DescribeChecklist(ctx context.Context, userId uint64, checklistId string) (*types.Checklist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeChecklist", ctx, userId, checklistId)
	ret0, _ := ret[0].(*types.Checklist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeChecklist indicates an expected call of DescribeChecklist.
func (mr *MockRepoMockRecorder) DescribeChecklist(ctx, userId, checklistId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeChecklist", reflect.TypeOf((*MockRepo)(nil).DescribeChecklist), ctx, userId, checklistId)
}

// ListChecklists mocks base method.
func (m *MockRepo) ListChecklists(ctx context.Context, userId, limit, offset uint64) ([]types.Checklist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChecklists", ctx, userId, limit, offset)
	ret0, _ := ret[0].([]types.Checklist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListChecklists indicates an expected call of ListChecklists.
func (mr *MockRepoMockRecorder) ListChecklists(ctx, userId, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChecklists", reflect.TypeOf((*MockRepo)(nil).ListChecklists), ctx, userId, limit, offset)
}

// RemoveChecklist mocks base method.
func (m *MockRepo) RemoveChecklist(ctx context.Context, userId uint64, checklistId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveChecklist", ctx, userId, checklistId)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveChecklist indicates an expected call of RemoveChecklist.
func (mr *MockRepoMockRecorder) RemoveChecklist(ctx, userId, checklistId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveChecklist", reflect.TypeOf((*MockRepo)(nil).RemoveChecklist), ctx, userId, checklistId)
}

// UpdateChecklist mocks base method.
func (m *MockRepo) UpdateChecklist(ctx context.Context, checklist types.Checklist) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateChecklist", ctx, checklist)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateChecklist indicates an expected call of UpdateChecklist.
func (mr *MockRepoMockRecorder) UpdateChecklist(ctx, checklist interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateChecklist", reflect.TypeOf((*MockRepo)(nil).UpdateChecklist), ctx, checklist)
}