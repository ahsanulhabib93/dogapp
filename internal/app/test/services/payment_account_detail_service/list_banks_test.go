package payment_account_detail_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("ListBanks", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("ListBanks", func() {
		It("Should Respond with all the Banks", func() {
			bank1 := test_helper.CreateBank(ctx, &models.Bank{})
			bank2 := test_helper.CreateBank(ctx, &models.Bank{})

			res, err := new(services.PaymentAccountDetailService).ListBanks(ctx, &paymentpb.ListBankParams{})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(2))

			Expect(res.Data[0].Name).To(Equal(bank1.Name))
			Expect(res.Data[1].Name).To(Equal(bank2.Name))
		})
	})
})
