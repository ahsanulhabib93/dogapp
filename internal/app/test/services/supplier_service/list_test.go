package supplier_service_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	opcPb "github.com/voonik/goConnect/api/go/oms/processing_center"
	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/misc"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
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

	Context("Supplier List", func() {
		It("Should Respond with all the suppliers", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 1},
					{CategoryID: 2},
				},
				SupplierType: utils.Hlc,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 3},
					{ProcessingCenterID: 4},
				},
			})

			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})

			res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{})
			Expect(err).To(BeNil())
			Expect(res.TotalCount).To(Equal(uint64(2)))
			Expect(len(res.Data)).To(Equal(2))
			supplierData1 := res.Data[0]
			Expect(supplierData1.Email).To(Equal(supplier1.Email))
			Expect(supplierData1.Name).To(Equal(supplier1.Name))
			Expect(supplierData1.CategoryIds).To(Equal([]uint64{1, 2}))
			Expect(supplierData1.OpcIds).To(Equal([]uint64{3, 4}))
			Expect(supplierData1.SupplierType).To(Equal(uint64(utils.Hlc)))
			Expect(supplierData1.Status).To(Equal(models.SupplierStatusPending))

			supplierData2 := res.Data[1]
			Expect(supplierData2.Email).To(Equal(supplier2.Email))
			Expect(supplierData2.Name).To(Equal(supplier2.Name))
			Expect(supplierData2.CategoryIds).To(Equal([]uint64{}))
			Expect(supplierData2.OpcIds).To(Equal([]uint64{}))
			Expect(supplierData2.SupplierType).To(Equal(uint64(utils.L1)))
			Expect(supplierData2.Status).To(Equal(models.SupplierStatusPending))
		})

		It("Should Respond with all the suppliers and non-deleted opc/category ids", func() {
			deletedAt := time.Now()
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierCategoryMappings: []models.SupplierCategoryMapping{
					{CategoryID: 1},
					{CategoryID: 2, DeletedAt: &deletedAt},
					{CategoryID: 3},
				},
				SupplierType: utils.Hlc,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 3},
					{ProcessingCenterID: 4, DeletedAt: &deletedAt},
				},
			})

			test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})

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
			Expect(supplierData2.CategoryIds).To(Equal([]uint64{}))
			Expect(supplierData2.OpcIds).To(Equal([]uint64{}))
		})
	})

	It("Should Respond with success with OPC filter (multiple)", func() {
		deletedAt := time.Now()
		suppliers := []*models.Supplier{
			test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierType: utils.Hlc,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 1},
					{ProcessingCenterID: 2},
				},
			}),
			test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierType: utils.Hlc,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 2},
					{ProcessingCenterID: 3, DeletedAt: &deletedAt},
				},
			}),
			test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierType:        utils.Hlc,
				SupplierOpcMappings: []models.SupplierOpcMapping{{ProcessingCenterID: 3}},
			}),
			test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierType:        utils.Hlc,
				SupplierOpcMappings: []models.SupplierOpcMapping{{ProcessingCenterID: 4}},
			}),
		}

		test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})
		res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{OpcIds: []uint64{1, 3}})
		Expect(err).To(BeNil())
		Expect(res.TotalCount).To(Equal(uint64(2)))
		Expect(len(res.Data)).To(Equal(2))

		Expect(res.Data[0].Id).To(Equal(suppliers[0].ID))
		Expect(res.Data[0].OpcIds).To(Equal([]uint64{1}))
		Expect(res.Data[1].Id).To(Equal(suppliers[2].ID))
		Expect(res.Data[1].OpcIds).To(Equal([]uint64{3}))
	})

	It("Should Respond with success with OPC filter (single)", func() {
		deletedAt := time.Now()
		suppliers := []*models.Supplier{
			test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierType: utils.Hlc,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 1},
					{ProcessingCenterID: 2},
				},
			}),
			test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierType: utils.Hlc,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 2},
					{ProcessingCenterID: 3, DeletedAt: &deletedAt},
				},
			}),
			test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierType:        utils.Hlc,
				SupplierOpcMappings: []models.SupplierOpcMapping{{ProcessingCenterID: 3}},
			}),
			test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierType:        utils.Hlc,
				SupplierOpcMappings: []models.SupplierOpcMapping{{ProcessingCenterID: 4}},
			}),
		}

		test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})
		res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{OpcId: uint64(1)})
		Expect(err).To(BeNil())
		Expect(res.TotalCount).To(Equal(uint64(1)))
		Expect(len(res.Data)).To(Equal(1))

		Expect(res.Data[0].Id).To(Equal(suppliers[0].ID))
		Expect(res.Data[0].OpcIds).To(Equal([]uint64{1}))
	})

	It("Should Respond with all active suppliers", func() {
		supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{
			SupplierCategoryMappings: []models.SupplierCategoryMapping{
				{CategoryID: 1},
				{CategoryID: 2},
			},
			SupplierType: utils.Hlc,
			SupplierOpcMappings: []models.SupplierOpcMapping{
				{ProcessingCenterID: 3},
				{ProcessingCenterID: 4},
			},
			Status: models.SupplierStatusActive,
		})

		test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1})
		res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{Status: models.SupplierStatusActive})
		Expect(err).To(BeNil())
		Expect(res.TotalCount).To(Equal(uint64(1)))
		Expect(len(res.Data)).To(Equal(1))
		supplierData1 := res.Data[0]
		Expect(supplierData1.Email).To(Equal(supplier1.Email))
		Expect(supplierData1.Name).To(Equal(supplier1.Name))
		Expect(supplierData1.CategoryIds).To(Equal([]uint64{1, 2}))
		Expect(supplierData1.OpcIds).To(Equal([]uint64{3, 4}))
		Expect(supplierData1.SupplierType).To(Equal(uint64(utils.Hlc)))
		Expect(supplierData1.Status).To(Equal(models.SupplierStatusActive))
	})

	It("Should Respond with success for pagination", func() {
		suppliers := []*models.Supplier{
			test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1}),
			test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1}),
			test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1}),
			test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.L1}),
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

	It("Should Respond with success for pagination with filter", func() {
		suppliers := []*models.Supplier{
			test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusActive}),
			test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusActive}),
			test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusPending}),
			test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusActive}),
		}

		res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{
			Page:    1,
			PerPage: 2,
			Status:  models.SupplierStatusActive,
		})
		Expect(err).To(BeNil())
		Expect(len(res.Data)).To(Equal(2))
		Expect(res.TotalCount).To(Equal(uint64(3)))
		Expect(res.Data[0].Id).To(Equal(suppliers[0].ID))
		Expect(res.Data[1].Id).To(Equal(suppliers[1].ID))
	})

	It("Should Respond with SA user related supplier", func() {
		mockOpc := mocks.SetOpcMock()
		mockOpc.On("ProcessingCenterList", ctx, userId).Return(&opcPb.ProcessingCenterListResponse{
			Data: []*opcPb.OpcDetail{{OpcId: 1}, {OpcId: 2}, {OpcId: 3}}}, nil)

		deletedAt := time.Now()
		suppliers := []*models.Supplier{
			test_helper.CreateSupplier(ctx, &models.Supplier{
				Status: models.SupplierStatusActive,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 1},
					{ProcessingCenterID: 2},
				},
			}),
			test_helper.CreateSupplier(ctx, &models.Supplier{
				Status: models.SupplierStatusActive,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 3, DeletedAt: &deletedAt},
				},
			}),
			test_helper.CreateSupplier(ctx, &models.Supplier{
				Status: models.SupplierStatusActive,
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{ProcessingCenterID: 3},
				},
			}),
			test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusActive}),
		}

		res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{
			AssociatedWithCurrentUser: true,
		})

		Expect(err).To(BeNil())
		Expect(len(res.Data)).To(Equal(2))
		Expect(res.TotalCount).To(Equal(uint64(2)))
		Expect(res.Data[0].Id).To(Equal(suppliers[0].ID))
		Expect(res.Data[0].OpcIds).To(Equal([]uint64{1, 2}))
		Expect(res.Data[1].Id).To(Equal(suppliers[2].ID))
		Expect(res.Data[1].OpcIds).To(Equal([]uint64{3}))
	})

	It("Should Respond with no supplier for user with no opc mapped", func() {
		mockOpc := mocks.SetOpcMock()
		mockOpc.On("ProcessingCenterList", ctx, userId).Return(&opcPb.ProcessingCenterListResponse{
			Data: []*opcPb.OpcDetail{{OpcId: 1}, {OpcId: 2}, {OpcId: 3}}}, nil)

		test_helper.CreateSupplier(ctx, &models.Supplier{
			Status: models.SupplierStatusActive,
			SupplierOpcMappings: []models.SupplierOpcMapping{
				{ProcessingCenterID: 11},
				{ProcessingCenterID: 12},
			},
		})
		test_helper.CreateSupplier(ctx, &models.Supplier{
			Status: models.SupplierStatusActive,
			SupplierOpcMappings: []models.SupplierOpcMapping{
				{ProcessingCenterID: 13},
			},
		})
		test_helper.CreateSupplier(ctx, &models.Supplier{
			Status: models.SupplierStatusActive,
			SupplierOpcMappings: []models.SupplierOpcMapping{
				{ProcessingCenterID: 13},
			},
		})
		test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusActive})

		res, err := new(services.SupplierService).List(ctx, &supplierpb.ListParams{
			AssociatedWithCurrentUser: true,
		})

		Expect(err).To(BeNil())
		Expect(len(res.Data)).To(Equal(0))
		Expect(res.TotalCount).To(Equal(uint64(0)))
	})
})
