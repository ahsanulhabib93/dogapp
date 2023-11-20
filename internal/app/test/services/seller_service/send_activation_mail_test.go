package seller_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	spb "github.com/voonik/goConnect/api/go/ss2/seller"
	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/goFramework/pkg/misc"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/helpers"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/mocks"
	"github.com/voonik/ss2/internal/app/utils"
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
		It("Should return status failure for invalid param", func() {
			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("Seller Ids and Action Should be Present"))
		})

		It("Should return status failure for seller not found", func() {
			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate"}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("Seller not found"))
		})

		It("Should return status failure for seller bank detail not found", func() {
			seller = models.Seller{UserID: 1}
			database.DBAPM(ctx).Create(&seller)

			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate"}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("1: Seller Pan Number, Bank Detail, MOU and Email should be confirmed"))
		})

		It("Should return status failure for > 1 vendor address", func() {
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
			Expect(res.Message).To(Equal("1: At least one address should be present Need ActivationState, StateReason to find Non Access Sellers"))
		})

		It("Should return status failure for no verify status address", func() {
			seller = models.Seller{
				UserID:          1,
				PanNumber:       "PAN123",
				EmailConfirmed:  true,
				MouAgreed:       true,
				ActivationState: 3,
				StateReason:     4,
			}
			database.DBAPM(ctx).Create(&seller)

			sellerBankDetail := models.SellerBankDetail{SellerID: int(seller.ID)}
			database.DBAPM(ctx).Create(&sellerBankDetail)

			vendorAddress := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED", VerificationStatus: "VERIFIED"}
			database.DBAPM(ctx).Create(&vendorAddress)

			vendorAddress2 := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "Status1"}
			database.DBAPM(ctx).Create(&vendorAddress2)

			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate"}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("1: Make at least one address as default"))
		})

		It("Should return status failure for no default status address", func() {
			seller = models.Seller{
				UserID:          1,
				PanNumber:       "PAN123",
				EmailConfirmed:  true,
				MouAgreed:       true,
				ActivationState: 3,
				StateReason:     4,
			}
			database.DBAPM(ctx).Create(&seller)

			sellerBankDetail := models.SellerBankDetail{SellerID: int(seller.ID)}
			database.DBAPM(ctx).Create(&sellerBankDetail)

			vendorAddress := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED"}
			database.DBAPM(ctx).Create(&vendorAddress)

			vendorAddress2 := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "Status2", DefaultAddress: true}
			database.DBAPM(ctx).Create(&vendorAddress2)

			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate"}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("1: Make at least one address as verified"))
		})

		It("Should return status failure for seller price details not present", func() {
			seller = models.Seller{
				UserID:          1,
				PanNumber:       "PAN123",
				EmailConfirmed:  true,
				MouAgreed:       true,
				ActivationState: 3,
				StateReason:     4,
			}
			database.DBAPM(ctx).Create(&seller)

			sellerBankDetail := models.SellerBankDetail{SellerID: int(seller.ID)}
			database.DBAPM(ctx).Create(&sellerBankDetail)

			vendorAddress := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED"}
			database.DBAPM(ctx).Create(&vendorAddress)

			vendorAddress2 := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED", VerificationStatus: "VERIFIED"}
			database.DBAPM(ctx).Create(&vendorAddress2)

			vendorAddress3 := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED", DefaultAddress: true}
			database.DBAPM(ctx).Create(&vendorAddress3)

			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate"}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("1: Seller pricing details are not present"))
		})

		It("Should return status failure for seller price details not present", func() {
			seller = models.Seller{
				UserID:          1,
				PanNumber:       "PAN123",
				EmailConfirmed:  true,
				MouAgreed:       true,
				ActivationState: 3,
				StateReason:     4,
			}
			seller.SellerPricingDetails = []*models.SellerPricingDetail{{Verified: utils.SellerPriceVerified(utils.NotVerified),
				SellerID: int(seller.ID)}}
			database.DBAPM(ctx).Create(&seller)

			sellerBankDetail := models.SellerBankDetail{SellerID: int(seller.ID)}
			database.DBAPM(ctx).Create(&sellerBankDetail)

			vendorAddress := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED"}
			database.DBAPM(ctx).Create(&vendorAddress)

			vendorAddress2 := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED", VerificationStatus: "VERIFIED"}
			database.DBAPM(ctx).Create(&vendorAddress2)

			vendorAddress3 := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED", DefaultAddress: true}
			database.DBAPM(ctx).Create(&vendorAddress3)

			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate"}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("1: Seller pricing details are not verified"))
		})
	})

	Context("success case", func() {
		BeforeEach(func() {
			ctx = misc.SetInContextThreadObject(ctx, &misc.ThreadObject{VaccountId: 100, PortalId: 100,
				UserData: &misc.UserData{UserId: 11}})

			seller = models.Seller{
				UserID:          2,
				PanNumber:       "PAN123",
				EmailConfirmed:  true,
				MouAgreed:       true,
				ActivationState: 3,
				StateReason:     4,
			}
			seller.SellerPricingDetails = []*models.SellerPricingDetail{{Verified: utils.SellerPriceVerified(utils.Verified),
				SellerID: int(seller.ID)}}
			database.DBAPM(ctx).Create(&seller)

			sellerBankDetail := models.SellerBankDetail{SellerID: int(seller.ID)}
			database.DBAPM(ctx).Create(&sellerBankDetail)

			vendorAddress := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED"}
			database.DBAPM(ctx).Create(&vendorAddress)

			vendorAddress2 := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED", VerificationStatus: "VERIFIED"}
			database.DBAPM(ctx).Create(&vendorAddress2)

			vendorAddress3 := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED", DefaultAddress: true}
			database.DBAPM(ctx).Create(&vendorAddress3)
		})

		It("Should return status success for seller onboarding team", func() {

			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate", IsSellerOnboardingTeam: true}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("Seller account activated successfully You don't have access to activate this Seller(s) - 1"))
		})

		It("Should return status success for quality team", func() {
			seller.StateReason = 2
			seller.BusinessType = utils.Manufacturer
			seller.ColorCode = utils.Black
			database.DBAPM(ctx).Save(&seller)

			//for coverage to handle IsQualityTeam
			param2 := spb.SendActivationMailParams{Ids: []uint64{2}, Action: "activate", IsQualityTeam: true}
			res2, err2 := new(services.SellerService).SendActivationMail(ctx, &param2)
			Expect(err2).To(BeNil())
			Expect(res2.Status).To(Equal("success"))
			Expect(res2.Message).To(Equal("Seller account activated successfully You don't have access to activate this Seller(s) - 1"))
		})
	})

	Context("success case", func() {
		BeforeEach(func() {
			ctx = misc.SetInContextThreadObject(ctx, &misc.ThreadObject{VaccountId: 100, PortalId: 100,
				UserData: &misc.UserData{UserId: 11}})

			seller = models.Seller{
				UserID:          2,
				PanNumber:       "PAN123",
				EmailConfirmed:  true,
				MouAgreed:       true,
				ActivationState: 1,
				StateReason:     4,
			}
			seller.SellerPricingDetails = []*models.SellerPricingDetail{{Verified: utils.SellerPriceVerified(utils.Verified),
				SellerID: int(seller.ID)}}
			database.DBAPM(ctx).Create(&seller)

			sellerBankDetail := models.SellerBankDetail{SellerID: int(seller.ID)}
			database.DBAPM(ctx).Create(&sellerBankDetail)

			vendorAddress := models.VendorAddress{SellerID: int(seller.ID), GSTStatus: "VERIFIED", VerificationStatus: "VERIFIED", DefaultAddress: true}
			database.DBAPM(ctx).Create(&vendorAddress)
		})

		It("Should return status success and update vendor address verification status", func() {
			param := spb.SendActivationMailParams{Ids: []uint64{1, 2, 3}, Action: "activate", IsSellerOnboardingTeam: true}
			res, err := new(services.SellerService).SendActivationMail(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("Seller account activated successfully"))
		})
	})
})
