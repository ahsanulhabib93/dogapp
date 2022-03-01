package mocks

import (
	"context"

	opcPb "github.com/voonik/goConnect/api/go/oms/processing_center"
	"github.com/voonik/ss2/internal/app/helpers"

	mock "github.com/stretchr/testify/mock"
)

type MockOpcHelper struct {
	mock.Mock
}

func SetOpcMock() *MockOpcHelper {
	mock := &MockOpcHelper{}
	helpers.InjectMockOpcClientInstance(mock)

	return mock
}

func UnsetOpcMock() {
	helpers.InjectMockOpcClientInstance(nil)
}

func (_m *MockOpcHelper) ProcessingCenterList(ctx context.Context, userId uint64) (*opcPb.ProcessingCenterListResponse, error) {
	args := _m.Called(ctx, userId)
	return args.Get(0).(*opcPb.ProcessingCenterListResponse), args.Error(1)
}
