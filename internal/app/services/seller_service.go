package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/shopuptech/go-libs/logger"
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
	if condition == nil {
		logger.FromContext(ctx).Info("No condition specified")
		response.Status = utils.Failure
		response.Message = "no condition specified"
		return &response, nil
	}
	conditionString := make([]string, 0)
	for key, value := range condition {
		conditionString = append(conditionString, fmt.Sprintf("%s = '%s'", key, value))
	}
	queryCondition := strings.Join(conditionString, " AND ")
	query := database.DBAPM(ctx).Model(&models.Seller{}).Where(queryCondition)
	if fields != nil {
		query = query.Select(fields)
	}
	if err := query.Scan(&sellers).Error; err != nil {
		logger.FromContext(ctx).Info("Error in seller service GetSellerByCondition API", err.Error())
		response.Status = utils.Failure
		response.Message = err.Error()
		return &response, nil
	}
	if len(sellers) == 0 {
		logger.FromContext(ctx).Info("Seller not found")
		response.Status = utils.Success
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
	phone := params.GetPhone()
	response := spb.GetSellersResponse{}
	response.Status = utils.Failure
	if len(phone) == 0 {
		response.Message = "no valid param"
		return &response, nil
	}
	sellers := []*spb.SellerObject{}
	query := database.DBAPM(ctx).Model(&models.Seller{}).Where("primary_phone in (?)", phone)
	err := query.Scan(&sellers).Error
	if err != nil {
		response.Message = "error in SellerPhoneRelation"
		return &response, nil
	}
	if len(sellers) == 0 {
		response.Status = "success"
		response.Message = "seller not found"
		return &response, nil
	}
	response.Seller = sellers
	response.Status = "success"
	response.Message = "fetched seller details successfully"
	return &response, nil

}

func (ss *SellerService) ApproveProducts(ctx context.Context, params *spb.ApproveProductsParams) (*spb.BasicApiResponse, error) {
	return nil, nil
}

func (ss *SellerService) ConfirmEmailFromAdminPanel(ctx context.Context, params *spb.GetByUserIDParams) (*spb.BasicApiResponse, error) {
	response := spb.BasicApiResponse{Status: utils.Failure}
	userId := params.GetUserId()
	if userId == 0 {
		response.Message = "param not specified"
		return &response, nil
	}
	seller := models.Seller{}
	query := database.DBAPM(ctx).Model(&models.Seller{}).Where("user_id = ?", userId)
	err := query.Scan(&seller).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			response.Message = "seller not found"
			return &response, nil
		}
		logger.FromContext(ctx).Error("Error in seller service ConfirmEmailFromAdminPanel API", err.Error())
		response.Message = fmt.Sprint("not able to confirm email: ", err.Error())
		return &response, nil
	}
	if seller.EmailConfirmed {
		response.Status = "success"
		response.Message = "email already confirmed"
		return &response, nil
	}
	err = query.Update("email_confirmed", true).Error
	if err != nil {
		logger.FromContext(ctx).Error("Error in seller service ConfirmEmailFromAdminPanel API", err.Error())
		response.Message = fmt.Sprint("not able to confirm email: ", err.Error())
		return &response, nil
	}
	response.Status = utils.Success
	response.Message = "email confirmed successfully"
	return &response, nil
}

func (ss *SellerService) Update(ctx context.Context, params *spb.UpdateParams) (*spb.BasicApiResponse, error) {
	response := spb.BasicApiResponse{Status: utils.Failure}
	id := params.GetId()
	sellerParam := params.GetSeller()
	if id == 0 || sellerParam == nil {
		response.Message = "param not specified"
		return &response, nil
	}
	seller := models.Seller{}
	err := database.DBAPM(ctx).Model(&models.Seller{}).Where("user_id = ?", id).First(&seller).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			response.Status = "failure"
			response.Message = "seller not found"
			return &response, nil
		}
		logger.FromContext(ctx).Info("Error in seller service Update API", err.Error())
		response.Message = fmt.Sprint("unable to update seller1: ", err.Error())
		return &response, nil
	}
	paramJson, err := json.Marshal(sellerParam)
	if err != nil {
		logger.FromContext(ctx).Info("Error marshaling SellerObject to JSON", err.Error())
		response.Message = fmt.Sprintf("unable to update seller: %s", err.Error())
		return &response, nil
	}
	sellerParamValue := models.Seller{}
	err = json.Unmarshal(paramJson, &sellerParamValue)
	if err != nil {
		logger.FromContext(ctx).Info("Error unmarshaling JSON to Seller struct", err.Error())
		response.Message = fmt.Sprintf("unable to update seller: %s", err.Error())
		return &response, nil
	}
	seller = sellerParamValue
	seller.UserID = id
	err = database.DBAPM(ctx).Save(&seller).Error
	if err != nil {
		logger.FromContext(ctx).Info("Error in seller service Update API", err.Error())
		response.Message = fmt.Sprint("unable to update seller: ", err.Error())
		return &response, nil
	}
	response.Status = utils.Success
	response.Message = "seller details updated successfully"
	return &response, nil
}

func (ss *SellerService) SendActivationMail(ctx context.Context, params *spb.SendActivationMailParams) (*spb.BasicApiResponse, error) {
	return nil, nil
}
