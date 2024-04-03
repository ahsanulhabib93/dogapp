package helper_tests

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/shopuptech/work"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/goFramework/pkg/worker"

	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("ChangePendingState", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	It("Should update status successfully", func() {
		lastMonth := time.Now().Add(-time.Hour * 24 * 31)
		s1 := test_helper.CreateSupplierWithDateTime(ctx, &models.Supplier{
			SupplierCategoryMappings: []models.SupplierCategoryMapping{{CategoryID: 1}, {CategoryID: 2}},
			SupplierOpcMappings:      []models.SupplierOpcMapping{{ProcessingCenterID: 3}, {ProcessingCenterID: 4}},
		}, lastMonth)
		isPhoneVerified := true
		s2 := test_helper.CreateSupplierWithDateTime(ctx, &models.Supplier{
			IsPhoneVerified: &isPhoneVerified,
			Status:          models.SupplierStatusPending,
		}, lastMonth)

		suppliers := []models.Supplier{}
		database.DBAPM(ctx).Model(&models.Supplier{}).Where("status = ?", models.SupplierStatusFailed).Scan(&suppliers)

		err := helpers.ChangePendingState(&worker.VaccountContext{VaccountID: 1, PortalID: 1}, &work.Job{})
		Expect(err).To(BeNil())

		var count int
		database.DBAPM(ctx).Model(&models.Supplier{}).Where("status = ?", models.SupplierStatusFailed).Count(&count)
		Expect(count).To(Equal(2))

		suppliers = []models.Supplier{}
		database.DBAPM(ctx).Model(&models.Supplier{}).Where("status = ?", models.SupplierStatusFailed).Scan(&suppliers)
		Expect(len(suppliers)).To(Equal(2))

		Expect(suppliers[0].ID).To(Equal(s1.ID))
		Expect(suppliers[1].ID).To(Equal(s2.ID))

		supplierData1 := suppliers[1]
		Expect(supplierData1.Email).To(Equal(s2.Email))
		Expect(supplierData1.Name).To(Equal(s2.Name))
		Expect(supplierData1.Phone).To(Equal(s2.Phone))
		Expect(supplierData1.AlternatePhone).To(Equal(s2.AlternatePhone))
		Expect(supplierData1.BusinessName).To(Equal(s2.BusinessName))
		Expect(*supplierData1.IsPhoneVerified).To(Equal(true))
		Expect(supplierData1.Status).To(Equal(models.SupplierStatusFailed))
		Expect(supplierData1.CreatedAt.Day()).To(Equal(lastMonth.Day()))
		Expect(supplierData1.UpdatedAt.Day()).To(Equal(time.Now().Day()))
	})

	It("Should not update supplier status if not more than 1 week old", func() {
		aaaModels.InjectMockAppPreferenceServiceInstance(mocks.GetAppPreferenceMock(map[string]interface{}{
			"supplier_auto_status_change_duration": int64(20),
		}))

		date := time.Now().Add(-time.Hour * 24 * 31)
		isPhoneVerified := true
		test_helper.CreateSupplierWithDateTime(ctx, &models.Supplier{
			IsPhoneVerified: &isPhoneVerified,
			Status:          models.SupplierStatusVerified,
		}, date)

		date = time.Now().Add(-time.Hour * 24 * 8)
		test_helper.CreateSupplierWithDateTime(ctx, &models.Supplier{
			IsPhoneVerified: &isPhoneVerified,
			Status:          models.SupplierStatusPending,
		}, date)

		err := helpers.ChangePendingState(&worker.VaccountContext{VaccountID: 1, PortalID: 1}, &work.Job{})
		Expect(err).To(BeNil())

		var count int
		database.DBAPM(ctx).Model(&models.Supplier{}).Where("status = ?", models.SupplierStatusFailed).Count(&count)
		Expect(count).To(Equal(0))
	})
})
