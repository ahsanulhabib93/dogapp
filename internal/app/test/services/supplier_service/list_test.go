package supplier_service_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	opcPb "github.com/voonik/goConnect/api/go/oms/processing_center"
	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("ListSupplier", func() {
	var ctx context.Context
	var userId uint64 = uint64(101)

	BeforeEach(func() {
		mocks.UnsetOpcMock()
		test_utils.GetContext(&ctx)
		test_helper.SetContextUser(&ctx, userId, []string{"supplierpanel:allservices:view"})
	})

	Context("Supplier List", func() {
		It("Should Respond with all the suppliers and attachments", func() {
			hlcServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Supplier, "Hlc")
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{{CategoryID: 1}, {CategoryID: 2}},
				SupplierOpcMappings:      []models.SupplierOpcMapping{{ProcessingCenterID: 3}, {ProcessingCenterID: 4}},
				PartnerServiceMappings:   []models.PartnerServiceMapping{{PartnerServiceLevelID: hlcServiceLevel.ID}},
			})
			driverServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Transporter, "Driver")
			test_helper.CreatePartnerServiceMapping(ctx, &models.PartnerServiceMapping{
				SupplierId:            supplier1.ID,
				ServiceType:           utils.Transporter,
				PartnerServiceLevelID: driverServiceLevel.ID,
			})

			isPhoneVerified := true
			l1ServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Supplier, "L1")
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{
				IsPhoneVerified:        &isPhoneVerified,
				PartnerServiceMappings: []models.PartnerServiceMapping{{PartnerServiceLevelID: l1ServiceLevel.ID}},
			})

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(2)))
			Expect(len(res.Data)).To(Equal(2))

			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier1.Email))
			Expect(supplierData1.Name).To(Equal(supplier1.Name))
			Expect(supplierData1.Phone).To(Equal(supplier1.Phone))
			Expect(supplierData1.AlternatePhone).To(Equal(supplier1.AlternatePhone))
			Expect(supplierData1.BusinessName).To(Equal(supplier1.BusinessName))
			Expect(supplierData1.ShopImageUrl).To(Equal(supplier1.ShopImageURL))
			Expect(supplierData1.Reason).To(Equal(supplier1.Reason))
			Expect(supplierData1.IsPhoneVerified).To(Equal(false))
			Expect(supplierData1.CategoryIds).To(Equal([]uint64{1, 2}))
			Expect(supplierData1.OpcIds).To(Equal([]uint64{3, 4}))
			Expect(supplierData1.Status).To(Equal(string(models.SupplierStatusPending)))

			Expect(supplierData1.PartnerServices).To(HaveLen(2))
			Expect(supplierData1.PartnerServices[0].ServiceType).To(Equal("Supplier"))
			Expect(supplierData1.PartnerServices[0].ServiceLevel).To(Equal("Hlc"))
			Expect(supplierData1.PartnerServices[1].ServiceType).To(Equal("Transporter"))
			Expect(supplierData1.PartnerServices[1].ServiceLevel).To(Equal("Driver"))

			supplierData2 := res.Data[1]
			Expect(supplierData2.Email).To(Equal(supplier2.Email))
			Expect(supplierData2.Name).To(Equal(supplier2.Name))
			Expect(supplierData2.Phone).To(Equal(supplier2.Phone))
			Expect(supplierData2.AlternatePhone).To(Equal(supplier2.AlternatePhone))
			Expect(supplierData2.BusinessName).To(Equal(supplier2.BusinessName))
			Expect(supplierData2.ShopImageUrl).To(Equal(supplier2.ShopImageURL))
			Expect(supplierData2.Reason).To(Equal(supplier2.Reason))
			Expect(supplierData2.IsPhoneVerified).To(Equal(true))
			Expect(supplierData2.CategoryIds).To(Equal([]uint64{}))
			Expect(supplierData2.OpcIds).To(Equal([]uint64{}))
			Expect(supplierData2.Status).To(Equal(string(models.SupplierStatusPending)))

			Expect(supplierData2.PartnerServices).To(HaveLen(1))
			Expect(supplierData2.PartnerServices[0].ServiceType).To(Equal("Supplier"))
			Expect(supplierData2.PartnerServices[0].ServiceLevel).To(Equal("L1"))
		})
	})

	Context("When deleted OPC and category mapping is present", func() {
		It("Should Respond with all the suppliers with non-deleted opc/category ids", func() {
			deletedAt := time.Now()
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 1},
					{CategoryID: 2, DeletedAt: &deletedAt},
					{CategoryID: 3},
				},
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 3},
					{ProcessingCenterID: 4, DeletedAt: &deletedAt},
				},
			})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{})

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(2)))
			Expect(len(res.Data)).To(Equal(2))
			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier1.Email))
			Expect(supplierData1.Name).To(Equal(supplier1.Name))
			Expect(supplierData1.CategoryIds).To(Equal([]uint64{1, 3}))
			Expect(supplierData1.OpcIds).To(Equal([]uint64{3}))

			supplierData2 := res.Data[1]
			Expect(supplierData2.Email).To(Equal(supplier2.Email))
			Expect(supplierData2.CategoryIds).To(Equal([]uint64{}))
			Expect(supplierData2.OpcIds).To(Equal([]uint64{}))
		})
	})

	Context("When all OPC and category mapping are deleted", func() {
		It("Should Respond with all the suppliers with empty opc/category ids", func() {
			deletedAt := time.Now()
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{{CategoryID: 1, DeletedAt: &deletedAt}},
				SupplierOpcMappings:      []models.SupplierOpcMapping{{ProcessingCenterID: 4, DeletedAt: &deletedAt}},
			})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{})

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(2)))
			Expect(len(res.Data)).To(Equal(2))

			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier1.Email))
			Expect(supplierData1.Name).To(Equal(supplier1.Name))
			Expect(supplierData1.CategoryIds).To(Equal([]uint64{}))
			Expect(supplierData1.OpcIds).To(Equal([]uint64{}))

			supplierData2 := res.Data[1]
			Expect(supplierData2.Email).To(Equal(supplier2.Email))
			Expect(supplierData2.CategoryIds).To(Equal([]uint64{}))
			Expect(supplierData2.OpcIds).To(Equal([]uint64{}))
		})
	})

	Context("When OPC filter is applied with multiple OPC ids", func() {
		It("Should Respond with corresponding suppliers", func() {
			deletedAt := time.Now()
			suppliers := []*models.Supplier{
				test_helper.CreateSupplier(ctx, &models.Supplier{
					SupplierOpcMappings: []models.SupplierOpcMapping{
						{ProcessingCenterID: 1},
						{ProcessingCenterID: 2},
					},
				}),
				test_helper.CreateSupplier(ctx, &models.Supplier{
					SupplierOpcMappings: []models.SupplierOpcMapping{
						{ProcessingCenterID: 2},
						{ProcessingCenterID: 3, DeletedAt: &deletedAt},
					},
				}),
				test_helper.CreateSupplier(ctx, &models.Supplier{
					SupplierOpcMappings: []models.SupplierOpcMapping{{ProcessingCenterID: 3}},
				}),
				test_helper.CreateSupplier(ctx, &models.Supplier{
					SupplierOpcMappings: []models.SupplierOpcMapping{{ProcessingCenterID: 4}},
				}),
			}

			test_helper.CreateSupplier(ctx, &models.Supplier{})
			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{OpcIds: []uint64{1, 3}})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(2)))
			Expect(len(res.Data)).To(Equal(2))

			Expect(res.Data[0].Id).To(Equal(suppliers[0].ID))
			Expect(res.Data[0].OpcIds).To(Equal([]uint64{1}))
			Expect(res.Data[1].Id).To(Equal(suppliers[2].ID))
			Expect(res.Data[1].OpcIds).To(Equal([]uint64{3}))
		})
	})

	Context("When OPC filter is applied with single OPC id", func() {
		It("Should Respond with corresponding suppliers", func() {
			deletedAt := time.Now()
			suppliers := []*models.Supplier{
				test_helper.CreateSupplier(ctx, &models.Supplier{
					SupplierOpcMappings: []models.SupplierOpcMapping{
						{ProcessingCenterID: 1},
						{ProcessingCenterID: 2},
					},
				}),
				test_helper.CreateSupplier(ctx, &models.Supplier{
					SupplierOpcMappings: []models.SupplierOpcMapping{
						{ProcessingCenterID: 2},
						{ProcessingCenterID: 3, DeletedAt: &deletedAt},
					},
				}),
				test_helper.CreateSupplier(ctx, &models.Supplier{
					SupplierOpcMappings: []models.SupplierOpcMapping{{ProcessingCenterID: 3}},
				}),
				test_helper.CreateSupplier(ctx, &models.Supplier{
					SupplierOpcMappings: []models.SupplierOpcMapping{{ProcessingCenterID: 4}},
				}),
			}

			test_helper.CreateSupplier(ctx, &models.Supplier{})
			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{OpcId: uint64(1)})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(1)))
			Expect(len(res.Data)).To(Equal(1))

			Expect(res.Data[0].Id).To(Equal(suppliers[0].ID))
			Expect(res.Data[0].OpcIds).To(Equal([]uint64{1}))
		})
	})

	Context("When active status filter is applied", func() {
		It("Should Respond with all active suppliers", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 1},
					{CategoryID: 2},
				},
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 3},
					{ProcessingCenterID: 4},
				},
				Status: models.SupplierStatusVerified,
			})

			test_helper.CreateSupplier(ctx, &models.Supplier{})
			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{Status: string(models.SupplierStatusVerified)})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(1)))
			Expect(len(res.Data)).To(Equal(1))
			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier1.Email))
			Expect(supplierData1.Name).To(Equal(supplier1.Name))
			Expect(supplierData1.CategoryIds).To(Equal([]uint64{1, 2}))
			Expect(supplierData1.OpcIds).To(Equal([]uint64{3, 4}))
			Expect(supplierData1.Status).To(Equal(string(models.SupplierStatusVerified)))
		})
	})

	Context("With pagination", func() {
		It("Should Respond with corresponding suppliers", func() {
			suppliers := []*models.Supplier{
				test_helper.CreateSupplier(ctx, &models.Supplier{}),
				test_helper.CreateSupplier(ctx, &models.Supplier{}),
				test_helper.CreateSupplier(ctx, &models.Supplier{}),
				test_helper.CreateSupplier(ctx, &models.Supplier{}),
			}

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{
				Page:    2,
				PerPage: 2,
			})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(2))
			Expect(res.TotalCount).To(Equal(uint64(4)))
			Expect(res.Data[0].Id).To(Equal(suppliers[2].ID))
			Expect(res.Data[1].Id).To(Equal(suppliers[3].ID))
		})
	})

	Context("When status filter is applied with pagination", func() {
		It("Should Respond with corresponding suppliers", func() {
			suppliers := []*models.Supplier{
				test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusVerified}),
				test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusFailed}),
				test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusPending}),
				test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusVerified}),
			}
			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{
				Page:    1,
				PerPage: 2,
				Status:  string(models.SupplierStatusVerified) + "," + string(models.SupplierStatusFailed),
			})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(2))
			Expect(res.TotalCount).To(Equal(uint64(3)))
			Expect(res.Data[0].Id).To(Equal(suppliers[0].ID))
			Expect(res.Data[1].Id).To(Equal(suppliers[1].ID))
		})
	})

	Context("While fetching suppliers related with SA user", func() {
		It("Should Respond with corresponding suppliers", func() {
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithUserId", ctx, userId).Return(&opcPb.ProcessingCenterListResponse{
				Data: []*opcPb.OpcDetail{{OpcId: 1}, {OpcId: 2}, {OpcId: 3}}}, nil)

			deletedAt := time.Now()
			suppliers := []*models.Supplier{
				test_helper.CreateSupplier(ctx, &models.Supplier{
					Status: models.SupplierStatusVerified,
					SupplierOpcMappings: []models.SupplierOpcMapping{
						{ProcessingCenterID: 1},
						{ProcessingCenterID: 2},
					},
				}),
				test_helper.CreateSupplier(ctx, &models.Supplier{
					Status: models.SupplierStatusVerified,
					SupplierOpcMappings: []models.SupplierOpcMapping{
						{ProcessingCenterID: 3, DeletedAt: &deletedAt},
					},
				}),
				test_helper.CreateSupplier(ctx, &models.Supplier{
					Status: models.SupplierStatusVerified,
					SupplierOpcMappings: []models.SupplierOpcMapping{
						{ProcessingCenterID: 3},
					},
				}),
				test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusVerified}),
			}

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{AssociatedWithCurrentUser: true})

			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(2))
			Expect(res.TotalCount).To(Equal(uint64(2)))
			Expect(res.Data[0].Id).To(Equal(suppliers[0].ID))
			Expect(res.Data[0].OpcIds).To(Equal([]uint64{1, 2}))
			Expect(res.Data[1].Id).To(Equal(suppliers[2].ID))
			Expect(res.Data[1].OpcIds).To(Equal([]uint64{3}))
		})
	})

	Context("When no OPC is mapped with current SA user", func() {
		It("Should Respond with no suppliers", func() {
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithUserId", ctx, userId).Return(&opcPb.ProcessingCenterListResponse{
				Data: []*opcPb.OpcDetail{{OpcId: 1}, {OpcId: 2}, {OpcId: 3}}}, nil)

			test_helper.CreateSupplier(ctx, &models.Supplier{
				Status: models.SupplierStatusVerified,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 11},
					{ProcessingCenterID: 12},
				},
			})
			test_helper.CreateSupplier(ctx, &models.Supplier{
				Status: models.SupplierStatusVerified,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 13},
				},
			})
			test_helper.CreateSupplier(ctx, &models.Supplier{
				Status: models.SupplierStatusVerified,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 13},
				},
			})
			test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusVerified})

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{
				AssociatedWithCurrentUser: true,
			})

			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(0))
			Expect(res.TotalCount).To(Equal(uint64(0)))
		})
	})

	Context("When created_at filter is applied", func() {
		It("Should Respond with corresponding suppliers", func() {
			test_helper.CreateSupplier(ctx, &models.Supplier{VaccountGorm: database.VaccountGorm{VModel: database.VModel{CreatedAt: time.Date(2021, 01, 10, 10, 0, 0, 0, time.UTC)}}})
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{VaccountGorm: database.VaccountGorm{VModel: database.VModel{CreatedAt: time.Date(2021, 01, 12, 10, 0, 0, 0, time.UTC)}}})
			test_helper.CreateSupplier(ctx, &models.Supplier{VaccountGorm: database.VaccountGorm{VModel: database.VModel{CreatedAt: time.Date(2021, 01, 14, 10, 0, 0, 0, time.UTC)}}})

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{CreatedAtGte: "2021-01-12", CreatedAtLte: "2021-01-13"})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(1)))
			Expect(len(res.Data)).To(Equal(1))

			supplierData := res.Data[0]
			Expect(supplierData.Email).To(Equal(supplier.Email))
			Expect(supplierData.Name).To(Equal(supplier.Name))
		})
	})

	Context("When supplier name is given", func() {
		It("Should Respond with suppliers with same name, same phone number and same account number", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{Name: "4444"})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{Phone: "8801444444444"})
			supplier3 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier3.ID, AccountNumber: "44404444", IsDefault: true})

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{Name: "444"})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(3)))
			Expect(len(res.Data)).To(Equal(3))

			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier1.Email))
			Expect(supplierData1.Name).To(Equal(supplier1.Name))
			Expect(supplierData1.Phone).To(Equal(supplier1.Phone))
			Expect(supplierData1.AlternatePhone).To(Equal(supplier1.AlternatePhone))
			Expect(supplierData1.BusinessName).To(Equal(supplier1.BusinessName))
			Expect(supplierData1.ShopImageUrl).To(Equal(supplier1.ShopImageURL))
			Expect(supplierData1.Reason).To(Equal(supplier1.Reason))

			supplierData2 := res.Data[1]
			Expect(supplierData2.Email).To(Equal(supplier2.Email))
			Expect(supplierData2.Name).To(Equal(supplier2.Name))
			Expect(supplierData2.Phone).To(Equal(supplier2.Phone))
			Expect(supplierData2.AlternatePhone).To(Equal(supplier2.AlternatePhone))
			Expect(supplierData2.BusinessName).To(Equal(supplier2.BusinessName))
			Expect(supplierData2.ShopImageUrl).To(Equal(supplier2.ShopImageURL))
			Expect(supplierData2.Reason).To(Equal(supplier2.Reason))

			supplierData3 := res.Data[2]
			Expect(supplierData3.Email).To(Equal(supplier3.Email))
			Expect(supplierData3.Name).To(Equal(supplier3.Name))
			Expect(supplierData3.Phone).To(Equal(supplier3.Phone))
			Expect(supplierData3.AlternatePhone).To(Equal(supplier3.AlternatePhone))
			Expect(supplierData3.BusinessName).To(Equal(supplier3.BusinessName))
			Expect(supplierData3.ShopImageUrl).To(Equal(supplier3.ShopImageURL))
			Expect(supplierData3.Reason).To(Equal(supplier3.Reason))

		})
	})

	Context("When ServiceLevels filter is applied", func() {
		It("Should Respond with corresponding suppliers", func() {
			mapping := helpers.GetServiceLevelIdMapping(ctx)
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{PartnerServiceMappings: []models.PartnerServiceMapping{{PartnerServiceLevelID: mapping["Hlc"]}}})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{PartnerServiceMappings: []models.PartnerServiceMapping{{PartnerServiceLevelID: mapping["L0"]}}})
			test_helper.CreateSupplier(ctx, &models.Supplier{PartnerServiceMappings: []models.PartnerServiceMapping{{PartnerServiceLevelID: mapping["L1"]}}})

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{ServiceLevels: []string{"Hlc", "L0"}})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(2)))
			Expect(len(res.Data)).To(Equal(2))

			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier1.Email))
			Expect(supplierData1.PartnerServices[0].ServiceLevel).To(Equal("Hlc"))

			supplierData2 := res.Data[1]
			Expect(supplierData2.Email).To(Equal(supplier2.Email))
			Expect(supplierData2.PartnerServices[0].ServiceLevel).To(Equal("L0"))
		})
	})

	Context("When ServiceTypes filter is applied", func() {
		It("Should Respond with supplier service type data", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{PartnerServiceMappings: []models.PartnerServiceMapping{{ServiceType: utils.Supplier}}})
			test_helper.CreateSupplier(ctx, &models.Supplier{PartnerServiceMappings: []models.PartnerServiceMapping{{ServiceType: utils.Transporter}}})

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{ServiceTypes: []string{"Supplier"}})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(1)))
			Expect(len(res.Data)).To(Equal(1))

			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier1.Email))
			Expect(supplierData1.PartnerServices[0].ServiceType).To(Equal("Supplier"))
		})
	})

	Context("When User has only supplier permission", func() {
		BeforeEach(func() {
			test_helper.SetContextUser(&ctx, 1, []string{"supplierpanel:supplierservice:view"})
		})
		It("Should Respond with only supplier service type data", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{PartnerServiceMappings: []models.PartnerServiceMapping{{ServiceType: utils.Supplier}}})
			driverServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Transporter, "Driver")
			test_helper.CreatePartnerServiceMapping(ctx, &models.PartnerServiceMapping{
				SupplierId:            supplier1.ID,
				ServiceType:           utils.Transporter,
				PartnerServiceLevelID: driverServiceLevel.ID,
			})
			test_helper.CreateSupplier(ctx, &models.Supplier{PartnerServiceMappings: []models.PartnerServiceMapping{{ServiceType: utils.Transporter}}})
			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(1)))
			Expect(len(res.Data)).To(Equal(1))

			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier1.Email))
			Expect(supplierData1.PartnerServices).To(HaveLen(1))
			Expect(supplierData1.PartnerServices[0].ServiceType).To(Equal("Supplier"))
		})
	})

	Context("When User has only supplier permission and transporter service type filter is applied", func() {
		BeforeEach(func() {
			test_helper.SetContextUser(&ctx, 1, []string{"supplierpanel:supplierservice:view"})
		})
		It("Should Respond with no supplier data", func() {
			test_helper.CreateSupplier(ctx, &models.Supplier{PartnerServiceMappings: []models.PartnerServiceMapping{{ServiceType: utils.Supplier}}})
			test_helper.CreateSupplier(ctx, &models.Supplier{PartnerServiceMappings: []models.PartnerServiceMapping{{ServiceType: utils.Transporter}}})

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{ServiceTypes: []string{"Transporter"}})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(0)))
			Expect(len(res.Data)).To(Equal(0))
		})
	})

})
