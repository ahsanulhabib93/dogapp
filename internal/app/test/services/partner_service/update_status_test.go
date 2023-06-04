package partner_service_service_test

import (
	"context"

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
)

var _ = Describe("UpdateStatus", func() {
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

	Context("When service is deactivated and proper service type and level are given", func() {
		It("Should return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			partnerservice := test_helper.CreatePartnerService(ctx, &models.PartnerServiceMapping{}, supplier.ID)

			param := psmpb.PartnerServiceObject{
				SupplierId:       supplier.ID,
				PartnerServiceId: partnerservice.ID,
				Active:           false,
			}

			res, _ := new(services.PartnerServiceMappingService).UpdateStatus(ctx, &param)

			Expect(res.Message).To(Equal("Partner Service Updated Successfully"))
			Expect(res.Success).To(Equal(true))

			service := models.PartnerServiceMapping{}
			database.DBAPM(ctx).Model(models.PartnerServiceMapping{}).Where("id = ?", partnerservice.ID).First(&service)

			Expect(service.Active).To(Equal(false))
		})
	})
	Context("When service is activated and proper service type and level are given", func() {
		It("Should return success response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			partnerservice := test_helper.CreatePartnerService(ctx, &models.PartnerServiceMapping{}, supplier.ID)

			param := psmpb.PartnerServiceObject{
				SupplierId:       supplier.ID,
				PartnerServiceId: partnerservice.ID,
				Active:           true,
			}

			res, _ := new(services.PartnerServiceMappingService).UpdateStatus(ctx, &param)

			Expect(res.Message).To(Equal("Partner Service Updated Successfully"))
			Expect(res.Success).To(Equal(true))

			service := models.PartnerServiceMapping{}
			database.DBAPM(ctx).Model(models.PartnerServiceMapping{}).Where("id = ?", partnerservice.ID).First(&service)

			Expect(service.Active).To(Equal(true))
		})
	})
	Context("When partner doesn't exist", func() {
		It("Should return failure response", func() {
			param := psmpb.PartnerServiceObject{
				SupplierId:       100,
				PartnerServiceId: 101,
				Active:           false,
			}

			res, _ := new(services.PartnerServiceMappingService).UpdateStatus(ctx, &param)

			Expect(res.Message).To(Equal("Partner Not Found"))
			Expect(res.Success).To(Equal(false))
		})
	})
	Context("When partner service doesn't exist", func() {
		It("Should return failure response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := psmpb.PartnerServiceObject{
				SupplierId:       supplier.ID,
				PartnerServiceId: 101,
				Active:           false,
			}

			res, _ := new(services.PartnerServiceMappingService).UpdateStatus(ctx, &param)

			Expect(res.Message).To(Equal("Partner Service Not Found"))
			Expect(res.Success).To(Equal(false))
		})
	})
})
