package key_account_manager_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	kampb "github.com/voonik/goConnect/api/go/ss2/key_account_manager"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("ListKeyAccountManager", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("List", func() {
		It("Should Respond with all the Key Account Managers", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			accountManager1 := test_helper.CreateKeyAccountManager(ctx, &models.KeyAccountManager{SupplierID: supplier1.ID})
			accountManager2 := test_helper.CreateKeyAccountManager(ctx, &models.KeyAccountManager{SupplierID: supplier1.ID})

			res, err := new(services.KeyAccountManagerService).List(ctx, &kampb.ListParams{SupplierId: supplier1.ID})
			Expect(err).To(BeNil())
			Expect(len(res.Data)).To(Equal(2))

			kamData1 := res.Data[0]
			Expect(kamData1.Name).To(Equal(accountManager1.Name))
			Expect(kamData1.Email).To(Equal(accountManager1.Email))
			Expect(kamData1.Phone).To(Equal(accountManager1.Phone))

			kamData2 := res.Data[1]
			Expect(kamData2.Name).To(Equal(accountManager2.Name))
			Expect(kamData2.Email).To(Equal(accountManager2.Email))
			Expect(kamData2.Phone).To(Equal(accountManager2.Phone))
		})
	})
})
