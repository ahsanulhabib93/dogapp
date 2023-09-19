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

	Context("ValidateSericeLevelEdit", func() {
		var incomingServiceLevel, existingServiceLevel utils.SupplierType

		BeforeEach(func() {
			aaaModels.InjectMockAppPreferenceServiceInstance(mocks.GetAppPreferenceMock(map[string]interface{}{
				"edit_allowed_service_levels": []string{"CashVendor", "CreditVendor"},
			}))
		})

		It("Should return true if incomingServiceLevel is same as existingServiceLevel", func() {
			incomingServiceLevel = utils.CashVendor
			existingServiceLevel = utils.CashVendor
			Expect(helpers.ValidateSericeLevelEdit(ctx, incomingServiceLevel, existingServiceLevel)).To(Equal(true))
		})

		It("Should return true if incomingServiceLevel is different from existingServiceLevel and both are allowed", func() {
			incomingServiceLevel = utils.CashVendor
			existingServiceLevel = utils.CreditVendor
			Expect(helpers.ValidateSericeLevelEdit(ctx, incomingServiceLevel, existingServiceLevel)).To(Equal(true))
		})

		It("Should return false if incomingServiceLevel is different from existingServiceLevel and either is not allowed", func() {
			incomingServiceLevel = utils.Driver
			existingServiceLevel = utils.CreditVendor
			Expect(helpers.ValidateSericeLevelEdit(ctx, incomingServiceLevel, existingServiceLevel)).To(Equal(false))
		})

		It("Should return false if incomingServiceLevel is different from existingServiceLevel and both are not allowed", func() {
			incomingServiceLevel = utils.Driver
			existingServiceLevel = utils.Captive
			Expect(helpers.ValidateSericeLevelEdit(ctx, incomingServiceLevel, existingServiceLevel)).To(Equal(false))
		})
	})
})
