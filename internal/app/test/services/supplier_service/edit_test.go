package supplier_service_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("EditSupplier", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Editing existing Supplier", func() {
		It("Should update supplier and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 1},
					{CategoryID: 2},
					{CategoryID: 3},
				},
				SupplierSaMappings: []models.SupplierSaMapping{
					{SourcingAssociateId: 4},
					{SourcingAssociateId: 5},
				},
			})
			param := &supplierpb.SupplierObject{
				Id:           supplier.ID,
				Name:         "Name",
				Email:        "Email",
				SupplierType: uint64(utils.L1),
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierCategoryMappings").Preload("SupplierSaMappings").
				First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Email).To(Equal(param.Email))
			Expect(updatedSupplier.Name).To(Equal(param.Name))
			Expect(updatedSupplier.SupplierType).To(Equal(utils.L1))
			Expect(len(updatedSupplier.SupplierCategoryMappings)).To(Equal(3))
			Expect(len(updatedSupplier.SupplierSaMappings)).To(Equal(2))
			Expect(updatedSupplier.SupplierCategoryMappings[1].CategoryID).To(Equal(uint64(2)))
		})
	})

	Context("Editing only one field of existing Supplier", func() {
		It("Should update supplier name and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.SupplierObject{
				Id:   supplier.ID,
				Name: "Name",
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Email).To(Equal(supplier.Email))
			Expect(updatedSupplier.SupplierType).To(Equal(utils.Hlc))
			Expect(updatedSupplier.Name).To(Equal(param.Name))
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusPending))
		})

		It("Should update supplier status and return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				Name: "test-supplier",
			})
			param := &supplierpb.SupplierObject{
				Id:     supplier.ID,
				Status: models.SupplierStatusActive,
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Email).To(Equal(supplier.Email))
			Expect(updatedSupplier.SupplierType).To(Equal(utils.Hlc))
			Expect(updatedSupplier.Name).To(Equal(supplier.Name))
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusActive))
		})

		It("Should return error on invalid status", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.SupplierObject{
				Id:     supplier.ID,
				Status: "no idea",
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while updating Supplier: Status should be Active/Pending"))
		})
	})

	Context("Editing invalid supplier", func() {
		It("Should return error response", func() {
			param := &supplierpb.SupplierObject{Id: 1000}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Not Found"))
		})
	})

	Context("Editing with other supplier name", func() {
		It("Should return error response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.SupplierObject{
				Id:   supplier1.ID,
				Name: supplier2.Name,
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while updating Supplier: Name should be unique"))
		})
	})

	Context("Editing with category ids", func() {
		It("Should delete old mapping", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 101},
					{CategoryID: 201},
				},
			})
			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierCategoryMappings").First(&updatedSupplier, supplier.ID)
			Expect(len(updatedSupplier.SupplierCategoryMappings)).To(Equal(2))

			param := &supplierpb.SupplierObject{
				Id:           supplier.ID,
				Name:         "Name",
				Email:        "Email",
				SupplierType: uint64(utils.L1),
				CategoryIds:  []uint64{101, 102, 100},
			}
			res, err := new(services.SupplierService).Edit(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier = models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierCategoryMappings").First(&updatedSupplier, supplier.ID)
			Expect(len(updatedSupplier.SupplierCategoryMappings)).To(Equal(3))
			Expect(updatedSupplier.SupplierCategoryMappings[0].CategoryID).To(Equal(uint64(101)))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Unscoped().Where("supplier_category_mappings.supplier_id = ?", supplier.ID).Count(&count)
			Expect(count).To(Equal(4))
		})

		It("Should restore deleted mapping", func() {
			t := time.Now()
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 101},
					{CategoryID: 200},
					{
						CategoryID: 567,
						DeletedAt:  &t,
					},
				},
			})
			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierCategoryMappings").First(&updatedSupplier, supplier.ID)
			Expect(len(updatedSupplier.SupplierCategoryMappings)).To(Equal(2))

			param := &supplierpb.SupplierObject{
				Id:           supplier.ID,
				Name:         "Name",
				Email:        "Email",
				SupplierType: uint64(utils.L1),
				CategoryIds:  []uint64{101, 200, 567},
			}
			res, err := new(services.SupplierService).Edit(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier = models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierCategoryMappings").First(&updatedSupplier, supplier.ID)
			Expect(len(updatedSupplier.SupplierCategoryMappings)).To(Equal(3))
			Expect(utils.Int64Min(updatedSupplier.SupplierCategoryMappings[0].CategoryID,
				utils.Int64Min(updatedSupplier.SupplierCategoryMappings[1].CategoryID,
					updatedSupplier.SupplierCategoryMappings[2].CategoryID))).To(Equal(uint64(101)))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Unscoped().Where("supplier_category_mappings.supplier_id = ?", supplier.ID).Count(&count)
			Expect(count).To(Equal(3))
		})
	})
	Context("Editing with Sourcing Associate Id ids", func() {
		It("Should delete old mapping", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierSaMappings: []models.SupplierSaMapping{
					{SourcingAssociateId: 202},
					{SourcingAssociateId: 203},
				},
			})
			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierSaMappings").First(&updatedSupplier, supplier.ID)
			Expect(len(updatedSupplier.SupplierSaMappings)).To(Equal(2))

			param := &supplierpb.SupplierObject{
				Id:           supplier.ID,
				Name:         "Name",
				Email:        "Email",
				SupplierType: uint64(utils.L1),
				SaIds:        []uint64{503, 504, 505},
			}
			res, err := new(services.SupplierService).Edit(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier = models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierSaMappings").First(&updatedSupplier, supplier.ID)
			Expect(len(updatedSupplier.SupplierSaMappings)).To(Equal(3))
			Expect(utils.Int64Min(updatedSupplier.SupplierSaMappings[0].SourcingAssociateId,
				utils.Int64Min(updatedSupplier.SupplierSaMappings[1].SourcingAssociateId,
					updatedSupplier.SupplierSaMappings[2].SourcingAssociateId))).To(Equal(uint64(503)))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierSaMapping{}).Unscoped().Where("supplier_sa_mappings.supplier_id = ?", supplier.ID).Count(&count)
			// so the data will be 202, 203, 503, 504, 505
			Expect(count).To(Equal(5))
		})

		It("Should restore deleted mapping", func() {
			t := time.Now()
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{

				SupplierSaMappings: []models.SupplierSaMapping{
					{SourcingAssociateId: 10001},
					{SourcingAssociateId: 10002},
					{
						SourcingAssociateId: 10009,
						DeletedAt:           &t,
					},
				},
			})
			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierSaMappings").First(&updatedSupplier, supplier.ID)
			Expect(len(updatedSupplier.SupplierSaMappings)).To(Equal(2))

			param := &supplierpb.SupplierObject{
				Id:           supplier.ID,
				Name:         "Name",
				Email:        "Email",
				SupplierType: uint64(utils.L1),
				CategoryIds:  []uint64{101, 200, 567},
				SaIds:        []uint64{10001, 10002, 10009},
			}
			res, err := new(services.SupplierService).Edit(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier = models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierSaMappings").First(&updatedSupplier, supplier.ID)
			Expect(len(updatedSupplier.SupplierSaMappings)).To(Equal(3))
			Expect(updatedSupplier.SupplierSaMappings[2].SourcingAssociateId).To(Equal(uint64(10009)))
			var count int
			database.DBAPM(ctx).Model(&models.SupplierSaMapping{}).Unscoped().Where("supplier_sa_mappings.supplier_id = ?", supplier.ID).Count(&count)
			Expect(count).To(Equal(3))
		})
	})
})
