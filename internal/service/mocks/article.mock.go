// Code generated by MockGen. DO NOT EDIT.
// Source: webook/internal/service/article.go

// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/jw803/webook/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockArticleService is a mock of ArticleService interface.
type MockArticleService struct {
	ctrl     *gomock.Controller
	recorder *MockArticleServiceMockRecorder
}

// MockArticleServiceMockRecorder is the mock recorder for MockArticleService.
type MockArticleServiceMockRecorder struct {
	mock *MockArticleService
}

// NewMockArticleService creates a new mock instance.
func NewMockArticleService(ctrl *gomock.Controller) *MockArticleService {
	mock := &MockArticleService{ctrl: ctrl}
	mock.recorder = &MockArticleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockArticleService) EXPECT() *MockArticleServiceMockRecorder {
	return m.recorder
}

// Publish mocks base method.
func (m *MockArticleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, art)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Publish indicates an expected call of Publish.
func (mr *MockArticleServiceMockRecorder) Publish(ctx, art interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockArticleService)(nil).Publish), ctx, art)
}

// PublishV1 mocks base method.
func (m *MockArticleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishV1", ctx, art)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PublishV1 indicates an expected call of PublishV1.
func (mr *MockArticleServiceMockRecorder) PublishV1(ctx, art interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishV1", reflect.TypeOf((*MockArticleService)(nil).PublishV1), ctx, art)
}

// Save mocks base method.
func (m *MockArticleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, art)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockArticleServiceMockRecorder) Save(ctx, art interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockArticleService)(nil).Save), ctx, art)
}
