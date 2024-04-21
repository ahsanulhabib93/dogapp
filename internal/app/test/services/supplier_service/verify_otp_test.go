package supplier_service_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	eventBus "github.com/voonik/goConnect/api/go/event_bus/publisher"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	aaaMocks "github.com/voonik/goFramework/pkg/aaa/models/mocks"
	"github.com/voonik/ss2/internal/app/publisher"
	mockPublisher "github.com/voonik/ss2/internal/app/publisher/mocks"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	otpPb "github.com/voonik/goConnect/api/go/vigeon2/otp"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("VerifyOtp", func() {
	var ctx context.Context
	var apiHelperInstance *mocks.APIHelperInterface
	var mockAudit *mocks.AuditLogMock
	var appPreferenceMockInstance *aaaMocks.AppPreferenceInterface

	BeforeEach(func() {
		test_utils.GetContext(&ctx)

		mockAudit = mocks.SetAuditLogMock()
		mockAudit.On("RecordAuditAction", ctx, mock.Anything).Return(nil)

		appPreferenceMockInstance = new(aaaMocks.AppPreferenceInterface)
		aaaModels.InjectMockAppPreferenceServiceInstance(appPreferenceMockInstance)
		appPreferenceMockInstance.On("GetValue", ctx, "should_send_supplier_log", "true").Return("true")
	})

	AfterEach(func() {
		aaaModels.InjectMockAppPreferenceServiceInstance(nil)
	})

	Context("For Invalid Supplier", func() {
		It("Should return error", func() {
			param := &supplierpb.VerifyOtpParam{SupplierId: 100, OtpCode: "1234"}
			t := &testing.T{}
			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).VerifyOtp(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Not Found"))
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("For Supplier with phone verified already", func() {
		It("Should return error", func() {
			isPhoneVerified := true
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{IsPhoneVerified: &isPhoneVerified})
			param := &supplierpb.VerifyOtpParam{SupplierId: supplier.ID, OtpCode: "1234"}

			t := &testing.T{}
			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).VerifyOtp(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Phone number is already verified"))
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("For valid Supplier", func() {
		var supplier *models.Supplier

		BeforeEach(func() {
			apiHelperInstance = new(mocks.APIHelperInterface)
			helpers.InjectMockAPIHelperInstance(apiHelperInstance)
			supplier = test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &otpPb.VerifyOtpParam{
				Service:    "ss2",
				SourceType: "SupplierVerification",
				SourceId:   supplier.ID,
				OtpCode:    "1234",
			}
			resp := &otpPb.OtpResponse{
				Success: true,
				Message: "Verified OTP successfully",
				Uuid:    "9876765",
			}
			apiHelperInstance.On("VerifyOtpAPI", ctx, param).Return(resp)
		})

		It("Should call vigeon service and verify OTP", func() {
			param := &supplierpb.VerifyOtpParam{SupplierId: supplier.ID, OtpCode: "1234"}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).VerifyOtp(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Verified OTP successfully"))

			updatedSupplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(*updatedSupplier.IsPhoneVerified).To(Equal(true))
			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(1))
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("When Vigeon returned error while verifying otp", func() {
		var supplier *models.Supplier
		BeforeEach(func() {
			apiHelperInstance = new(mocks.APIHelperInterface)
			helpers.InjectMockAPIHelperInstance(apiHelperInstance)
			supplier = test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &otpPb.VerifyOtpParam{
				Service:    "ss2",
				SourceType: "SupplierVerification",
				SourceId:   supplier.ID,
				OtpCode:    "1234",
			}
			resp := &otpPb.OtpResponse{
				Success: false,
				Message: "Invalid OTP",
			}
			apiHelperInstance.On("VerifyOtpAPI", ctx, param).Return(resp)
		})

		It("Should return error", func() {
			param := &supplierpb.VerifyOtpParam{SupplierId: supplier.ID, OtpCode: "1234"}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).VerifyOtp(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Invalid OTP"))

			updatedSupplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(*updatedSupplier.IsPhoneVerified).To(Equal(false))
			mockedEventBus.AssertExpectations(t)
		})
	})
})
