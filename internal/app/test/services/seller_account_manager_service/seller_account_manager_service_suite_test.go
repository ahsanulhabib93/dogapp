package seller_account_manager_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/ss2/internal/app/test"
)

func TestSellerAccountManagerService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SellerAccountManagerService Suite")
}

var _ = AfterEach(func() {
	test.Cleaner.Clean("seller_account_managers")
})
