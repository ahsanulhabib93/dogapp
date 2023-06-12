package supplier_service_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	eventBus "github.com/voonik/goConnect/api/go/event_bus/publisher"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	aaaMocks "github.com/voonik/goFramework/pkg/aaa/models/mocks"
	"github.com/voonik/ss2/internal/app/publisher"
	mockPublisher "github.com/voonik/ss2/internal/app/publisher/mocks"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("RemoveDocument", func() {
	var ctx context.Context
	var mockAudit *mocks.AuditLogMock
	var supplier *models.Supplier
	var appPreferenceMockInstance *aaaMocks.AppPreferenceInterface

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		ctx = test_helper.SetContextUser(ctx, 101, []string{"supplierpanel:editverifiedblockedsupplieronly:admin"})

		supplierData := models.Supplier{
			GuarantorNidFrontImageUrl: "abc/xyz.jpg",
			Status:                    models.SupplierStatusVerified,
			PartnerServiceMappings: []models.PartnerServiceMapping{{
				AgreementUrl: "abc/xyz.pdf",
			}},
		}
		supplier = test_helper.CreateSupplier(ctx, &supplierData)

		mockAudit = mocks.SetAuditLogMock()
		mockAudit.On("RecordAuditAction", ctx, mock.Anything).Return(nil)

		appPreferenceMockInstance = new(aaaMocks.AppPreferenceInterface)
		aaaModels.InjectMockAppPreferenceServiceInstance(appPreferenceMockInstance)
		appPreferenceMockInstance.On("GetValue", ctx, "should_send_supplier_log", "true").Return("true")
		permission := "supplierpanel:editverifiedblockedsupplieronly:admin"
		appPreferenceMockInstance.On("GetValue", mock.Anything, "supplier_update_allowed_permission", permission).Return(permission)
	})

	AfterEach(func() {
		aaaModels.InjectMockAppPreferenceServiceInstance(nil)
	})

	Context("Removing primary document", func() {
		It("Should remove document successfully", func() {
			param := &supplierpb.RemoveDocumentParam{
				Id:           supplier.ID,
				DocumentType: "agreement_url",
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).RemoveDocument(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier agreement_url Removed Successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusPending))

			partnerServices := []*models.PartnerServiceMapping{{}}
			database.DBAPM(ctx).Model(supplier).Association("PartnerServiceMappings").Find(&partnerServices)
			Expect(len(partnerServices)).To(Equal(1))
			Expect(partnerServices[0].AgreementUrl).To(Equal(""))
			Expect(partnerServices[0].Active).To(Equal(false))

			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(1))
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("Removing secondary document", func() {
		It("Should remove document successfully", func() {
			param := &supplierpb.RemoveDocumentParam{
				Id:           supplier.ID,
				DocumentType: "guarantor_nid_front_image_url",
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).RemoveDocument(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier guarantor_nid_front_image_url Removed Successfully"))

			updatedSupplier := models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).First(&updatedSupplier, supplier.ID)
			Expect(updatedSupplier.Status).To(Equal(models.SupplierStatusVerified))
			Expect(updatedSupplier.GuarantorNidFrontImageUrl).To(Equal(""))

			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(1))
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("Removing agreement_url for given partner service mapping", func() {
		It("Should remove document successfully", func() {
			partnerService := test_helper.CreatePartnerServiceMapping(ctx, &models.PartnerServiceMapping{
				ServiceType:  utils.Transporter,
				ServiceLevel: utils.Driver,
				SupplierId:   supplier.ID,
				AgreementUrl: "abc/xyz.pdf",
				Active:       true,
			})

			param := &supplierpb.RemoveDocumentParam{
				Id:               supplier.ID,
				DocumentType:     "agreement_url",
				PartnerServiceId: partnerService.ID,
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).RemoveDocument(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Supplier agreement_url Removed Successfully"))

			partnerServices := []*models.PartnerServiceMapping{{}}
			database.DBAPM(ctx).Model(supplier).Association("PartnerServiceMappings").Find(&partnerServices)
			Expect(len(partnerServices)).To(Equal(2))

			Expect(partnerServices[0].AgreementUrl).To(Equal("abc/xyz.pdf"))
			Expect(partnerServices[0].Active).To(Equal(true))
			Expect(partnerServices[1].AgreementUrl).To(Equal(""))
			Expect(partnerServices[1].Active).To(Equal(false))

			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(1))
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("for invalid document type", func() {
		It("Should return error", func() {
			param := &supplierpb.RemoveDocumentParam{
				Id:           supplier.ID,
				DocumentType: "agreement_url_abc",
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).RemoveDocument(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Invalid Document Type"))

			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(0))
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("for un-allowed permission", func() {
		It("Should return error", func() {
			test_utils.SetPermission(&ctx, []string{"per:missi:on"})

			param := &supplierpb.RemoveDocumentParam{
				Id:           supplier.ID,
				DocumentType: "agreement_url",
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

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
			Expect(partnerServices[0].Active).To(Equal(true))
			mockedEventBus.AssertExpectations(t)
		})
	})

	Context("for invalid partner service ID", func() {
		It("Should return error", func() {
			param := &supplierpb.RemoveDocumentParam{
				Id:               supplier.ID,
				DocumentType:     "agreement_url",
				PartnerServiceId: 100,
			}

			t := &testing.T{}

			mockedEventBus, resetEventBus := mockPublisher.SetupMockPublisherClient(t, &publisher.EventBusClient)
			defer resetEventBus()

			mockedEventBus.On("Publish", ctx, mock.Anything, mock.Anything, mock.Anything).Return(&eventBus.PublishResponse{Success: true}, nil)

			res, err := new(services.SupplierService).RemoveDocument(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("ParnerServiceMapping not found"))

			Expect(mockAudit.Count["RecordAuditAction"]).To(Equal(0))
			mockedEventBus.AssertExpectations(t)
		})
	})
})
