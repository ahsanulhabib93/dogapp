package partner_service_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	psmpb "github.com/voonik/goConnect/api/go/ss2/partner_service_mapping"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/services"
)

var _ = Describe("PartnerTypesList", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		test_utils.SetPermission(&ctx, []string{"supplierpanel:allservices:view"})
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
			Expect(len(transport.ServiceTypes)).To(Equal(4))
			Expect(transport.ServiceTypes).To(Equal([]string{"Captive", "Driver", "CashVendor", "RedxHubVendor"}))
		})
	})
})
