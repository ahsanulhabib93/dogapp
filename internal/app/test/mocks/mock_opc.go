package mocks

import (
	"context"

	opcPb "github.com/voonik/goConnect/api/go/oms/processing_center"
	"github.com/voonik/ss2/internal/app/helpers"
)

type MockOpcHelper struct{}

func SetOpcMock() {
	mock := &MockOpcHelper{}
	helpers.InjectMockOpcClientInstance(mock)
}

func UnsetOpcMock() {
	helpers.InjectMockOpcClientInstance(nil)
}

func (s *MockOpcHelper) ProcessingCenterList(ctx context.Context) (*opcPb.ProcessingCenterListResponse, error) {
	return &opcPb.ProcessingCenterListResponse{
		Data: []*opcPb.OpcDetail{
			{OpcId: 201},
			{OpcId: 202},
		},
	}, nil
}
