// Code generated by MockGen. DO NOT EDIT.
// Source: ./ad.go

// Package handlermock is a generated GoMock package.
package handlermock

import (
	context "context"
	models "homework10/internal/domain/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAdService is a mock of AdService interface.
type MockAdService struct {
	ctrl     *gomock.Controller
	recorder *MockAdServiceMockRecorder
}

// MockAdServiceMockRecorder is the mock recorder for MockAdService.
type MockAdServiceMockRecorder struct {
	mock *MockAdService
}

// NewMockAdService creates a new mock instance.
func NewMockAdService(ctrl *gomock.Controller) *MockAdService {
	mock := &MockAdService{ctrl: ctrl}
	mock.recorder = &MockAdServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdService) EXPECT() *MockAdServiceMockRecorder {
	return m.recorder
}

// ChangeAdStatus mocks base method.
func (m *MockAdService) ChangeAdStatus(ctx context.Context, adID, userID int64, published bool) (*models.Ad, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeAdStatus", ctx, adID, userID, published)
	ret0, _ := ret[0].(*models.Ad)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangeAdStatus indicates an expected call of ChangeAdStatus.
func (mr *MockAdServiceMockRecorder) ChangeAdStatus(ctx, adID, userID, published interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeAdStatus", reflect.TypeOf((*MockAdService)(nil).ChangeAdStatus), ctx, adID, userID, published)
}

// CreateAd mocks base method.
func (m *MockAdService) CreateAd(ctx context.Context, title, text string, authorID int64) (*models.Ad, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAd", ctx, title, text, authorID)
	ret0, _ := ret[0].(*models.Ad)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAd indicates an expected call of CreateAd.
func (mr *MockAdServiceMockRecorder) CreateAd(ctx, title, text, authorID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAd", reflect.TypeOf((*MockAdService)(nil).CreateAd), ctx, title, text, authorID)
}

// DeleteAd mocks base method.
func (m *MockAdService) DeleteAd(ctx context.Context, adID, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAd", ctx, adID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAd indicates an expected call of DeleteAd.
func (mr *MockAdServiceMockRecorder) DeleteAd(ctx, adID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAd", reflect.TypeOf((*MockAdService)(nil).DeleteAd), ctx, adID, userID)
}

// GetAdByID mocks base method.
func (m *MockAdService) GetAdByID(ctx context.Context, adID int64) (*models.Ad, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdByID", ctx, adID)
	ret0, _ := ret[0].(*models.Ad)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdByID indicates an expected call of GetAdByID.
func (mr *MockAdServiceMockRecorder) GetAdByID(ctx, adID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdByID", reflect.TypeOf((*MockAdService)(nil).GetAdByID), ctx, adID)
}

// GetAdsByTitle mocks base method.
func (m *MockAdService) GetAdsByTitle(ctx context.Context, text string) ([]*models.Ad, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdsByTitle", ctx, text)
	ret0, _ := ret[0].([]*models.Ad)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdsByTitle indicates an expected call of GetAdsByTitle.
func (mr *MockAdServiceMockRecorder) GetAdsByTitle(ctx, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdsByTitle", reflect.TypeOf((*MockAdService)(nil).GetAdsByTitle), ctx, text)
}

// ListAds mocks base method.
func (m *MockAdService) ListAds(ctx context.Context, published, userIDRaw, dateCreationRaw string) ([]*models.Ad, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAds", ctx, published, userIDRaw, dateCreationRaw)
	ret0, _ := ret[0].([]*models.Ad)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAds indicates an expected call of ListAds.
func (mr *MockAdServiceMockRecorder) ListAds(ctx, published, userIDRaw, dateCreationRaw interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAds", reflect.TypeOf((*MockAdService)(nil).ListAds), ctx, published, userIDRaw, dateCreationRaw)
}

// UpdateAd mocks base method.
func (m *MockAdService) UpdateAd(ctx context.Context, adID, userID int64, title, text string) (*models.Ad, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAd", ctx, adID, userID, title, text)
	ret0, _ := ret[0].(*models.Ad)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAd indicates an expected call of UpdateAd.
func (mr *MockAdServiceMockRecorder) UpdateAd(ctx, adID, userID, title, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAd", reflect.TypeOf((*MockAdService)(nil).UpdateAd), ctx, adID, userID, title, text)
}
