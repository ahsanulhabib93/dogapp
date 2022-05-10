package helpers

import (
	"context"

	otpPb "github.com/voonik/goConnect/api/go/vigeon2/otp"
	Vigeon2Service "github.com/voonik/goConnect/vigeon2/otp"
)

// APIHelper ...
type APIHelper struct{}

// APIHelperInterface ...
type APIHelperInterface interface {
	SendOtpAPI(context.Context, otpPb.OtpParam) *otpPb.OtpResponse
	VerifyOtpAPI(context.Context, otpPb.VerifyOtpParam) *otpPb.OtpResponse
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
