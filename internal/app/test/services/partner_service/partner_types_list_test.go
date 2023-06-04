package partner_service_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	psmpb "github.com/voonik/goConnect/api/go/ss2/partner_service_mapping"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("PartnerTypesList", func() {
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

	Context("When no params are given", func() {
		It("Should return success response", func() {
			param := psmpb.PartnerServiceObject{}

			res, _ := new(services.PartnerServiceMappingService).PartnerTypesList(ctx, &param)

			Expect(len(res.PartnerServiceTypeMappings)).To(Equal(2))

			supplier := res.PartnerServiceTypeMappings[0]

			Expect(supplier.PartnerType).To(Equal("Supplier"))
			Expect(len(supplier.ServiceTypes)).To(Equal(5))
			Expect(supplier.ServiceTypes).To(Equal([]string{"L0", "L1", "L2", "L3", "Hlc"}))

			transport := res.PartnerServiceTypeMappings[1]

			Expect(transport.PartnerType).To(Equal("Transporter"))
			Expect(len(transport.ServiceTypes)).To(Equal(3))
			Expect(transport.ServiceTypes).To(Equal([]string{"Captive", "Driver", "CashVendor"}))
		})
	})
})
