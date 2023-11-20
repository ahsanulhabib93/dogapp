package seller_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	cmtPb "github.com/voonik/goConnect/api/go/cmt/product"
	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("Approve Products", func() {
	var ctx context.Context
	var mockAudit *mocks.AuditLogMock
	var apiHelperInstance *mocks.APIHelperInterface

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

	Context("Fail Case", func() {
		It("Should return invalid msg", func() {
			res, err := new(services.SellerService).ApproveProducts(ctx, &spb.ApproveProductsParams{})
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal(utils.Failure))
			Expect(res.Message).To(Equal("Failed to approve the products - atleast one id should be present"))
		})

		It("Should return seller not found msg", func() {
			test_helper.SetContextUser(&ctx, 11, []string{})

			res, err := new(services.SellerService).ApproveProducts(ctx, &spb.ApproveProductsParams{Ids: []uint64{1312}})
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal(utils.Failure))
			Expect(res.Message).To(Equal("Seller Not Found"))
		})

		It("Should return pickup address or pan number missing msg", func() {
			test_helper.SetContextUser(&ctx, 10, []string{})

			seller := models.Seller{UserID: 10}
			database.DBAPM(ctx).Create(&seller)
			res, err := new(services.SellerService).ApproveProducts(ctx, &spb.ApproveProductsParams{Ids: []uint64{1312}})
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal(utils.Failure))
			Expect(res.Message).To(Equal("Pick Up Address or Pan number is missing"))
		})

		It("Should return product approved count 0", func() {
			test_helper.SetContextUser(&ctx, 10, []string{})

			seller := models.Seller{UserID: 10, PanNumber: "24123413", ActivationState: 3}
			database.DBAPM(ctx).Create(&seller)

			seller.VendorAddresses = []*models.VendorAddress{{SellerID: int(seller.ID)}}
			seller.BusinessType = utils.Manufacturer
			seller.ColorCode = utils.Black
			database.DBAPM(ctx).Save(&seller)

			apiHelperInstance = new(mocks.APIHelperInterface)
			helpers.InjectMockAPIHelperInstance(apiHelperInstance)
			apiHelperInstance.On("CmtApproveItems", ctx, &cmtPb.ApproveItemParams{ProductIds: []uint64{13, 31}, UserId: seller.UserID, State: 3}).Return(nil)

			res, err := new(services.SellerService).ApproveProducts(ctx, &spb.ApproveProductsParams{Ids: []uint64{13, 31}})
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal(utils.Success))
			Expect(res.Message).To(Equal("The total number of products approved are 0"))
		})
	})

	Context("Success Case", func() {
		It("Should return product approved count 1", func() {
			test_helper.SetContextUser(&ctx, 10, []string{})

			seller := models.Seller{UserID: 10, PanNumber: "24123413", ActivationState: 3}
			database.DBAPM(ctx).Create(&seller)

			seller.VendorAddresses = []*models.VendorAddress{{SellerID: int(seller.ID)}}
			seller.BusinessType = utils.Manufacturer
			seller.ColorCode = utils.Black
			database.DBAPM(ctx).Save(&seller)

			apiHelperInstance = new(mocks.APIHelperInterface)
			helpers.InjectMockAPIHelperInstance(apiHelperInstance)
			apiHelperInstance.On("CmtApproveItems", ctx, &cmtPb.ApproveItemParams{ProductIds: []uint64{12, 31}, UserId: seller.UserID, State: 3}).Return(&cmtPb.ItemCountResponse{Count: 1})

			res, err := new(services.SellerService).ApproveProducts(ctx, &spb.ApproveProductsParams{Ids: []uint64{12, 31}})
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal(utils.Success))
			Expect(res.Message).To(Equal("The total number of products approved are 1"))
		})
	})

})
