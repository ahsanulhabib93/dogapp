package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"github.com/shopuptech/go-libs/logger"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	"github.com/voonik/goFramework/pkg/database"

	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type SellerService struct{}

func (ss *SellerService) GetByUserID(ctx context.Context, params *spb.GetByUserIDParams) (*spb.GetByUserIDResponse, error) {
	userId := params.GetUserId()
	seller := helpers.GetSellerByUserId(ctx, userId)
	if seller.ID == utils.Zero {
		return &spb.GetByUserIDResponse{}, nil
	}
	// convert model data into proto object
	sellerData := &spb.SellerObject{}
	copier.Copy(&sellerData, &seller) //nolint:errcheck
	sellerData.SellerConfig = helpers.GetDefaultSellerConfig()
	sellerData.ReturnExchangePolicy = helpers.DefaultsellerReturnExchangePolicy()
	resp := &spb.GetByUserIDResponse{Seller: sellerData}
	return resp, nil
}

func (SellerService) GetSellerByCondition(ctx context.Context, params *spb.GetSellerByConditionParams) (*spb.GetSellersResponse, error) {
	response := spb.GetSellersResponse{Status: utils.Failure}
	sellers := []*spb.SellerObject{}
	fields := params.GetFields()
	condition := params.GetCondition()
	if condition == nil {
		logger.FromContext(ctx).Info("No condition specified")
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
		response.Message = fmt.Sprint("Error in seller service GetSellerByCondition API", err.Error())
		logger.FromContext(ctx).Error(response.Message)
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
	userIds := params.GetSellerIds()
	response := spb.GetSellersResponse{
		Status: utils.Failure,
	}
	if len(userIds) == 0 {
		response.Message = "no valid param"
		return &response, nil
	}
	sellers := []*models.Seller{}
	sellerData := []*spb.SellerObject{}
	query := database.DBAPM(ctx).Model(&models.Seller{}).Where("user_id in (?)", userIds)
	err := query.Scan(&sellers).Error
	if err != nil {
		response.Message = fmt.Sprint("Error in seller service GetSellersRelatedToOrder API", err.Error())
		logger.FromContext(ctx).Error(response.Message)
		return &response, nil
	}
	copier.Copy(&sellerData, &sellers) //nolint:errcheck
	if len(sellerData) == 0 {
		response.Status = utils.Success
		response.Message = "seller not found"
		return &response, nil
	}
	for _, seller := range sellerData {
		seller.SellerConfig = helpers.GetDefaultSellerConfig()
		seller.ReturnExchangePolicy = helpers.DefaultsellerReturnExchangePolicy()
	}
	response.Seller = sellerData
	response.Status = utils.Success
	response.Message = "fetched seller details successfully"
	return &response, nil
}

func (ss *SellerService) SmallReport(ctx context.Context, params *spb.SmallReportParams) (*spb.GetSellersResponse, error) {
	return &spb.GetSellersResponse{}, nil
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
	logger.Log().Info("Approve Products API Params: %+v", params)
	resp := &spb.BasicApiResponse{Status: utils.Failure}
	if params.GetId() == utils.Zero || len(params.GetIds()) == utils.Zero {
		resp.Message = "UserID & Product IDs Should Not Empty to Approve Products"
	} else {
		resp = helpers.PerformApproveProductFunc(ctx, params)
	}
	return resp, nil
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
	query := database.DBAPM(ctx).Model(&models.Seller{}).Where("user_id = ?", id)
	err := query.First(&seller).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			response.Status = "failure"
			response.Message = "seller not found"
			return &response, nil
		}
		logger.FromContext(ctx).Info("Error in seller service Update API", err.Error())
		response.Message = fmt.Sprint("unable to update seller: ", err.Error())
		return &response, nil
	}
	paramJSON, err := json.Marshal(sellerParam)
	if err != nil {
		logger.FromContext(ctx).Info("Error marshaling SellerObject to JSON", err.Error())
		response.Message = fmt.Sprintf("unable to update seller: %s", err.Error())
		return &response, nil
	}
	var sellerUpdates map[string]interface{}
	if err := json.Unmarshal(paramJSON, &sellerUpdates); err != nil {
		logger.FromContext(ctx).Info("Error unmarshaling JSON to map", err.Error())
		response.Message = fmt.Sprintf("unable to update seller: %s", err.Error())
		return &response, nil
	}
	err = query.Updates(sellerUpdates).Error
	if err != nil {
		logger.FromContext(ctx).Info("Error in seller service Update API", err.Error())
		response.Message = fmt.Sprintf("unable to update seller: %s", err.Error())
		return &response, nil
	}
	response.Status = utils.Success
	response.Message = "seller details updated successfully"
	return &response, nil
}

func (ss *SellerService) SendActivationMail(ctx context.Context, params *spb.SendActivationMailParams) (*spb.BasicApiResponse, error) {
	logger.Log().Infof("Send Activation Mail API Params: %+v", params)
	resp := &spb.BasicApiResponse{Status: utils.Failure}
	if params.GetId() == utils.EmptyString {
		resp.Message = "Seller UserIds Should be Present"
	} else {
		msg, sellerIDs := helpers.GetArrayIdsFromString(params.GetId())
		if msg == utils.EmptyString && len(sellerIDs) > utils.Zero { // TODO: validate params.GetAction()
			resp = helpers.PerformSendActivationMail(ctx, sellerIDs, params)
		} else {
			resp.Message = msg
		}
	}
	return resp, nil
}

func (ss *SellerService) Create(ctx context.Context, params *spb.CreateParams) (*spb.CreateResponse, error) {
	resp := &spb.CreateResponse{Status: false, Message: "Failed to register the seller. Please try again."}

	if params.Seller == nil {
		resp.Message = "Missing Seller Params"
		return resp, nil
	}

	existingSeller := helpers.GetSellerByUserId(ctx, params.Seller.UserId)
	if existingSeller != nil {
		resp.Status = true
		resp.Message = "Seller already registered."
		resp.UserId = existingSeller.UserID
		return resp, nil
	}

	seller := helpers.FormatAndAssignData(params)
	err := database.DBAPM(ctx).Model(&models.Seller{}).Create(seller).Error

	if err == nil {
		resp.Message = "Seller registered successfully."
		resp.Status = true
		resp.UserId = seller.UserID
	}
	return resp, nil
}
