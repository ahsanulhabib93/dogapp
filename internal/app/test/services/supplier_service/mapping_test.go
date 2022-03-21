package supplier_service_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	opcPb "github.com/voonik/goConnect/api/go/oms/processing_center"
	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("MapSupplier", func() {
	var ctx context.Context

	BeforeEach(func() {
		mocks.UnsetOpcMock()
		test_utils.GetContext(&ctx)
	})

	Context("Supplier-OPC map add", func() {
		It("Should Respond with success while mapping with OPC", func() {
			opcId := uint64(101)
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, []uint64{uint64(opcId)}).Return(&opcPb.ProcessingCenterListResponse{Data: []*opcPb.OpcDetail{
				{OpcId: opcId},
			}}, nil)
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			resp, err := new(services.SupplierService).SupplierMap(ctx, &supplierpb.SupplierMappingParams{
				SupplierId:    supplier.ID,
				Id:            opcId,
				MapWith:       "OPC",
				OperationType: "Add",
			})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(true))
			Expect(resp.Message).To(Equal("Supplier Mapped with OPC"))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierOpcMapping{}).Where("processing_center_id = ? AND supplier_id = ?", opcId, supplier.ID).Count(&count)
			Expect(count).To(Equal(1))
		})
		It("Should Respond with success while mapping with existing OPC", func() {
			opcId := uint64(101)
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, []uint64{uint64(opcId)}).Return(&opcPb.ProcessingCenterListResponse{Data: []*opcPb.OpcDetail{
				{OpcId: opcId},
			}}, nil)
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{
						ProcessingCenterID: opcId,
					},
				},
			})
			resp, err := new(services.SupplierService).SupplierMap(ctx, &supplierpb.SupplierMappingParams{
				SupplierId:    supplier.ID,
				Id:            opcId,
				MapWith:       "OPC",
				OperationType: "Add",
			})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(true))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierOpcMapping{}).Where("processing_center_id = ? AND supplier_id = ?", opcId, supplier.ID).Count(&count)
			Expect(count).To(Equal(1))
		})
		It("Should Respond with success while mapping with delete OPC mapping", func() {
			opcId := uint64(101)
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, []uint64{uint64(opcId)}).Return(&opcPb.ProcessingCenterListResponse{Data: []*opcPb.OpcDetail{
				{OpcId: opcId},
			}}, nil)

			deletedAt := time.Now()
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{
						DeletedAt:          &deletedAt,
						ProcessingCenterID: opcId,
					},
				},
			})
			var count int
			database.DBAPM(ctx).Unscoped().Model(&models.SupplierOpcMapping{}).Where("processing_center_id = ? AND supplier_id = ?", opcId, supplier.ID).Count(&count)
			Expect(count).To(Equal(1))

			resp, err := new(services.SupplierService).SupplierMap(ctx, &supplierpb.SupplierMappingParams{
				SupplierId:    supplier.ID,
				Id:            opcId,
				MapWith:       "OPC",
				OperationType: "Add",
			})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(true))

			database.DBAPM(ctx).Model(&models.SupplierOpcMapping{}).Where("processing_center_id = ? AND supplier_id = ?", opcId, supplier.ID).Count(&count)
			Expect(count).To(Equal(1))
		})

		It("Should Respond with error for invalid supplier ID", func() {
			opcId := uint64(101)
			resp, err := new(services.SupplierService).SupplierMap(ctx, &supplierpb.SupplierMappingParams{
				SupplierId:    123,
				Id:            opcId,
				MapWith:       "OPC",
				OperationType: "Add",
			})

			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(false))
			Expect(resp.Message).To(Equal("Supplier Not Found"))
		})

		It("Should Respond with error for invalid OPC ID", func() {
			opcId := uint64(101)
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, []uint64{uint64(opcId)}).Return(&opcPb.ProcessingCenterListResponse{Data: []*opcPb.OpcDetail{
				{OpcId: 123},
			}}, nil)
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			resp, err := new(services.SupplierService).SupplierMap(ctx, &supplierpb.SupplierMappingParams{
				SupplierId:    supplier.ID,
				Id:            opcId,
				MapWith:       "OPC",
				OperationType: "Add",
			})

			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(false))
			Expect(resp.Message).To(Equal("invalid opc id #(101)"))
		})

		It("Should Respond with error failed oms remote call", func() {
			opcId := uint64(101)
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, []uint64{uint64(opcId)}).Return(&opcPb.ProcessingCenterListResponse{}, errors.New("failed"))
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			resp, err := new(services.SupplierService).SupplierMap(ctx, &supplierpb.SupplierMappingParams{
				SupplierId:    supplier.ID,
				Id:            opcId,
				MapWith:       "OPC",
				OperationType: "Add",
			})

			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(false))
			Expect(resp.Message).To(Equal("failed to fetch opc list"))
		})
	})

	Context("Supplier-OPC map delete", func() {
		It("Should Respond with success while deleting OPC mapping", func() {
			opcId := uint64(101)
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, []uint64{uint64(opcId)}).Return(&opcPb.ProcessingCenterListResponse{Data: []*opcPb.OpcDetail{}}, nil)

			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{
						ProcessingCenterID: uint64(opcId),
					},
				},
			})
			resp, err := new(services.SupplierService).SupplierMap(ctx, &supplierpb.SupplierMappingParams{
				SupplierId:    supplier.ID,
				Id:            uint64(opcId),
				MapWith:       "OPC",
				OperationType: "Delete",
			})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(true))
			Expect(resp.Message).To(Equal("Supplier Unmapped with OPC"))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierOpcMapping{}).Where("processing_center_id = ? AND supplier_id = ?", opcId, supplier.ID).Count(&count)
			Expect(count).To(Equal(0))
		})
		It("Should Respond with success while deleting non-existing OPC mapping", func() {
			opcId := uint64(101)
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, []uint64{uint64(opcId)}).Return(&opcPb.ProcessingCenterListResponse{Data: []*opcPb.OpcDetail{}}, nil)

			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			resp, err := new(services.SupplierService).SupplierMap(ctx, &supplierpb.SupplierMappingParams{
				SupplierId:    supplier.ID,
				Id:            uint64(opcId),
				MapWith:       "OPC",
				OperationType: "Delete",
			})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(true))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierOpcMapping{}).Where("processing_center_id = ? AND supplier_id = ?", opcId, supplier.ID).Count(&count)
			Expect(count).To(Equal(0))
		})
		It("Should Respond with success while deleting already deleted OPC mapping", func() {
			opcId := uint64(101)
			mockOpc := mocks.SetOpcMock()
			mockOpc.On("GetProcessingCenterListWithOpcIds", ctx, []uint64{uint64(opcId)}).Return(&opcPb.ProcessingCenterListResponse{Data: []*opcPb.OpcDetail{}}, nil)

			deletedAt := time.Now()
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
				SupplierOpcMappings: []models.SupplierOpcMapping{
					{
						DeletedAt:          &deletedAt,
						ProcessingCenterID: uint64(opcId),
					},
				},
			})
			var count int
			database.DBAPM(ctx).Unscoped().Model(&models.SupplierOpcMapping{}).Where("processing_center_id = ? AND supplier_id = ?", opcId, supplier.ID).Count(&count)
			Expect(count).To(Equal(1))

			resp, err := new(services.SupplierService).SupplierMap(ctx, &supplierpb.SupplierMappingParams{
				SupplierId:    supplier.ID,
				Id:            uint64(opcId),
				MapWith:       "OPC",
				OperationType: "Delete",
			})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(true))

			database.DBAPM(ctx).Unscoped().Model(&models.SupplierOpcMapping{}).Where("processing_center_id = ? AND supplier_id = ?", opcId, supplier.ID).Count(&count)
			Expect(count).To(Equal(1))
		})
	})
})
