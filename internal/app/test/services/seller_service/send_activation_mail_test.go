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
)

var _ = Describe("Send Activation Mail", func() {
	var ctx context.Context
	var mockAudit *mocks.AuditLogMock
	var seller models.Seller

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

	Context("fail Case", func() {
		It("Should return status false for invalid param", func() {
			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("Seller Ids and Action Should be Present"))
		})

		It("Should return status false for seller not found", func() {
			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate"}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("Seller not found"))
		})

		It("Should return status false for seller bank detail not found", func() {
			seller = models.Seller{UserID: 1}
			database.DBAPM(ctx).Create(&seller)

			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate"}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("1: Seller Pan Number, Bank Detail, MOU and Email should be confirmed"))
		})

		It("Should return status false for > 1 vendor address", func() {
			seller = models.Seller{
				UserID:         1,
				PanNumber:      "PAN123",
				EmailConfirmed: true,
				MouAgreed:      true,
			}
			database.DBAPM(ctx).Create(&seller)

			sellerBankDetail := models.SellerBankDetail{SellerID: int(seller.ID)}
			database.DBAPM(ctx).Create(&sellerBankDetail)

			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate"}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("1: At least one address should be present"))
		})

		It("Should return status false for no verify status address", func() {
			seller = models.Seller{
				UserID:         1,
				PanNumber:      "PAN123",
				EmailConfirmed: true,
				MouAgreed:      true,
			}
			database.DBAPM(ctx).Create(&seller)

			sellerBankDetail := models.SellerBankDetail{SellerID: int(seller.ID)}
			database.DBAPM(ctx).Create(&sellerBankDetail)

			vendorAddress := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED", DefaultAddress: true}
			database.DBAPM(ctx).Create(&vendorAddress)

			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate"}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("1: Make at least one address as default"))
		})

		It("Should return status false for no default status address", func() {
			seller = models.Seller{
				UserID:         1,
				PanNumber:      "PAN123",
				EmailConfirmed: true,
				MouAgreed:      true,
			}
			database.DBAPM(ctx).Create(&seller)

			sellerBankDetail := models.SellerBankDetail{SellerID: int(seller.ID)}
			database.DBAPM(ctx).Create(&sellerBankDetail)

			vendorAddress := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED"}
			database.DBAPM(ctx).Create(&vendorAddress)

			vendorAddress2 := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED", VerificationStatus: "VERIFIED"}
			database.DBAPM(ctx).Create(&vendorAddress2)

			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate"}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("1: Make at least one address as verified"))
		})
	})
})
