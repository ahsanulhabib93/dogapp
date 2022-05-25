package supplier_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("UpdateStatus", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
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
		})
	})

	Context("Updating without reason", func() {
		It("Should be updated", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusBlocked),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier status updated successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusBlocked))
			Expect(updatedSupplier.Reason).To(Equal(""))
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
		It("Should return error", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusVerified),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Required details for verification are not present"))
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

	Context("Updating with status for which transition not allowed", func() {
		It("Should return error", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusFailed})
			param := &supplierpb.UpdateStatusParam{
				Id:     supplier.ID,
				Status: string(models.SupplierStatusBlocked),
			}
			res, err := new(services.SupplierService).UpdateStatus(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Status transition not allowed"))
		})
	})
})
