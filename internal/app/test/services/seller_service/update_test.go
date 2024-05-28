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

var _ = Describe("Update", func() {
	var ctx context.Context
	var mockAudit *mocks.AuditLogMock

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
			param := spb.UpdateParams{}
			res, err := new(services.SellerService).Update(ctx, &param)
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("param not specified"))
			Expect(err).To(BeNil())
		})
	})

	Context("When invalid userID is given", func() {
		It("Should return error", func() {
			test_helper.CreateSeller(ctx, &models.Seller{})
			param := spb.UpdateParams{
				Id:     100,
				Seller: &spb.SellerObject{},
			}
			res, err := new(services.SellerService).Update(ctx, &param)
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("seller not found"))
			Expect(err).To(BeNil())
		})
	})

	Context("When params are given", func() {
		It("Should update seller", func() {
			seller := test_helper.CreateSeller(ctx, &models.Seller{UserID: 101})
			newSeller := &spb.SellerObject{PrimaryEmail: "abc@gmail.com"}
			param := spb.UpdateParams{
				Id:     101,
				Seller: newSeller,
			}
			res, err := new(services.SellerService).Update(ctx, &param)
			database.DBAPM(ctx).Model(&models.Seller{}).Where("user_id = ?", 101).Scan(&seller)
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("seller details updated successfully"))
			Expect(seller.PrimaryEmail).To(Equal("abc@gmail.com"))
			Expect(seller.PrimaryPhone).To(Equal(seller.PrimaryPhone))
			Expect(err).To(BeNil())
		})
	})

	Context("When params with multiple fields are given", func() {
		It("Should update seller", func() {
			seller := test_helper.CreateSeller(ctx, &models.Seller{UserID: 101})
			newSeller := &spb.SellerObject{PrimaryEmail: "abc@gmail.com", PrimaryPhone: "880123456789"}
			param := spb.UpdateParams{
				Id:     101,
				Seller: newSeller,
			}
			res, err := new(services.SellerService).Update(ctx, &param)
			database.DBAPM(ctx).Model(&models.Seller{}).Where("user_id = ?", 101).Scan(&seller)
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("seller details updated successfully"))
			Expect(seller.PrimaryEmail).To(Equal("abc@gmail.com"))
			Expect(seller.PrimaryPhone).To(Equal("880123456789"))
			Expect(err).To(BeNil())
		})
	})
})
