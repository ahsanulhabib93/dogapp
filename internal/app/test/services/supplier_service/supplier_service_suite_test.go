package supplier_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/ss2/internal/app/test"
)

func TestSupplierService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SupplierService Suite")
}

var _ = AfterEach(func() {
	test.Cleaner.Clean("suppliers", "supplier_addresses", "supplier_category_mappings",
		"supplier_sa_mappings", "payment_account_details", "supplier_opc_mappings")
})
