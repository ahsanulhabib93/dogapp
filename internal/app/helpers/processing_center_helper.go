package helpers

import (
	"context"
	"reflect"

	opcPb "github.com/voonik/goConnect/api/go/oms/processing_center"
	opcService "github.com/voonik/goConnect/oms/processing_center"
)

type OpcHelper struct{}

var opcClient OpcClientInterface

func InjectMockOpcClientInstance(mockObj OpcClientInterface) {
	opcClient = mockObj
}

type OpcClientInterface interface {
	GetProcessingCenterListWithUserId(ctx context.Context, userId uint64) (*opcPb.ProcessingCenterListResponse, error)
	GetProcessingCenterListWithOpcIds(ctx context.Context, opcIds []uint64) (*opcPb.ProcessingCenterListResponse, error)
}

func getOpcClient() OpcClientInterface {
	if opcClient == nil || reflect.ValueOf(opcClient).IsNil() {
		return new(OpcHelper)
	}
	return opcClient
}

func (s *OpcHelper) GetProcessingCenterListWithUserId(ctx context.Context, userId uint64) (*opcPb.ProcessingCenterListResponse, error) {
	return opcService.ProcessingCenter().ProcessingCenterList(ctx, &opcPb.OpcListParams{UserId: userId})
}

func (s *OpcHelper) GetProcessingCenterListWithOpcIds(ctx context.Context, opcIds []uint64) (*opcPb.ProcessingCenterListResponse, error) {
	return opcService.ProcessingCenter().ProcessingCenterList(ctx, &opcPb.OpcListParams{OpcId: opcIds})
}
