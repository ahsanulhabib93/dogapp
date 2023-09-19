package helper_tests

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("PartnerServiceHelpers", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("ParseServiceLevels", func() {
		It("Should return parsed service levels", func() {
			serviceLevels := []string{"CashVendor", "CreditVendor", "GoldVendor"}
			Expect(helpers.ParseServiceLevels(serviceLevels)).To(Equal(
				[]utils.SupplierType{utils.CashVendor, utils.CreditVendor, utils.SupplierType(0)}))
		})
	})

	Context("ValidatePartnerSericeEdit", func() {
		var incoming, existing helpers.PartnerServiceEditEntity

		BeforeEach(func() {
			aaaModels.InjectMockAppPreferenceServiceInstance(mocks.GetAppPreferenceMock(map[string]interface{}{
				"edit_allowed_service_levels": []string{"CashVendor", "CreditVendor"},
			}))
			incoming = helpers.PartnerServiceEditEntity{
				ServiceType:  utils.Transporter,
				ServiceLevel: utils.CashVendor,
			}
			existing = helpers.PartnerServiceEditEntity{
				ServiceType:  utils.Transporter,
				ServiceLevel: utils.CreditVendor,
			}
		})

		It("Should return true if incomingServiceLevel is same as existingServiceLevel", func() {
			incoming.ServiceLevel = utils.CashVendor
			existing.ServiceLevel = utils.CashVendor
			Expect(helpers.ValidatePartnerSericeEdit(ctx, incoming, existing)).To(Equal(true))
		})

		It("Should return true if incomingServiceLevel is different from existingServiceLevel and both are allowed", func() {
			incoming.ServiceLevel = utils.CashVendor
			existing.ServiceLevel = utils.CreditVendor
			Expect(helpers.ValidatePartnerSericeEdit(ctx, incoming, existing)).To(Equal(true))
		})

		It("Should return false if incomingServiceLevel is different from existingServiceLevel and either is not allowed", func() {
			incoming.ServiceLevel = utils.Driver
			existing.ServiceLevel = utils.CreditVendor
			Expect(helpers.ValidatePartnerSericeEdit(ctx, incoming, existing)).To(Equal(false))
		})

		It("Should return false if incomingServiceLevel is different from existingServiceLevel and both are not allowed", func() {
			incoming.ServiceLevel = utils.Driver
			existing.ServiceLevel = utils.Captive
			Expect(helpers.ValidatePartnerSericeEdit(ctx, incoming, existing)).To(Equal(false))
		})
	})
})
