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
)

var _ = Describe("EditKeyAccountManager", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Editing existing Account Manager", func() {
		It("Should update and return success response", func() {
			accountManager := test_helper.CreateKeyAccountManager(ctx, &models.KeyAccountManager{})
			param := &kampb.KeyAccountManagerObject{
				Id:    accountManager.ID,
				Name:  "Name",
				Email: "Email",
				Phone: "Phone",
			}
			res, err := new(services.KeyAccountManagerService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("KeyAccountManager Edited Successfully"))

			database.DBAPM(ctx).Model(&models.Supplier{}).First(&accountManager, accountManager.ID)
			Expect(accountManager.Email).To(Equal(param.Email))
			Expect(accountManager.Name).To(Equal(param.Name))
			Expect(accountManager.Phone).To(Equal(param.Phone))
		})
	})

	Context("Editing only name of existing account manager", func() {
		It("Should return error response", func() {
			accountManager := test_helper.CreateKeyAccountManager(ctx, &models.KeyAccountManager{})
			param := &kampb.KeyAccountManagerObject{
				Id:   accountManager.ID,
				Name: "Name",
			}
			res, err := new(services.KeyAccountManagerService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("KeyAccountManager Edited Successfully"))

			updatedMananager := &models.KeyAccountManager{}
			database.DBAPM(ctx).Model(accountManager).First(&updatedMananager, accountManager.ID)
			Expect(updatedMananager.Name).To(Equal(param.Name))
			Expect(updatedMananager.Email).To(Equal(accountManager.Email))
			Expect(updatedMananager.Phone).To(Equal(accountManager.Phone))
		})
	})

	Context("Editing invalid Account Manager", func() {
		It("Should return error response", func() {
			param := &kampb.KeyAccountManagerObject{Id: 1000}
			res, err := new(services.KeyAccountManagerService).Edit(ctx, param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("KeyAccountManager Not Found"))
		})
	})
})
