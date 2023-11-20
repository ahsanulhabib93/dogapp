package services

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/shopuptech/go-libs/logger"
	vapb "github.com/voonik/goConnect/api/go/ss2/vendor_address"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/utils"
)

type VendorAddressService struct{}

func (va *VendorAddressService) GetData(ctx context.Context, params *vapb.GetDataParams) (*vapb.GetDataResponse, error) {
	return nil, nil
}

// def verify_address
// 	params.permit!
// 	current_user_id = current_agent_id
// 	vendor_address = VendorAddress.unscoped.existing_address.find_by_uuid(params[:id])
// 	vendor_address.verification_status = 'VERIFIED'
// 	@seller_id = vendor_address.seller.user_id
// 	SellerActivityLog.create({"seller_id"=>@seller_id,"action"=>"verify_address","user_id"=>current_user_id,"notes"=>{"status"=>'verified'}.to_json})
// 	json_response = {:status => "success", :message=> "Address status was changed to verified successfully"}
// 	unless vendor_address.save!
// 		json_response[:status] = "failure"
// 		json_response[:message] = vendor_address.errors.full_messages.to_sentence
// 	else
// 		from_email = "#{APP_CONFIG['brand'].camelize}<#{APP_CONFIG['reply_to_id']}>"
// 		to_email = vendor_address.seller.primary_email
// 		subject = "Address Verified"
// 		content = "Dear #{vendor_address.seller.brand_name} ,<br>Your address request for #{vendor_address.address1} has been accepted and it has been verified .<br><br>Thanks,<br>Team #{APP_CONFIG['brand'].camelize} "
// 		VMailer.send_html_mail(from_email, to_email, subject, content).deliver
// 	end
// 	json_response
// end

func (vas *VendorAddressService) VerifyAddress(ctx context.Context, params *vapb.VerifyAddressParams) (*vapb.BasicApiResponse, error) {
	response := vapb.BasicApiResponse{Status: utils.Failure}
	uuid := params.GetId()
	if uuid == "" {
		response.Message = "param not specified"
		return &response, nil
	}
	vendorAddress := models.VendorAddress{}
	query := database.DBAPM(ctx).Model(&models.VendorAddress{}).Where("uuid = ?", uuid)
	err := query.Scan(&vendorAddress).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			response.Message = "vendor address not found"
			return &response, nil
		}
		response.Message = fmt.Sprint("error in vendor address service VerifyAddress API1: ", err.Error())
		logger.FromContext(ctx).Error(response.Message)
		return &response, nil
	}
	err = query.Update("verification_status", utils.Verified).Error
	if err != nil {
		response.Message = fmt.Sprint("not able to update verification status: ", err.Error())
		logger.FromContext(ctx).Error(response.Message)
		return &response, nil
	}
	sellerID := vendorAddress.SellerID
	seller := models.Seller{}
	err = database.DBAPM(ctx).Model(&models.Seller{}).Where("id = ?", sellerID).Scan(&seller).Error
	if err != nil {
		response.Message = fmt.Sprint("error in vendor address service VerifyAddress API2: ", err.Error())
		logger.FromContext(ctx).Error(response.Message)
		return &response, nil
	}
	sellerActivityLog := models.SellerActivityLog{
		SellerID: seller.UserID,
		Action:   "verify_address",
		Notes:    `{"status":"verified"}`,
	}
	if uid := utils.GetCurrentUserID(ctx); uid != nil {
		sellerActivityLog.UserID = *uid
	}
	err = database.DBAPM(ctx).Model(&models.SellerActivityLog{}).Create(&sellerActivityLog).Error
	if err != nil {
		response.Message = fmt.Sprint("error in vendor address service VerifyAddress API3: ", err.Error())
		logger.FromContext(ctx).Error(response.Message)
		return &response, nil
	}
	response.Status = utils.Success
	response.Message = "vendor address verified successfully"
	return &response, nil
}
