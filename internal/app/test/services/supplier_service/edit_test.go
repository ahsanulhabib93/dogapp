package supplier_service_test

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	eventBus "github.com/voonik/goConnect/api/go/event_bus/publisher"
	"github.com/voonik/ss2/internal/app/publisher"
	mockPublisher "github.com/voonik/ss2/internal/app/publisher/mocks"

	categoryPb "github.com/voonik/goConnect/api/go/cmt/category"
	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	aaaMocks "github.com/voonik/goFramework/pkg/aaa/models/mocks"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("EditSupplier", func() {
	var ctx context.Context
	var mockAudit *mocks.AuditLogMock
	var appPreferenceMockInstance *aaaMocks.AppPreferenceInterface

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		test_helper.SetContextUser(&ctx, 101, []string{"supplierpanel:editverifiedblockedsupplieronly:admin"})

		mocks.SetAuditLogMock()
		mockAudit = mocks.SetAuditLogMock()
		mockAudit.On("RecordAuditAction", ctx, mock.Anything).Return(nil)

		appPreferenceMockInstance = new(aaaMocks.AppPreferenceInterface)
		aaaModels.InjectMockAppPreferenceServiceInstance(appPreferenceMockInstance)
		appPreferenceMockInstance.On("GetValue", ctx, "allowed_supplier_types", []string{"L0", "L1", "L2", "L3", "Hlc", "Captive", "Driver"}).Return([]string{"L1", "Hlc"})
		appPreferenceMockInstance.On("GetValue", ctx, "supplier_update_allowed_permission", mock.Anything).Return("supplierpanel:editverifiedblockedsupplieronly:admin")
		appPreferenceMockInstance.On("GetValue", ctx, "should_send_supplier_log", "true").Return("true")
	})

	AfterEach(func() {
		mocks.UnsetAuditLogMock()
		aaaModels.InjectMockAppPreferenceServiceInstance(nil)
	})

	Context("Editing existing Supplier", func() {
		It("Should update supplier and return success response", func() {
			isPhoneVerified := true
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierType:    utils.Hlc,
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
				Id:                        supplier.ID,
				Name:                      "Name",
				Email:                     "Email",
				SupplierType:              uint64(utils.L1),
				BusinessName:              "BusinessName",
				Phone:                     "8801234567890",
				AlternatePhone:            "8801234567891",
				ShopImageUrl:              "ss2/shop_images/test.png",
				NidNumber:                 "12345",
				NidFrontImageUrl:          "ss2/shop_images/test.png",
				NidBackImageUrl:           "ss2/shop_images/test.png",
				TradeLicenseUrl:           "ss2/shop_images/test.pdf",
				AgreementUrl:              "ss2/shop_images/test.pdf",
				ShopOwnerImageUrl:         "ss2/shop_images/test.png",
				GuarantorImageUrl:         "ss2/shop_images/test.png",
				GuarantorNidNumber:        "12345",
				GuarantorNidFrontImageUrl: "ss2/shop_images/test.png",
				ChequeImageUrl:            "ss2/shop_images/test.png",
			}

			t := &testing.T{}
			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Preload("SupplierCategoryMappings").Preload("SupplierOpcMappings").
				First(&updatedSupplier, supplier.ID)

			Expect(updatedSupplier.Email).To(Equal(param.Email))
			Expect(updatedSupplier.Name).To(Equal(param.Name))
			Expect(updatedSupplier.BusinessName).To(Equal(param.BusinessName))
			Expect(updatedSupplier.Phone).To(Equal(param.Phone))
			Expect(updatedSupplier.AlternatePhone).To(Equal(param.AlternatePhone))
			Expect(updatedSupplier.ShopImageURL).To(Equal(param.ShopImageUrl))
			Expect(updatedSupplier.NidNumber).To(Equal(param.NidNumber))
			Expect(updatedSupplier.NidFrontImageUrl).To(Equal(param.NidFrontImageUrl))
			Expect(updatedSupplier.NidBackImageUrl).To(Equal(param.NidBackImageUrl))
			Expect(updatedSupplier.ShopOwnerImageUrl).To(Equal(param.ShopOwnerImageUrl))
			Expect(updatedSupplier.GuarantorImageUrl).To(Equal(param.GuarantorImageUrl))
			Expect(updatedSupplier.GuarantorNidNumber).To(Equal(param.GuarantorNidNumber))
			Expect(updatedSupplier.GuarantorNidFrontImageUrl).To(Equal(param.GuarantorNidFrontImageUrl))
			Expect(updatedSupplier.ChequeImageUrl).To(Equal(param.ChequeImageUrl))
			Expect(*updatedSupplier.IsPhoneVerified).To(Equal(false))
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusPending))

			Expect(len(updatedSupplier.SupplierCategoryMappings)).To(Equal(3))
			Expect(len(updatedSupplier.SupplierOpcMappings)).To(Equal(2))
			Expect(updatedSupplier.SupplierCategoryMappings[1].CategoryID).To(Equal(uint64(2)))

			partnerService := models.PartnerServiceMapping{}
			database.DBAPM(ctx).Model(&models.PartnerServiceMapping{}).Where("supplier_id = ?", supplier.ID).First(&partnerService)
			Expect(partnerService.ServiceLevel).To(Equal(utils.L1))
			Expect(partnerService.TradeLicenseUrl).To(Equal(param.TradeLicenseUrl))
			Expect(partnerService.AgreementUrl).To(Equal(param.AgreementUrl))

			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(1))
			mockedEventBus.AssertExpectations(t)
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
				Id:           supplier.ID,
				Name:         "Name",
				SupplierType: uint64(utils.L1),
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Email).To(Equal(supplier.Email))
			Expect(updatedSupplier.Name).To(Equal(param.Name))
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusBlocked))
			Expect(*updatedSupplier.IsPhoneVerified).To(Equal(true))
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("Editing allowed for limited permission", func() {

		BeforeEach(func() {
			test_utils.SetPermission(&ctx, []string{})
			mockAudit.On("RecordAuditAction", ctx, mock.Anything).Return(nil)

			appPreferenceMockInstance.On("GetValue", ctx, "allowed_supplier_types", []string{"L0", "L1", "L2", "L3", "Hlc", "Captive", "Driver"}).Return([]string{"L1", "Hlc"})
			appPreferenceMockInstance.On("GetValue", ctx, "supplier_update_allowed_permission", mock.Anything).Return("supplierpanel:editverifiedblockedsupplieronly:admin")
		})

		AfterEach(func() {
			aaaModels.InjectMockAppPreferenceServiceInstance(nil)
		})

		It("Should return success on updating pending supplier", func() {
			isPhoneVerified := true
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				IsPhoneVerified: &isPhoneVerified,
				Status:          models.SupplierStatusPending,
			})
			param := &supplierpb.SupplierObject{
				Id:           supplier.ID,
				Name:         "Name",
				SupplierType: uint64(utils.L1),
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Email).To(Equal(supplier.Email))
			Expect(updatedSupplier.Name).To(Equal(param.Name))
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusPending))
			Expect(*updatedSupplier.IsPhoneVerified).To(Equal(true))
		})

		It("Should return error on updating verified supplier", func() {
			isPhoneVerified := true
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				IsPhoneVerified: &isPhoneVerified,
				Status:          models.SupplierStatusVerified,
			})
			param := &supplierpb.SupplierObject{
				Id:   supplier.ID,
				Name: "Name",
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Change Not Allowed"))
			mockedEventBus.AssertExpectations(t)
		})

		It("Should return error on updating blocked supplier", func() {
			isPhoneVerified := true
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				IsPhoneVerified: &isPhoneVerified,
				Status:          models.SupplierStatusBlocked,
			})
			param := &supplierpb.SupplierObject{
				Id:   supplier.ID,
				Name: "Name",
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Change Not Allowed"))
			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(0))
			mockedEventBus.AssertExpectations(t)
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
				Id:           supplier.ID,
				Name:         "Name",
				SupplierType: uint64(utils.Hlc),
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			updatedSupplier := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Email).To(Equal(supplier.Email))
			Expect(updatedSupplier.Name).To(Equal(param.Name))
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusPending))
			Expect(*updatedSupplier.IsPhoneVerified).To(Equal(true))
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("Editing invalid supplier", func() {
		It("Should return error response", func() {
			param := &supplierpb.SupplierObject{Id: 1000}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Not Found"))
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("Editing with new set of category ids", func() {
		It("Should delete old mapping and add new mapping", func() {
			category_ids := []uint64{100, 101, 102}
			mockCategory := mocks.SetCategoryMock()
			mockCategory.On("GetCategoriesData", ctx, category_ids).Return(&categoryPb.CategoryDataList{Data: []*categoryPb.CategoryData{
				{Id: 100},
				{Id: 101},
				{Id: 102},
			}}, nil)
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
				CategoryIds:  []uint64{100, 101, 102},
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

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
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("Editing with category ids from cmt response", func() {
		It("Should delete old mapping and add new mapping", func() {
			category_ids := []uint64{100, 101, 102}
			mockCategory := mocks.SetCategoryMock()
			mockCategory.On("GetCategoriesData", ctx, category_ids).Return(&categoryPb.CategoryDataList{Data: []*categoryPb.CategoryData{
				{Id: 100},
				{Id: 101},
			}}, nil)
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
				CategoryIds:  []uint64{100, 101, 102},
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).Edit(ctx, param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier Edited Successfully"))

			var categoryIds []uint64
			database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Pluck("category_id", &categoryIds)
			Expect(len(categoryIds)).To(Equal(2))
			Expect(categoryIds).To(ContainElements([]uint64{100, 101}))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierCategoryMapping{}).Unscoped().Where("supplier_category_mappings.supplier_id = ? and supplier_category_mappings.deleted_at IS NULL", supplier.ID).Count(&count)
			Expect(count).To(Equal(2))
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("Editing with new set of category ids which got removed before", func() {
		It("Should restore deleted mapping", func() {
			category_ids := []uint64{101, 200, 567}
			mockCategory := mocks.SetCategoryMock()
			mockCategory.On("GetCategoriesData", ctx, category_ids).Return(&categoryPb.CategoryDataList{Data: []*categoryPb.CategoryData{
				{Id: 101},
				{Id: 200},
				{Id: 567},
			}}, nil)
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

			test := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(test, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

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
			mockedEventBus.AssertExpectations(test)
		})
	})

	Context("Editing with invalid phone number", func() {
		It("Should return error response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := &supplierpb.SupplierObject{
				Id:    supplier1.ID,
				Phone: "1234567890",
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while updating Supplier: Phone Number should have 13 digits"))
			mockedEventBus.AssertExpectations(t)
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

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while updating Supplier: Phone Number Already Exists"))
			mockedEventBus.AssertExpectations(t)
		})
	})
})
