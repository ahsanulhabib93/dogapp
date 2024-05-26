package helper_tests

import (
	"context"
	"fmt"

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

	Context("ValidatePartnerSericeEdit", func() {
		var incoming, existing helpers.PartnerServiceEditEntity
		BeforeEach(func() {
			cashVendorServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Transporter, "Cash Vendor")
			creditVendorServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Transporter, "Credit Vendor")
			aaaModels.InjectMockAppPreferenceServiceInstance(mocks.GetAppPreferenceMock(map[string]interface{}{
				"edit_allowed_service_levels": []string{fmt.Sprint(cashVendorServiceLevel.ID), fmt.Sprint(creditVendorServiceLevel.ID)},
			}))
			incoming = helpers.PartnerServiceEditEntity{
				ServiceType:    utils.Transporter,
				ServiceLevelId: cashVendorServiceLevel.ID,
			}
			existing = helpers.PartnerServiceEditEntity{
				ServiceType:    utils.Transporter,
				ServiceLevelId: creditVendorServiceLevel.ID,
			}
		})

		It("Should return true if incomingServiceLevel is same as existingServiceLevel", func() {
			cashVendorServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Transporter, "Cash Vendor")
			incoming.ServiceLevelId = cashVendorServiceLevel.ID
			existing.ServiceLevelId = cashVendorServiceLevel.ID
			Expect(helpers.ValidatePartnerSericeEdit(ctx, incoming, existing)).To(Equal(true))
		})

		It("Should return true if incomingServiceLevel is different from existingServiceLevel and both are allowed", func() {
			cashVendorServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Transporter, "Cash Vendor")
			creditVendorServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Transporter, "Credit Vendor")
			incoming.ServiceLevelId = cashVendorServiceLevel.ID
			existing.ServiceLevelId = creditVendorServiceLevel.ID
			Expect(helpers.ValidatePartnerSericeEdit(ctx, incoming, existing)).To(Equal(true))
		})

		It("Should return false if incomingServiceLevel is different from existingServiceLevel and either is not allowed", func() {
			driverServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Transporter, "Driver")
			creditVendorServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Transporter, "Credit Vendor")
			incoming.ServiceLevelId = driverServiceLevel.ID
			existing.ServiceLevelId = creditVendorServiceLevel.ID
			Expect(helpers.ValidatePartnerSericeEdit(ctx, incoming, existing)).To(Equal(false))
		})

		It("Should return false if incomingServiceLevel is different from existingServiceLevel and both are not allowed", func() {
			driverServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Transporter, "Driver")
			captiveServiceLevel := helpers.GetServiceLevelByTypeAndName(ctx, utils.Transporter, "Captive")
			incoming.ServiceLevelId = driverServiceLevel.ID
			existing.ServiceLevelId = captiveServiceLevel.ID
			Expect(helpers.ValidatePartnerSericeEdit(ctx, incoming, existing)).To(Equal(false))
		})
	})
})
