package services

import (
	"context"
	"errors"

	sampb "github.com/voonik/goConnect/api/go/ss2/seller_account_manager"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type SellerAccountManagerService struct{}

func (sams *SellerAccountManagerService) List(ctx context.Context, params *sampb.ListParams) (*sampb.ListResponse, error) {
	resp := &sampb.ListResponse{}

	if params.GetSellerId() == utils.Zero && params.GetId() == utils.Zero {
		return resp, nil
	}

	var samList []models.SellerAccountManager
	var accountManagers []*sampb.AccountManagerObject
	query := database.DBAPM(ctx).Model(&models.SellerAccountManager{})
	if params.GetSellerId() != utils.Zero {
		query = query.Where(`seller_id =?`, params.SellerId)
	}
	if params.GetId() != utils.Zero {
		query = query.Where(`id =?`, params.GetId())
	}
	query = query.Order("role, priority").Scan(&samList)
	for _, sam := range samList {
		accountManagers = append(accountManagers, &sampb.AccountManagerObject{
			Id:       sam.ID,
			Email:    sam.Email,
			Phone:    uint64(sam.Phone),
			Name:     sam.Name,
			Priority: uint64(sam.Priority),
			Role:     sam.Role,
		})
	}

	resp.Status = "success"
	resp.AccountManager = accountManagers
	return resp, nil
}

func (sams *SellerAccountManagerService) Update(ctx context.Context, params *sampb.AccountManagerObject) (*sampb.BasicApiResponse, error) {
	response := &sampb.BasicApiResponse{Success: true, Message: "update successfull"}

	sam, err := getSamFromParams(ctx, params)
	if err != nil {
		return &sampb.BasicApiResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	updateParams := &models.SellerAccountManager{
		Phone: int64(params.GetPhone()),
		Role:  params.GetRole(),
		Name:  params.GetName(),
		Email: params.GetEmail(),
	}
	err = database.DBAPM(ctx).Model(sam).Updates(updateParams).Error
	if err != nil {

	}
	return response, nil
}

func (sams *SellerAccountManagerService) Delete(ctx context.Context, params *sampb.AccountManagerObject) (*sampb.BasicApiResponse, error) {
	response := &sampb.BasicApiResponse{Success: true, Message: "deletion successfull"}

	sam, err := getSamFromParams(ctx, params)
	if err != nil {
		return &sampb.BasicApiResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	err = database.DBAPM(ctx).Delete(sam).Error
	if err != nil {
		return &sampb.BasicApiResponse{
			Success: false,
			Message: "unable to delete seller account manager" + err.Error(),
		}, nil
	}
	return response, nil
}

func getSamFromParams(ctx context.Context, params *sampb.AccountManagerObject) (*models.SellerAccountManager, error) {
	sam := &models.SellerAccountManager{}
	if params.GetId() == utils.Zero {
		return nil, errors.New("id cannot be empty")
	}
	err := database.DBAPM(ctx).Model(&models.SellerAccountManager{}).Where("id = ? ", params.GetId()).Find(sam).Error
	if err != nil {
		return nil, err
	}
	return sam, nil
}
