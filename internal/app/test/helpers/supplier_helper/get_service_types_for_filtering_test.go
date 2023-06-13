package helper_tests

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("GetServiceTypesForFiltering", func() {
	var ctx context.Context
	BeforeEach(func() {
		test_utils.GetContext(&ctx)

	})

	Context("When user has global permission and service types are blank", func() {
		BeforeEach(func() {
			test_helper.SetContextUser(&ctx, 1, []string{"supplierpanel:allservices:view"})
		})

		It("Should return all service types", func() {
			serviceTypes := []string{}
			allowedServiceTypes := helpers.GetServiceTypesForFiltering(ctx, serviceTypes)

			Expect(len(allowedServiceTypes)).To(Equal(2))
			Expect(allowedServiceTypes[0]).To(Equal("Supplier"))
			Expect(allowedServiceTypes[1]).To(Equal("Transporter"))
		})
	})

	Context("For framework user", func() {
		BeforeEach(func() {
			test_utils.SetPermission(&ctx, []string{})
		})

		It("Should return all service types", func() {
			serviceTypes := []string{}
			allowedServiceTypes := helpers.GetServiceTypesForFiltering(ctx, serviceTypes)

			Expect(utils.GetCurrentUserID(ctx)).To(Equal(nil))

			Expect(len(allowedServiceTypes)).To(Equal(2))
			Expect(allowedServiceTypes[0]).To(Equal("Supplier"))
			Expect(allowedServiceTypes[1]).To(Equal("Transporter"))
		})
	})

	Context("When user has global permission and service types are present", func() {
		BeforeEach(func() {
			test_helper.SetContextUser(&ctx, 1, []string{"supplierpanel:allservices:view"})
		})

		It("Should return given service types", func() {
			serviceTypes := []string{"Supplier"}
			allowedServiceTypes := helpers.GetServiceTypesForFiltering(ctx, serviceTypes)

			Expect(len(allowedServiceTypes)).To(Equal(1))
			Expect(allowedServiceTypes[0]).To(Equal("Supplier"))
		})
	})

	Context("When user has Supplier permission and service types are blank", func() {
		BeforeEach(func() {
			test_helper.SetContextUser(&ctx, 1, []string{"supplierpanel:supplierservice:view"})
		})

		It("Should return Supplier service types", func() {
			serviceTypes := []string{}
			allowedServiceTypes := helpers.GetServiceTypesForFiltering(ctx, serviceTypes)

			Expect(len(allowedServiceTypes)).To(Equal(1))
			Expect(allowedServiceTypes[0]).To(Equal("Supplier"))
		})
	})

	Context("When user has Supplier permission and service types is different", func() {
		BeforeEach(func() {
			test_helper.SetContextUser(&ctx, 1, []string{"supplierpanel:supplierservice:view"})
		})

		It("Should return empty service types", func() {
			serviceTypes := []string{"Transport"}
			allowedServiceTypes := helpers.GetServiceTypesForFiltering(ctx, serviceTypes)

			Expect(len(allowedServiceTypes)).To(Equal(0))
		})
	})
})
