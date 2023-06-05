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
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("AddPartnerService", func() {
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
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})

			param := psmpb.PartnerServiceObject{
				SupplierId:   supplier1.ID,
				ServiceType:  "Transporter",
				ServiceLevel: "Captive",
			}

			res, _ := new(services.PartnerServiceMappingService).Add(ctx, &param)

			Expect(res.Message).To(Equal("Partner Service Added Successfully"))
			Expect(res.Success).To(Equal(true))

			partnerServiceObj := &models.PartnerServiceMapping{}
			database.DBAPM(ctx).Model(&models.PartnerServiceMapping{}).Last(&partnerServiceObj)

			Expect(partnerServiceObj.SupplierId).To(Equal(supplier1.ID))
			Expect(partnerServiceObj.ServiceType).To(Equal(utils.Transporter))
			Expect(partnerServiceObj.ServiceLevel).To(Equal(utils.Captive))
		})
	})
	Context("When partner service already exist for a user", func() {
		It("Should return failure response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			partnerservice1 := test_helper.CreatePartnerServiceMapping(ctx, &models.PartnerServiceMapping{ServiceType: utils.Supplier}, supplier1.ID)

			param := psmpb.PartnerServiceObject{
				SupplierId:       supplier1.ID,
				PartnerServiceId: partnerservice1.ID,
				ServiceType:      "Supplier",
				ServiceLevel:     "Hlc",
			}

			res, _ := new(services.PartnerServiceMappingService).Add(ctx, &param)

			Expect(res.Message).To(Equal("Error while creating Partner Service: Error 1062: Duplicate entry '1-1-1' for key 'partner_service_mappings.idx_partner_service'"))
			Expect(res.Success).To(Equal(false))
		})
	})
	Context("When service level and service type are not passed", func() {
		It("Should return failure response", func() {
			param := psmpb.PartnerServiceObject{
				SupplierId: 100,
			}

			res, _ := new(services.PartnerServiceMappingService).Add(ctx, &param)

			Expect(res.Message).To(Equal("Invalid Service Type and/or Service Level"))
			Expect(res.Success).To(Equal(false))
		})
	})
	Context("When supplier doesn't exist", func() {
		It("Should return failure response", func() {
			param := psmpb.PartnerServiceObject{
				SupplierId:   100,
				ServiceType:  "Supplier",
				ServiceLevel: "L0",
			}

			res, _ := new(services.PartnerServiceMappingService).Add(ctx, &param)

			Expect(res.Message).To(Equal("Partner Not Found"))
			Expect(res.Success).To(Equal(false))
		})
	})
	Context("When service type and level are incompatible", func() {
		It("Should return failure response", func() {
			param := psmpb.PartnerServiceObject{
				SupplierId:   100,
				ServiceType:  "Supplier",
				ServiceLevel: "Cash Vendor",
			}

			res, _ := new(services.PartnerServiceMappingService).Add(ctx, &param)

			Expect(res.Message).To(Equal("Incompatible Service Type and Service Level"))
			Expect(res.Success).To(Equal(false))
		})
	})
})
