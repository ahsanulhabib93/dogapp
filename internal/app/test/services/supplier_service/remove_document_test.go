package supplier_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("RemoveDocument", func() {
	var ctx context.Context
	var mockAudit *mocks.AuditLogMock

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		ctx = test_helper.SetContextUser(ctx, 101, []string{"supplierpanel:editverifiedblockedsupplieronly:admin"})

		mocks.SetAuditLogMock()
		mockAudit = mocks.SetAuditLogMock()
		mockAudit.On("RecordAuditAction", ctx, mock.Anything).Return(nil)
	})

	AfterEach(func() {
		mocks.UnsetAuditLogMock()
	})

	Context("Removing supplier document", func() {
		It("Should remove primary document successfully", func() {
			supplierData := models.Supplier{
				Status: models.SupplierStatusVerified,
				PartnerServiceMappings: []models.PartnerServiceMapping{{
					AgreementUrl: "abc/xyz.pdf",
				}},
			}
			supplier := test_helper.CreateSupplier(ctx, &supplierData)

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
			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(1))
		})

		It("Should remove secondary document successfully", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				GuarantorNidFrontImageUrl: "abc/xyz.jpg",
				Status:                    models.SupplierStatusVerified,
			})
			param := &supplierpb.RemoveDocumentParam{
				Id:           supplier.ID,
				DocumentType: "guarantor_nid_front_image_url",
			}
			res, err := new(services.SupplierService).RemoveDocument(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier guarantor_nid_front_image_url Removed Successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)

			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusVerified))
			Expect(updatedSupplier.GuarantorNidFrontImageUrl).To(Equal(""))
			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(1))
		})

		It("Should return error for invalid document type", func() {
			supplierData := models.Supplier{
				Status: models.SupplierStatusVerified,
				PartnerServiceMappings: []models.PartnerServiceMapping{{
					AgreementUrl: "abc/xyz.pdf",
				}},
			}
			supplier := test_helper.CreateSupplier(ctx, &supplierData)
			param := &supplierpb.RemoveDocumentParam{
				Id:           supplier.ID,
				DocumentType: "agreement_url_abc",
			}
			res, err := new(services.SupplierService).RemoveDocument(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Invalid Document Type"))
			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(0))
		})

		It("Should return error for un-allowed permission", func() {
			test_utils.SetPermission(&ctx, []string{"per:missi:on"})
			supplierData := models.Supplier{
				Status: models.SupplierStatusVerified,
				PartnerServiceMappings: []models.PartnerServiceMapping{{
					AgreementUrl: "abc/xyz.pdf",
				}},
			}
			supplier := test_helper.CreateSupplier(ctx, &supplierData)

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
			partnerServices := []*models.PartnerServiceMapping{{}}
			database.DBAPM(ctx).Model(supplier).Association("PartnerServiceMappings").Find(&partnerServices)
			Expect(len(partnerServices)).To(Equal(1))
			Expect(partnerServices[0].AgreementUrl).To(Equal("abc/xyz.pdf"))
		})
	})
})
