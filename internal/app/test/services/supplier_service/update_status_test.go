package supplier_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goConnect/api/go/vigeon/notify"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/goFramework/pkg/rest"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"google.golang.org/grpc/metadata"
)

var _ = Describe("UpdateStatus", func() {
	var ctx context.Context
	var userId uint64 = uint64(101)
	var mock1 *mocks.ApiCallHelperInterface
	var mock2 *mocks.VigeonAPIHelperInterface
	var mockAudit *mocks.AuditLogMock

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		aaaModels.CreateAppPreferenceServiceInterface()

		header := map[string]string{"authorization": "random"}
		ctx = test_helper.SetContextUser(ctx, userId, []string{})
		ctx = metadata.NewIncomingContext(ctx, metadata.New(header))

		mock1, mock2, mockAudit = mocks.SetApiCallerMock(), mocks.SetVigeonAPIHelperMock(), mocks.SetAuditLogMock()
		mockAudit.On("RecordAuditAction", ctx, mock.Anything).Return(nil)
		mock1.On("Get", ctx, mock.Anything, mock.Anything).Return(&rest.Response{Body: "{\"data\":{\"users\":[{\"id\":101,\"email\":\"user_email@gmail.com\"}]}}"}, nil)
		mock2.On("SendEmailAPI", ctx, mock.Anything).Return(&notify.EmailResp{}, nil)
	})

	AfterEach(func() {
		mocks.UnsetApiCallerMock()
		mocks.UnsetVigeonHelperMock()
		mocks.UnsetAuditLogMock()
	})

	Context("Update Supplier status", func() {
		It("Should be updated and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				UserID: &userId,
			})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusFailed),
				Reason: "test reason",
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier status updated successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusFailed))
			Expect(updatedSupplier.Reason).To(Equal(param.Reason))
			Expect(updatedSupplier.AgentID).To(BeNil())
			Expect(mock1.Count["Get"]).To(Equal(1))
			Expect(mock2.Count["SendEmailAPI"]).To(Equal(1))
			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(1))
		})

		It("Should update status for blocked user", func() {
			isPhoneVerified := true
			supplier := test_helper.CreateSupplierWithAddress(ctx, &models.Supplier{
				IsPhoneVerified: &isPhoneVerified,
				Status:          models.SupplierStatusBlocked,
			})
			test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusVerified),
				Reason: "test reason",
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier status updated successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusVerified))
			Expect(updatedSupplier.Reason).To(Equal(param.Reason))
			Expect(*updatedSupplier.AgentID).To(Equal(userId))
			Expect(mock1.Count["Get"]).To(Equal(0))
			Expect(mock2.Count["SendEmailAPI"]).To(Equal(0))
			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(1))
		})

		It("Updating status to block with reason reason", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				UserID: &userId,
			})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusBlocked),
				Reason: "no reason",
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier status updated successfully"))
			Expect(mock1.Count["Get"]).To(Equal(1))
			Expect(mock2.Count["SendEmailAPI"]).To(Equal(1))
		})
	})

	Context("Update Supplier status as Verified", func() {
		It("Should be updated and return success response", func() {
			isPhoneVerified := true
			supplier := test_helper.CreateSupplierWithAddress(ctx, &models.Supplier{IsPhoneVerified: &isPhoneVerified})
			test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusVerified),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier status updated successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusVerified))
			Expect(*updatedSupplier.AgentID).To(Equal(userId))
		})
	})

	Context("Updating invalid supplier", func() {
		It("Should return error response", func() {
			param := &supplierpb.UpdateStatusParam{
				Id:     1000,
				Status: string(models.SupplierStatusVerified),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Not Found"))
			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(0))
		})
	})

	Context("Updating invalid status", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: "Test",
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Invalid Status"))
		})
	})

	Context("Updating without status", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: "",
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Invalid Status"))
		})
	})

	Context("Update Supplier status as Verified without required details", func() {
		It("Should return error for missing payment account", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusVerified),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("At least one payment account details should be present"))
		})

		It("Should return error for missing address", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusVerified),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("At least one supplier address should be present"))
		})
	})

	Context("Update same status", func() {
		It("Should return error", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusPending),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Status transition not allowed"))
		})
	})

	Context("Should return error", func() {
		It("Updating with status for which transition not allowed", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusFailed})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusBlocked),
				Reason: "here take the reason",
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Status transition not allowed"))
		})

		It("Updating status to block without reason", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusBlocked),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Status change reason missing"))
		})

		It("When no OTP verification or primary document given", func() {
			supplier := test_helper.CreateSupplierWithAddress(ctx, &models.Supplier{
				Status: models.SupplierStatusBlocked,
			})
			test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusVerified),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("At least one primary document or OTP verification needed"))
		})

		It("When at least one primary document required for given supplier type", func() {
			aaaModels.InjectMockAppPreferenceServiceInstance(mocks.GetAppPreferenceMock(map[string]interface{}{
				"enabled_primary_doc_verification": []string{"Hlc"},
			}))

			isPhoneVerified := true
			supplier := test_helper.CreateSupplierWithAddress(ctx, &models.Supplier{
				IsPhoneVerified: &isPhoneVerified,
				Status:          models.SupplierStatusBlocked,
			})
			test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusVerified),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("At least one primary document required for supplier type: Hlc"))
		})

		It("When otp verification required for given supplier type", func() {
			aaaModels.InjectMockAppPreferenceServiceInstance(mocks.GetAppPreferenceMock(map[string]interface{}{
				"enabled_otp_verification": []string{"Hlc"},
			}))

			supplier := test_helper.CreateSupplierWithAddress(ctx, &models.Supplier{
				Status:    models.SupplierStatusBlocked,
				NidNumber: "1234567890",
			})
			test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusVerified),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("OTP verification required for supplier type: Hlc"))
		})
	})
})
