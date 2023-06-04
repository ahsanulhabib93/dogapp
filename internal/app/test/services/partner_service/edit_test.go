package partner_service_service_test

import (
	"context"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	psmpb "github.com/voonik/goConnect/api/go/ss2/partner_service_mapping"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("EditPartnerService", func() {
	var ctx context.Context
	var mockAudit *mocks.AuditLogMock
	var userId uint64 = uint64(101)

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		mocks.UnsetOpcMock()

		ctx = test_helper.SetContextUser(ctx, userId, []string{})
		mockAudit = mocks.SetAuditLogMock()
		mockAudit.On("RecordAuditAction", ctx, mock.Anything).Return(nil)
	})

	AfterEach(func() {
		mocks.UnsetAuditLogMock()
		helpers.InjectMockAPIHelperInstance(nil)
		helpers.InjectMockIdentityUserApiHelperInstance(nil)
		aaaModels.InjectMockAppPreferenceServiceInstance(nil)
	})

	Context("When proper service type and level are given", func() {
		It("Should return success response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusVerified})
			partnerservice1 := test_helper.CreatePartnerService(ctx, &models.PartnerServiceMapping{ServiceType: utils.Supplier}, supplier1.ID)
			partnerservice2 := test_helper.CreatePartnerService(ctx, &models.PartnerServiceMapping{ServiceType: utils.Transporter, ServiceLevel: utils.Captive}, supplier1.ID)

			log.Printf("inside test case %+v || %+v", partnerservice1.ID, partnerservice1.ServiceType)
			param := psmpb.PartnerServiceObject{
				SupplierId:       supplier1.ID,
				PartnerServiceId: partnerservice1.ID,
				ServiceType:      "Supplier",
				ServiceLevel:     "L2",
			}

			res, _ := new(services.PartnerServiceMappingService).Edit(ctx, &param)

			Expect(res.Message).To(Equal("Partner Service Edited Successfully"))
			Expect(res.Success).To(Equal(true))

			partner := &models.Supplier{}
			database.DBAPM(ctx).Model(&models.Supplier{}).Where("id = ?", supplier1.ID).First(&partner)

			Expect(partner.Status).To(Equal(models.SupplierStatusPending))

			supplier := &models.PartnerServiceMapping{}
			database.DBAPM(ctx).Model(&models.PartnerServiceMapping{}).Where("id = ?", partnerservice2.ID).First(&supplier)

			Expect(supplier.SupplierId).To(Equal(supplier1.ID))
			Expect(supplier.Active).To(Equal(false))

			transporter := &models.PartnerServiceMapping{}
			database.DBAPM(ctx).Model(&models.PartnerServiceMapping{}).Where("id = ?", partnerservice2.ID).First(&transporter)

			Expect(transporter.SupplierId).To(Equal(supplier1.ID))
			Expect(transporter.Active).To(Equal(false))
		})
	})
	Context("When Partner doesn't exist", func() {
		It("Should return failure response", func() {
			param := psmpb.PartnerServiceObject{
				SupplierId:       100,
				PartnerServiceId: 1000,
				ServiceType:      "Supplier",
				ServiceLevel:     "L0",
				TradeLicenseUrl:  "trade_license_url",
				AgreementUrl:     "agreement_url",
			}

			res, _ := new(services.PartnerServiceMappingService).Edit(ctx, &param)

			Expect(res.Message).To(Equal("Partner/Partner Service Not Found"))
			Expect(res.Success).To(Equal(false))
		})
	})
	Context("When service type and level are incompatible", func() {
		It("Should return failure response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			partnerservice1 := test_helper.CreatePartnerService(ctx, &models.PartnerServiceMapping{ServiceType: utils.Supplier}, supplier1.ID)

			param := psmpb.PartnerServiceObject{
				SupplierId:       supplier1.ID,
				PartnerServiceId: partnerservice1.ID,
				ServiceType:      "Supplier",
				ServiceLevel:     "Cash Vendor",
				TradeLicenseUrl:  "trade_license_url",
				AgreementUrl:     "agreement_url",
			}

			res, _ := new(services.PartnerServiceMappingService).Edit(ctx, &param)

			Expect(res.Message).To(Equal("Incompatible Service Type and Service Level"))
			Expect(res.Success).To(Equal(false))
		})
	})
	Context("When service type is edited", func() {
		It("Should return failure response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			partnerservice1 := test_helper.CreatePartnerService(ctx, &models.PartnerServiceMapping{ServiceType: utils.Supplier}, supplier1.ID)

			param := psmpb.PartnerServiceObject{
				SupplierId:       supplier1.ID,
				PartnerServiceId: partnerservice1.ID,
				ServiceType:      "Transporter",
				ServiceLevel:     "Captive",
			}

			res, _ := new(services.PartnerServiceMappingService).Edit(ctx, &param)

			Expect(res.Message).To(Equal("Not allowed to edit Partner Type"))
			Expect(res.Success).To(Equal(false))
		})
	})
})
