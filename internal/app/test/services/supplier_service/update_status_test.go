package supplier_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/goFramework/pkg/misc"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("UpdateStatus", func() {
	var ctx context.Context
	var userId uint64 = uint64(101)

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		aaaModels.CreateAppPreferenceServiceInterface()

		threadObject := &misc.ThreadObject{
			VaccountId:    1,
			PortalId:      1,
			CurrentActId:  1,
			XForwardedFor: "5079327",
			UserData: &misc.UserData{
				UserId: userId,
				Name:   "John",
				Email:  "john@gmail.com",
				Phone:  "8801855533367",
			},
		}
		ctx = misc.SetInContextThreadObject(ctx, threadObject)
	})

	Context("Update Supplier status", func() {
		It("Should be updated and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
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
		})

		It("Updating status to block with reason reason", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusBlocked),
				Reason: "no reason",
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier status updated successfully"))
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
