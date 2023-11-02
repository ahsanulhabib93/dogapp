package services

import (
	"context"
	"errors"

	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	"github.com/voonik/goFramework/pkg/database"

	"github.com/voonik/ss2/internal/app/models"
)

type SellerService struct{}

func (ss *SellerService) GetByUserID(ctx context.Context, params *spb.GetByUserIDParams) (*spb.GetByUserIDResponse, error) {
	userId := params.GetUserId()
	if userId == 0 {
		return nil, errors.New("User ID is empty or zero")
	}
	seller := &spb.SellerObject{}
	query := database.DBAPM(ctx).Model(&models.Seller{}).Where("user_id = ?", userId)
	err := query.Scan(seller).Error
	if err != nil {
		return nil, errors.New("Seller not found")
	}
	response := spb.GetByUserIDResponse{
		Seller: seller,
	}

	return &response, nil
}
