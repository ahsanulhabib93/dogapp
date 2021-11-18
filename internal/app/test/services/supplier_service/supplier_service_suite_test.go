package supplier_service_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/voonik/supplier_service/internal/app/test"
)

func TestExpenseConfigurationService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SupplierService Suite")
}

var _ = AfterEach(func() {
	test.Cleaner.Clean("suppliers")
})
