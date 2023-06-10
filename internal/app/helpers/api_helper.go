package helpers

import (
	"context"
	"log"

	userPb "github.com/voonik/goConnect/api/go/cre_admin/users_detail"
	employeePb "github.com/voonik/goConnect/api/go/sr_service/attendance"
	otpPb "github.com/voonik/goConnect/api/go/vigeon2/otp"
	userSrv "github.com/voonik/goConnect/cre_admin/users_detail"
	employeeSrv "github.com/voonik/goConnect/sr_service/attendance"
	Vigeon2Service "github.com/voonik/goConnect/vigeon2/otp"
	"github.com/voonik/ss2/internal/app/utils"
)

// APIHelper ...
type APIHelper struct{}

// APIHelperInterface ...
type APIHelperInterface interface {
	SendOtpAPI(context.Context, otpPb.OtpParam) *otpPb.OtpResponse
	VerifyOtpAPI(context.Context, otpPb.VerifyOtpParam) *otpPb.OtpResponse
	FindUserByPhone(context.Context, string) *userPb.UserInfo
	FindTalentXUserByPhone(context.Context, string) []*employeePb.EmployeeRecord
}

var apiHelper APIHelperInterface

// InjectMockAPIHelperInstance ...
func InjectMockAPIHelperInstance(mockObj APIHelperInterface) {
	apiHelper = mockObj
}

// getAPIHelperInstance ...
func getAPIHelperInstance() APIHelperInterface {
	if apiHelper == nil {
		return new(APIHelper)
	}
	return apiHelper
}

//SendOtpAPI ...
func SendOtpAPI(ctx context.Context, supplierID uint64, phone string, content string, resend bool) *otpPb.OtpResponse {
	otpParam := otpPb.OtpParam{
		Service:    "ss2",
		SourceType: "SupplierVerification",
		SourceId:   supplierID,
		Phone:      phone,
		Content:    content,
		Resend:     resend,
	}
	return getAPIHelperInstance().SendOtpAPI(ctx, otpParam)
}

//SendOtpAPI ...
func (apiHelper *APIHelper) SendOtpAPI(ctx context.Context, otpParam otpPb.OtpParam) *otpPb.OtpResponse {
	resp, _ := Vigeon2Service.Otp().CreateOtp(ctx, &otpParam)
	return resp
}

//VerifyOtpAPI ...
func VerifyOtpAPI(ctx context.Context, supplierID uint64, otpCode string) *otpPb.OtpResponse {
	verifyOtpParam := otpPb.VerifyOtpParam{
		Service:    "ss2",
		SourceType: "SupplierVerification",
		SourceId:   supplierID,
		OtpCode:    otpCode,
	}
	return getAPIHelperInstance().VerifyOtpAPI(ctx, verifyOtpParam)
}

//VerifyOtpAPI ...
func (apiHelper *APIHelper) VerifyOtpAPI(ctx context.Context, verifyOtpParam otpPb.VerifyOtpParam) *otpPb.OtpResponse {
	resp, _ := Vigeon2Service.Otp().VerifyOtp(ctx, &verifyOtpParam)
	return resp
}

//FindCreUserByPhone ...
func FindCreUserByPhone(ctx context.Context, phone string) *userPb.UserInfo {
	log.Printf("FindCreUserByPhone: phone = %s\n", phone)
	if utils.IsEmptyStr(phone) {
		return nil
	}

	return getAPIHelperInstance().FindUserByPhone(ctx, phone)
}

func GetTalentXUser(ctx context.Context, phone string) []*employeePb.EmployeeRecord {
	log.Printf("GetIdentityUser: phone = %s\n", phone)
	if utils.IsEmptyStr(phone) {
		return nil
	}

	return getAPIHelperInstance().FindTalentXUserByPhone(ctx, phone)
}

//FindUserByPhone ...
func (apiHelper *APIHelper) FindUserByPhone(ctx context.Context, phone string) *userPb.UserInfo {
	resp, _ := userSrv.UsersDetail().FindByPhone(ctx, &userPb.UserParams{Phone: phone})
	log.Printf("FindUserByPhone: phone = %s response = %v\n", phone, resp)
	return resp.Data
}

//FindTalentXUserByPhone ...
func (apiHelper *APIHelper) FindTalentXUserByPhone(ctx context.Context, phone string) []*employeePb.EmployeeRecord {
	resp, _ := employeeSrv.Attendance().ListEmployee(ctx, &employeePb.ListEmployeeParams{Phone: phone, IgnoreWarehouseFilter: true})
	log.Printf("FindTalentXUserByPhone: phone = %s response = %v\n", phone, resp)
	return resp.Data
}
