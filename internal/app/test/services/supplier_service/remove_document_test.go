package supplier_service_test

import (
	"context"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("EditSupplier", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		test_utils.SetPermission(&ctx, []string{"supplierpanel:editverifiedblockedsupplieronly:admin"})
	})

	Context("Removing supplier document", func() {
		It("Should remove primary document successfully", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				AgreementUrl: "abc/xyz.pdf",
				Status:       models.SupplierStatusVerified,
			})
			param := &supplierpb.RemoveDocumentParam{
				Id:           supplier.ID,
				DocumentType: "agreement_url",
			}
			res, err := new(services.SupplierService).RemoveDocument(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier agreement_url Removed Successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)

			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusPending))
			Expect(updatedSupplier.AgreementUrl).To(Equal(""))
		})

		It("Should remove secondary document successfully", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				GuarantorImageUrl: "abc/xyz.jpg",
				Status:            models.SupplierStatusVerified,
			})
			param := &supplierpb.RemoveDocumentParam{
				Id:           supplier.ID,
				DocumentType: "guarantor_image_url",
			}
			res, err := new(services.SupplierService).RemoveDocument(ctx, param)
			log.Println("======>", res.Message)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier guarantor_image_url Removed Successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)

			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusVerified))
			Expect(updatedSupplier.GuarantorImageUrl).To(Equal(""))
		})

		It("Should return error for invalid document type", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				AgreementUrl: "abc/xyz.pdf",
			})
			param := &supplierpb.RemoveDocumentParam{
				Id:           supplier.ID,
				DocumentType: "agreement_url_abc",
			}
			res, err := new(services.SupplierService).RemoveDocument(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Invalid Document Type"))
		})

		It("Should return error for un-allowed permission", func() {
			test_utils.SetPermission(&ctx, []string{"per:missi:on"})
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				AgreementUrl: "abc/xyz.pdf",
				Status:       models.SupplierStatusVerified,
			})
			param := &supplierpb.RemoveDocumentParam{
				Id:           supplier.ID,
				DocumentType: "agreement_url",
			}
			res, err := new(services.SupplierService).RemoveDocument(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Change Not Allowed"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)

			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusVerified))
			Expect(updatedSupplier.AgreementUrl).To(Equal("abc/xyz.pdf"))
		})
	})
})
