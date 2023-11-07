package services

import (
	"context"
	"fmt"
	"log"
	"strings"

	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	"github.com/voonik/goFramework/pkg/database"

	"github.com/voonik/ss2/internal/app/models"
)

type SellerService struct{}

func (ss *SellerService) GetByUserID(ctx context.Context, params *spb.GetByUserIDParams) (*spb.GetByUserIDResponse, error) {
	userId := params.GetUserId()
	seller := &spb.SellerObject{}
	query := database.DBAPM(ctx).Model(&models.Seller{}).Where("user_id = ?", userId)
	err := query.Scan(seller).Error
	if err != nil {
		log.Println("Error in seller service:", err.Error())
		return &spb.GetByUserIDResponse{}, nil
	}
	response := spb.GetByUserIDResponse{
		Seller: seller,
	}

	return &response, nil
}

func (ss *SellerService) GetSellerByCondition(ctx context.Context, params *spb.GetSellerByConditionParams) (*spb.GetSellersResponse, error) {
	return nil, nil
}

func (ss *SellerService) GetSellersRelatedToOrder(ctx context.Context, params *spb.GetSellersRelatedToOrderParams) (*spb.GetSellersResponse, error) {
	return nil, nil
}

func (ss *SellerService) SmallReport(ctx context.Context, params *spb.SmallReportParams) (*spb.GetSellersResponse, error) {
	return nil, nil
}

func (ss *SellerService) ValidateField(ctx context.Context, params *spb.ValidateFieldParams) (*spb.StatusResponse, error) {
	response := spb.StatusResponse{}
	response.Status = false
	seller := &spb.SellerObject{}
	data := params.GetData()
	if data == nil {
		return &response, nil
	}
	conditionString := make([]string, 0)
	for key, value := range data {
		conditionString = append(conditionString, fmt.Sprintf("%s = '%s'", key, value))
	}
	queryCondition := strings.Join(conditionString, " AND ")
	query := database.DBAPM(ctx).Model(&models.Seller{}).Where(queryCondition)
	if err := query.Scan(seller).Error; err != nil {
		response.Status = true
		return &response, nil
	}
	return &response, nil
}

func (ss *SellerService) SellerPhoneRelation(ctx context.Context, params *spb.SellerPhoneRelationParams) (*spb.GetSellersResponse, error) {
	return nil, nil
}

func (ss *SellerService) ApproveProducts(ctx context.Context, params *spb.ApproveProductsParams) (*spb.BasicApiResponse, error) {
	return nil, nil
}

func (ss *SellerService) ConfirmEmailFromAdminPanel(ctx context.Context, params *spb.GetByUserIDParams) (*spb.BasicApiResponse, error) {
	return nil, nil
}

func (ss *SellerService) Update(ctx context.Context, params *spb.UpdateParams) (*spb.BasicApiResponse, error) {
	return nil, nil
}

func (ss *SellerService) SendActivationMail(ctx context.Context, params *spb.SendActivationMailParams) (*spb.BasicApiResponse, error) {
	return nil, nil
}
