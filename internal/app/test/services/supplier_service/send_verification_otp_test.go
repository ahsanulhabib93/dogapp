package supplier_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	otpPb "github.com/voonik/goConnect/api/go/vigeon2/otp"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("SendVerificationOtp", func() {
	var ctx context.Context
	var apiHelperInstance *mocks.APIHelperInterface

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("For Invalid Supplier", func() {
		It("Should return error", func() {
			param := &supplierpb.SendOtpParam{SupplierId: 100}
			res, err := new(services.SupplierService).SendVerificationOtp(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Not Found"))
		})
	})

	Context("For Supplier with phone verified already", func() {
		It("Should return error", func() {
			isPhoneVerified := true
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{IsPhoneVerified: &isPhoneVerified})
			param := &supplierpb.SendOtpParam{SupplierId: supplier.ID}
			res, err := new(services.SupplierService).SendVerificationOtp(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Phone number is already verified"))
		})
	})

	Context("For valid Supplier", func() {
		var supplier *models.Supplier
		BeforeEach(func() {
			apiHelperInstance = new(mocks.APIHelperInterface)
			helpers.InjectMockAPIHelperInstance(apiHelperInstance)
			supplier = test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := otpPb.OtpParam{
				Service:    "ss2",
				SourceType: "SupplierVerification",
				SourceId:   supplier.ID,
				Phone:      supplier.Phone,
				Content:    "OTP for supplier verification: $otp",
				Resend:     false,
			}
			resp := otpPb.OtpResponse{
				Success: true,
				Message: "OTP created and sent successfully",
				Uuid:    "1234",
			}
			apiHelperInstance.On("SendOtpAPI", ctx, &param).Return(&resp)
		})

		It("Should call vigeon service to send OTP", func() {
			param := &supplierpb.SendOtpParam{SupplierId: supplier.ID}
			res, err := new(services.SupplierService).SendVerificationOtp(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("OTP created and sent successfully"))
		})
	})

	Context("Resend OTP for valid Supplier ", func() {
		var supplier *models.Supplier
		BeforeEach(func() {
			apiHelperInstance = new(mocks.APIHelperInterface)
			helpers.InjectMockAPIHelperInstance(apiHelperInstance)
			supplier = test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := otpPb.OtpParam{
				Service:    "ss2",
				SourceType: "SupplierVerification",
				SourceId:   supplier.ID,
				Phone:      supplier.Phone,
				Content:    "OTP for supplier verification: $otp",
				Resend:     true,
			}
			resp := otpPb.OtpResponse{
				Success: true,
				Message: "OTP created and sent successfully",
				Uuid:    "1234",
			}
			apiHelperInstance.On("SendOtpAPI", ctx, &param).Return(&resp)
		})

		It("Should call vigeon service to resend otp", func() {
			param := &supplierpb.SendOtpParam{SupplierId: supplier.ID, Resend: true}
			res, err := new(services.SupplierService).SendVerificationOtp(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("OTP created and sent successfully"))
		})
	})

	Context("When Vigeon returned error while sending otp", func() {
		var supplier *models.Supplier
		BeforeEach(func() {
			apiHelperInstance = new(mocks.APIHelperInterface)
			helpers.InjectMockAPIHelperInstance(apiHelperInstance)
			supplier = test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := otpPb.OtpParam{
				Service:    "ss2",
				SourceType: "SupplierVerification",
				SourceId:   supplier.ID,
				Phone:      supplier.Phone,
				Content:    "OTP for supplier verification: $otp",
				Resend:     false,
			}
			resp := otpPb.OtpResponse{
				Success: false,
				Message: "Error while creating OTP",
			}
			apiHelperInstance.On("SendOtpAPI", ctx, &param).Return(&resp)
		})

		It("Should return error", func() {
			param := &supplierpb.SendOtpParam{SupplierId: supplier.ID}
			res, err := new(services.SupplierService).SendVerificationOtp(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating OTP"))
		})
	})
})
