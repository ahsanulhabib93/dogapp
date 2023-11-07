package seller_service_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
)

var _ = Describe("Get seller by condition", func() {
	var ctx context.Context
	var mockAudit *mocks.AuditLogMock
	userId := uint64(101)

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		mocks.UnsetOpcMock()

		test_helper.SetContextUser(&ctx, userId, []string{})
		mockAudit = mocks.SetAuditLogMock()
		mockAudit.On("RecordAuditAction", ctx, mock.Anything).Return(nil)
	})

	AfterEach(func() {
		mocks.UnsetAuditLogMock()
		helpers.InjectMockAPIHelperInstance(nil)
		helpers.InjectMockIdentityUserApiHelperInstance(nil)
		aaaModels.InjectMockAppPreferenceServiceInstance(nil)
	})

	Context("When condition is not passed in the param", func() {
		It("Should return error", func() {
			param := spb.GetSellerByConditionParams{}
			res, err := new(services.SellerService).GetSellerByCondition(ctx, &param)
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("no condition specified"))
			Expect(err).To(BeNil())
		})
	})

	Context("When condition is passed in the param", func() {
		It("Should return seller details", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{
				UserID:    userId,
				BrandName: "test_brand",
			})
			fmt.Print("user_id = ?", seller1.UserID)
			param := spb.GetSellerByConditionParams{
				Condition: fmt.Sprint("user_id = ", seller1.UserID),
			}

			res, err := new(services.SellerService).GetSellerByCondition(ctx, &param)
			Expect(len(res.Seller)).To(Equal(1))
			Expect(res.Seller[0].BrandName).To(Equal(seller1.BrandName))
			Expect(res.Seller[0].UserId).To(Equal(seller1.UserID))
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("fetched seller details successfully"))
			Expect(err).To(BeNil())
		})
	})

	Context("When condition and fields are passed in the param", func() {
		It("Should return seller details", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{
				UserID:    userId,
				BrandName: "test_brand",
			})
			seller2 := test_helper.CreateSeller(ctx, &models.Seller{
				UserID:    uint64(200),
				BrandName: "test_brand",
			})
			fmt.Print("user_id = ?", seller1.UserID)
			param := spb.GetSellerByConditionParams{
				Condition: fmt.Sprintf("brand_name = '%s'", seller2.BrandName),
				Fields:    fmt.Sprint("user_id, brand_name, primary_phone"),
			}

			res, err := new(services.SellerService).GetSellerByCondition(ctx, &param)
			Expect(len(res.Seller)).To(Equal(2))
			Expect(res.Seller[0].BrandName).To(Equal(seller1.BrandName))
			Expect(res.Seller[0].UserId).To(Equal(seller1.UserID))
			Expect(res.Seller[0].PrimaryPhone).To(Equal(seller1.PrimaryPhone))
			Expect(res.Seller[1].BrandName).To(Equal(seller2.BrandName))
			Expect(res.Seller[1].UserId).To(Equal(seller2.UserID))
			Expect(res.Seller[1].PrimaryPhone).To(Equal(seller2.PrimaryPhone))
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("fetched seller details successfully"))
			Expect(err).To(BeNil())
		})
	})

	Context("When invalid condition is passed in the param", func() {
		It("Should return seller details", func() {
			seller1 := test_helper.CreateSeller(ctx, &models.Seller{
				UserID:    userId,
				BrandName: "test_brand",
			})
			fmt.Print("user_id = ?", seller1.UserID)
			param := spb.GetSellerByConditionParams{
				Condition: fmt.Sprint("brand_name = 'test'"),
				Fields:    fmt.Sprint("user_id, brand_name"),
			}

			res, err := new(services.SellerService).GetSellerByCondition(ctx, &param)
			Expect(len(res.Seller)).To(Equal(0))
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("seller not found"))
			Expect(err).To(BeNil())
		})
	})
})
