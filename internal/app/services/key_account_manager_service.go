package services

import (
	"context"
	"fmt"
	"log"

	kampb "github.com/voonik/goConnect/api/go/ss2/key_account_manager"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
)

// KeyAccountManagerService ...
type KeyAccountManagerService struct{}

// List ...
func (kams *KeyAccountManagerService) List(ctx context.Context, params *kampb.ListParams) (*kampb.ListResponse, error) {
	resp := kampb.ListResponse{}
	database.DBAPM(ctx).Model(&models.KeyAccountManager{}).Where("supplier_id = ?", params.GetSupplierId()).Scan(&resp.Data)
	return &resp, nil
}

// Add ...
func (kams *KeyAccountManagerService) Add(ctx context.Context, params *kampb.KeyAccountManagerParam) (*kampb.BasicApiResponse, error) {
	resp := kampb.BasicApiResponse{Success: false}

	supplier := &models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).First(supplier, params.GetSupplierId())
	if result.RecordNotFound() {
		resp.Message = "Supplier Not Found"
	} else {
		keyAccountManager := models.KeyAccountManager{
			Supplier: *supplier,
			Name:     params.GetName(),
			Email:    params.GetEmail(),
			Phone:    params.GetPhone(),
		}
		err := database.DBAPM(ctx).Model(&models.KeyAccountManager{}).Create(&keyAccountManager)

		if err != nil && err.Error != nil {
			errorMsg := fmt.Sprintf("Error while creating KeyAccountManager: %s", err.Error)
			log.Println(errorMsg)
			resp.Message = errorMsg
		} else {
			resp.Message = "KeyAccountManager Added Successfully"
			resp.Success = true
		}
	}
	return &resp, nil
}
