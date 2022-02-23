package helpers

import (
	"context"

	opcPb "github.com/voonik/goConnect/api/go/oms/processing_center"
	opcService "github.com/voonik/goConnect/oms/processing_center"
)

type OpcHelper struct{}

var opcClient OpcClientInterface

func InjectMockOpcClientInstance(mockObj OpcClientInterface) {
	opcClient = mockObj
}

type OpcClientInterface interface {
	processingCenterList(ctx context.Context) (*opcPb.ProcessingCenterListResponse, error)
}

func getOpcClient() OpcClientInterface {
	if opcClient == nil {
		return new(OpcHelper)
	}
	return opcClient
}

func (s *OpcHelper) processingCenterList(ctx context.Context) (*opcPb.ProcessingCenterListResponse, error) {
	return opcService.ProcessingCenter().ProcessingCenterList(ctx, &opcPb.EmptyParams{})
}
