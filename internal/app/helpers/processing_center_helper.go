package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	opcPb "github.com/voonik/goConnect/api/go/oms/processing_center"
	opcService "github.com/voonik/goConnect/oms/processing_center"
	"github.com/voonik/ss2/internal/app/utils"
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

func GetOPCListForCurrentUser(ctx context.Context) []uint64 {
	opcList := []uint64{}

	userId := *utils.GetCurrentUserID(ctx)
	resp, err := getOpcClient().GetProcessingCenterListWithUserId(ctx, userId)
	if err != nil {
		log.Printf("GetOPCListForCurrentUser: Failed to fetch OPC list. Error: %v\n", err)
		return opcList
	}

	for _, opc := range resp.Data {
		opcList = append(opcList, opc.OpcId)
	}

	log.Printf("GetOPCListForCurrentUser: opc list = %v\n", opcList)
	return opcList
}

func IsOpcListValid(ctx context.Context, opcIds []uint64) error {
	if len(opcIds) == 0 {
		return nil
	}

	resp, err := getOpcClient().GetProcessingCenterListWithOpcIds(ctx, opcIds)
	if err != nil {
		log.Printf("IsOpcListValid: failed to fetch opc list. Error: %v\n", err)
		return errors.New("failed to fetch opc list")
	}

	opcMap := map[uint64]bool{}
	for _, opc := range resp.Data {
		opcMap[opc.OpcId] = true
	}

	for _, id := range opcIds {
		if _, found := opcMap[id]; !found {
			return fmt.Errorf("invalid opc id #(%v)", id)
		}
	}

	return nil
}
