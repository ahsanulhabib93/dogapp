package payment_account_detail_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/ss2/internal/app/test"
)

func TestPaymentAccountDetailService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PaymentAccountDetailService Suite")
}

var _ = AfterEach(func() {
	test.Cleaner.Clean("suppliers", "payment_account_details")
})
