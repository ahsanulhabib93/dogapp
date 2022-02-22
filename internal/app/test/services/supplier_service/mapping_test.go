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
)

var _ = Describe("MapSupplier", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Supplier-OPC map add", func() {
		It("Should Respond with success while mapping with OPC", func() {
			opcId := 101
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			resp, err := new(services.SupplierService).SupplierMap(ctx, &supplierpb.SupplierMappingParams{
				SupplierId:    supplier.ID,
				Id:            uint64(opcId),
				MapWith:       "OPC",
				OperationType: "Add",
			})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(true))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierOpcMapping{}).Where("processing_center_id = ? AND supplier_id = ?", opcId, supplier.ID).Count(&count)
			Expect(count).To(Equal(1))
		})
		It("Should Respond with success while mapping with existing OPC", func() {
			opcId := 101
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
				OperationType: "Add",
			})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(true))

			var count int
			database.DBAPM(ctx).Model(&models.SupplierOpcMapping{}).Where("processing_center_id = ? AND supplier_id = ?", opcId, supplier.ID).Count(&count)
			Expect(count).To(Equal(1))
		})
		It("Should Respond with success while mapping with delete OPC mapping", func() {
			opcId := 101
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
				OperationType: "Add",
			})
			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(true))

			database.DBAPM(ctx).Model(&models.SupplierOpcMapping{}).Where("processing_center_id = ? AND supplier_id = ?", opcId, supplier.ID).Count(&count)
			Expect(count).To(Equal(1))
		})
		It("Should Respond with error for invalid supplier ID", func() {
			opcId := 101
			resp, err := new(services.SupplierService).SupplierMap(ctx, &supplierpb.SupplierMappingParams{
				SupplierId:    123,
				Id:            uint64(opcId),
				MapWith:       "OPC",
				OperationType: "Add",
			})

			Expect(err).To(BeNil())
			Expect(resp.Success).To(Equal(false))
			Expect(resp.Message).To(Equal("Supplier Not Found"))
		})
	})

	Context("Supplier-OPC map delete", func() {
		It("Should Respond with success while deleting OPC mapping", func() {
			opcId := 101
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

			var count int
			database.DBAPM(ctx).Model(&models.SupplierOpcMapping{}).Where("processing_center_id = ? AND supplier_id = ?", opcId, supplier.ID).Count(&count)
			Expect(count).To(Equal(0))
		})
		It("Should Respond with success while deleting non-existing OPC mapping", func() {
			opcId := 101
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
			opcId := 101
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
