package seller_account_manager_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	sampb "github.com/voonik/goConnect/api/go/ss2/seller_account_manager"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/services"
)

var _ = Describe("SellerAccountManager", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("List", func() {
		It("Should return empty response for no params", func() {
			resp, err := new(services.SellerAccountManagerService).List(ctx, &sampb.ListParams{})

			Expect(err).To(BeNil())
			Expect(resp).To(BeNil())
		})
	})
})
