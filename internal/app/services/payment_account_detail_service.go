package services

import (
	"context"
	"fmt"
	"log"

	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

// PaymentAccountDetailService ...
type PaymentAccountDetailService struct{}

// List ...
func (ps *PaymentAccountDetailService) List(ctx context.Context, params *paymentpb.ListParams) (*paymentpb.ListResponse, error) {
	resp := paymentpb.ListResponse{}
	database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).Where("supplier_id = ?", params.GetSupplierId()).Scan(&resp.Data)
	return &resp, nil
}

// Add ...
func (ps *PaymentAccountDetailService) Add(ctx context.Context, params *paymentpb.PaymentAccountDetailParam) (*paymentpb.BasicApiResponse, error) {
	resp := paymentpb.BasicApiResponse{Success: false}

	supplier := &models.Supplier{}
	result := database.DBAPM(ctx).Model(&models.Supplier{}).First(supplier, params.GetSupplierId())
	if result.RecordNotFound() {
		resp.Message = "Supplier Not Found"
	} else {
		paymentAccountDetail := models.PaymentAccountDetail{
			Supplier:      *supplier,
			AccountType:   utils.AccountType(params.GetAccountType()),
			AccountName:   params.GetAccountName(),
			AccountNumber: params.GetAccountNumber(),
			BankName:      params.GetBankName(),
			BranchName:    params.GetBranchName(),
			RoutingNumber: params.GetRoutingNumber(),
			IsDefault:     params.GetIsDefault(),
		}
		err := database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).Create(&paymentAccountDetail)

		if err != nil && err.Error != nil {
			errorMsg := fmt.Sprintf("Error while creating PaymentAccountDetail: %s", err.Error)
			log.Println(errorMsg)
			resp.Message = errorMsg
		} else {
			resp.Message = "PaymentAccountDetail Added Successfully"
			resp.Success = true
		}
	}
	return &resp, nil
}
