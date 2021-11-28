package key_account_manager_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	kampb "github.com/voonik/goConnect/api/go/ss2/key_account_manager"
	"github.com/voonik/goFramework/pkg/database"
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

	Context("Add", func() {
		It("Should create key account manager and return success", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{SupplierType: utils.Hlc})
			param := kampb.KeyAccountManagerParam{
				SupplierId: supplier.ID,
				Name:       "Name",
				Email:      "Email",
				Phone:      "Phone",
			}
			res, err := new(services.KeyAccountManagerService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("KeyAccountManager Added Successfully"))

			accountManagers := []*models.KeyAccountManager{{}}
			database.DBAPM(ctx).Model(supplier).Association("KeyAccountManagers").Find(&accountManagers)
			Expect(len(accountManagers)).To(Equal(1))
			accountManager := accountManagers[0]

			Expect(accountManager.Name).To(Equal(param.Name))
			Expect(accountManager.Email).To(Equal(param.Email))
			Expect(accountManager.Phone).To(Equal(param.Phone))
		})
	})

	Context("While adding account manager for invalid Supplier ID", func() {
		It("Should return error response", func() {
			param := kampb.KeyAccountManagerParam{
				SupplierId: 1000,
				Name:       "Name",
			}
			res, err := new(services.KeyAccountManagerService).Add(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Not Found"))
		})
	})
})
