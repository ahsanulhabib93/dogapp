package partner_service_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	psmpb "github.com/voonik/goConnect/api/go/ss2/partner_service_mapping"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("PartnerTypesList", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("When user has global permission", func() {
		BeforeEach(func() {
			test_helper.SetContextUser(&ctx, 1, []string{"supplierpanel:allservices:view"})
		})

		It("Should return all service types data", func() {
			res, _ := new(services.PartnerServiceMappingService).PartnerTypesList(ctx, &psmpb.PartnerServiceObject{})

			Expect(len(res.PartnerServiceTypeMappings)).To(Equal(2))

			supplier := res.PartnerServiceTypeMappings[0]
			Expect(supplier.PartnerType).To(Equal("Supplier"))
			Expect(len(supplier.ServiceTypes)).To(Equal(5))
			Expect(supplier.ServiceTypes).To(Equal([]string{"L0", "L1", "L2", "L3", "Hlc"}))

			transport := res.PartnerServiceTypeMappings[1]
			Expect(transport.PartnerType).To(Equal("Transporter"))
			Expect(len(transport.ServiceTypes)).To(Equal(5))
			Expect(transport.ServiceTypes).To(Equal([]string{"Captive", "Driver", "CashVendor", "RedxHubVendor", "CreditVendor"}))
		})
	})

	Context("When user has supplier permission", func() {
		BeforeEach(func() {
			test_helper.SetContextUser(&ctx, 1, []string{"supplierpanel:supplierservice:view"})
		})

		It("Should return only supplier service type data", func() {
			res, _ := new(services.PartnerServiceMappingService).PartnerTypesList(ctx, &psmpb.PartnerServiceObject{})

			Expect(len(res.PartnerServiceTypeMappings)).To(Equal(1))

			supplier := res.PartnerServiceTypeMappings[0]
			Expect(supplier.PartnerType).To(Equal("Supplier"))
			Expect(len(supplier.ServiceTypes)).To(Equal(5))
			Expect(supplier.ServiceTypes).To(Equal([]string{"L0", "L1", "L2", "L3", "Hlc"}))
		})
	})
})
