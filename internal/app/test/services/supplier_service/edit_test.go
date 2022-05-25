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
		utils.SetMockPermissions([]string{models.AllowedPermission})
	})

	Context("Editing existing Supplier", func() {
		It("Should update supplier and return success response", func() {
			isPhoneVerified := true
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				IsPhoneVerified: &isPhoneVerified,
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 1},
					{CategoryID: 2},
					{CategoryID: 3},
				},
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 4},
					{ProcessingCenterID: 5},
				},
			})
			param := &supplierpb.SupplierObject{
				Id:             supplier.ID,
				Name:           "Name",
				Email:          "Email",
				SupplierType:   uint64(utils.L1),
				BusinessName:   "BusinessName",
				Phone:          "01123456789",
				AlternatePhone: "01123456111",
				ShopImageUrl:   "ss2/shop_images/test.png",
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierCategoryMappings").Preload("SupplierOpcMappings").
				First(&updatedSupplier, supplier.ID)

			Expect(updatedSupplier.Email).To(Equal(param.Email))
			Expect(updatedSupplier.Name).To(Equal(param.Name))
			Expect(updatedSupplier.SupplierType).To(Equal(utils.L1))
			Expect(updatedSupplier.BusinessName).To(Equal(param.BusinessName))
			Expect(updatedSupplier.Phone).To(Equal(param.Phone))
			Expect(updatedSupplier.AlternatePhone).To(Equal(param.AlternatePhone))
			Expect(updatedSupplier.ShopImageURL).To(Equal(param.ShopImageUrl))
			Expect(*updatedSupplier.IsPhoneVerified).To(Equal(false))
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusPending))

			Expect(len(updatedSupplier.SupplierCategoryMappings)).To(Equal(3))
			Expect(len(updatedSupplier.SupplierOpcMappings)).To(Equal(2))
			Expect(updatedSupplier.SupplierCategoryMappings[1].CategoryID).To(Equal(uint64(2)))
		})
	})

	Context("Editing only one field of existing Supplier", func() {
		It("Should update supplier name and return success response", func() {
			isPhoneVerified := true
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				IsPhoneVerified: &isPhoneVerified,
				Status:          models.SupplierStatusBlocked,
			})
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
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusBlocked))
			Expect(*updatedSupplier.IsPhoneVerified).To(Equal(true))
		})
	})

	Context("Editing allowed for limited permission", func() {
		It("Should return success on updating pending supplier", func() {
			utils.SetMockPermissions([]string{})
			isPhoneVerified := true
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				IsPhoneVerified: &isPhoneVerified,
				Status:          models.SupplierStatusPending,
			})
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
			Expect(*updatedSupplier.IsPhoneVerified).To(Equal(true))
		})

		It("Should return error on updating verified supplier", func() {
			utils.SetMockPermissions([]string{})
			isPhoneVerified := true
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				IsPhoneVerified: &isPhoneVerified,
				Status:          models.SupplierStatusVerified,
			})
			param := &supplierpb.SupplierObject{
				Id:   supplier.ID,
				Name: "Name",
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Update Not Allowed"))
		})

		It("Should return error on updating blocked supplier", func() {
			utils.SetMockPermissions([]string{})
			isPhoneVerified := true
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				IsPhoneVerified: &isPhoneVerified,
				Status:          models.SupplierStatusBlocked,
			})
			param := &supplierpb.SupplierObject{
				Id:   supplier.ID,
				Name: "Name",
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Update Not Allowed"))
		})
	})

	Context("Editing Supplier details in Verified status", func() {
		It("Should update supplier details and update status as Pending and return success response", func() {
			isPhoneVerified := true
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				IsPhoneVerified: &isPhoneVerified,
				Status:          models.SupplierStatusVerified,
			})
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
			Expect(*updatedSupplier.IsPhoneVerified).To(Equal(true))
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

	Context("Editing with new set of category ids", func() {
		It("Should delete old mapping and add new mapping", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 101},
					{CategoryID: 201},
				},
			})

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

			var categoryIds []uint64
			database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Pluck("category_id", &categoryIds)
			Expect(len(categoryIds)).To(Equal(3))
			Expect(categoryIds).To(ContainElements([]uint64{100, 101, 102}))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Unscoped().Where("supplier_category_mappings.supplier_id = ?", supplier.ID).Count(&count)
			Expect(count).To(Equal(4))
		})
	})

	Context("Editing with new set of category ids which got removed before", func() {
		It("Should restore deleted mapping", func() {
			t := time.Now()
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 101},
					{CategoryID: 200},
					{CategoryID: 567, DeletedAt: &t},
				},
			})

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

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierCategoryMappings").First(&updatedSupplier, supplier.ID)

			categoryMappings := updatedSupplier.SupplierCategoryMappings
			Expect(len(categoryMappings)).To(Equal(3))
			Expect(categoryMappings[0].CategoryID).To(Equal(uint64(101)))
			Expect(categoryMappings[1].CategoryID).To(Equal(uint64(200)))
			Expect(categoryMappings[2].CategoryID).To(Equal(uint64(567)))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Unscoped().Where("supplier_category_mappings.supplier_id = ?", supplier.ID).Count(&count)
			Expect(count).To(Equal(3))
		})
	})

	Context("Editing with invalid phone number", func() {
		It("Should return error response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.SupplierObject{
				Id:    supplier1.ID,
				Phone: "1234",
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while updating Supplier: Invalid Phone Number"))
		})
	})

	Context("Editing Supplier with duplicate phone number", func() {
		It("Should return error response", func() {
			test_helper.CreateSupplier(ctx, &models.Supplier{Phone: "8801234567890"})
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{Phone: "8801234567800"})
			param := &supplierpb.SupplierObject{
				Id:    supplier1.ID,
				Phone: "8801234567890",
			}
			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while updating Supplier: Supplier Already Exists"))
		})
	})
})
