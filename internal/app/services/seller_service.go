package services

import (
	"context"
	"log"

	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	"github.com/voonik/goFramework/pkg/database"

	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
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

func (SellerService) GetSellerByCondition(ctx context.Context, params *spb.GetSellerByConditionParams) (*spb.GetSellersResponse, error) {
	response := spb.GetSellersResponse{}
	sellers := []*spb.SellerObject{}
	fields := params.GetFields()
	condition := params.GetCondition()
	if condition == "" {
		log.Println("No condition specified")
		response.Status = utils.Failure
		response.Message = "no condition specified"
		return &response, nil
	}
	query := database.DBAPM(ctx).Model(&models.Seller{}).Where(condition)
	if fields != "" {
		query = query.Select(fields)
	}
	query.Scan(&sellers)
	if len(sellers) == 0 {
		log.Println("Seller not found")
		response.Status = utils.Failure
		response.Message = "seller not found"
		return &response, nil
	}
	response.Seller = sellers
	response.Status = utils.Success
	response.Message = "fetched seller details successfully"
	return &response, nil
}

func (ss *SellerService) GetSellersRelatedToOrder(ctx context.Context, params *spb.GetSellersRelatedToOrderParams) (*spb.GetSellersResponse, error) {
	return nil, nil
}

func (ss *SellerService) SmallReport(ctx context.Context, params *spb.SmallReportParams) (*spb.GetSellersResponse, error) {
	return nil, nil
}

func (ss *SellerService) ValidateField(ctx context.Context, params *spb.ValidateFieldParams) (*spb.StatusResponse, error) {
	return nil, nil
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
