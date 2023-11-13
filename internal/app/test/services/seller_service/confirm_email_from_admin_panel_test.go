package seller_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("Confirm email from admin panel", func() {
	var ctx context.Context
	var mockAudit *mocks.AuditLogMock
	userID := uint64(200)

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		mocks.UnsetOpcMock()

		mockAudit = mocks.SetAuditLogMock()
		mockAudit.On("RecordAuditAction", ctx, mock.Anything).Return(nil)
	})

	AfterEach(func() {
		mocks.UnsetAuditLogMock()
		helpers.InjectMockAPIHelperInstance(nil)
		helpers.InjectMockIdentityUserApiHelperInstance(nil)
		aaaModels.InjectMockAppPreferenceServiceInstance(nil)
	})

	Context("When no params are given", func() {
		It("Should return error", func() {
			param := spb.GetByUserIDParams{}
			res, err := new(services.SellerService).ConfirmEmailFromAdminPanel(ctx, &param)
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("param not specified"))
			Expect(err).To(BeNil())
		})
	})

	Context("When params are given", func() {
		It("Should confirm email", func() {
			seller := test_helper.CreateSeller(ctx, &models.Seller{UserID: userID})
			param := spb.GetByUserIDParams{
				UserId: userID,
			}
			res, err := new(services.SellerService).ConfirmEmailFromAdminPanel(ctx, &param)
			database.DBAPM(ctx).Model(&models.Seller{}).Where("user_id = ?", userID).Scan(&seller)
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("email confirmed successfully"))
			Expect(seller.EmailConfirmed).To(BeTrue())
			Expect(err).To(BeNil())
		})
	})

	Context("When invalid userID are given", func() {
		It("Should return error", func() {
			test_helper.CreateSeller(ctx, &models.Seller{UserID: userID})
			param := spb.GetByUserIDParams{
				UserId: 100,
			}
			res, err := new(services.SellerService).ConfirmEmailFromAdminPanel(ctx, &param)
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("seller not found"))
			Expect(err).To(BeNil())
		})
	})
})
