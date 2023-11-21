package vendor_address_service_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	vapb "github.com/voonik/goConnect/api/go/ss2/vendor_address"
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

var _ = Describe("Verify Address", func() {
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
			param := vapb.VerifyAddressParams{}
			res, err := new(services.VendorAddressService).VerifyAddress(ctx, &param)
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("param not specified"))
			Expect(err).To(BeNil())
		})
	})

	Context("When params is passed", func() {
		It("Should update verification status", func() {
			seller := test_helper.CreateSeller(ctx, &models.Seller{})
			vendorAddress1 := test_helper.CreateVendorAddress(ctx, &models.VendorAddress{SellerID: int(seller.ID)})
			param := vapb.VerifyAddressParams{
				Id: vendorAddress1.UUID,
			}
			fmt.Println("vendor ref:", vendorAddress1)
			res, err := new(services.VendorAddressService).VerifyAddress(ctx, &param)
			database.DBAPM(ctx).Model(&models.VendorAddress{}).Where("uuid = ?", vendorAddress1.UUID).Scan(&vendorAddress1)
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("vendor address verified successfully"))
			Expect(vendorAddress1.VerificationStatus).To(Equal(utils.Verified))
			Expect(err).To(BeNil())

			sellerActivityLog := models.SellerActivityLog{}
			err = database.DBAPM(ctx).Model(&models.SellerActivityLog{}).Where("seller_id = ?", seller.UserID).Scan(&sellerActivityLog).Error
			Expect(err).To(BeNil())
			Expect(sellerActivityLog.Action).To(Equal("verify_address"))
		})
	})

	Context("When invalid param is passed", func() {
		It("Should return error", func() {
			param := vapb.VerifyAddressParams{
				Id: "abc",
			}
			res, err := new(services.VendorAddressService).VerifyAddress(ctx, &param)
			Expect(res.Status).To(Equal("failure"))
			Expect(res.Message).To(Equal("vendor address not found"))
			Expect(err).To(BeNil())
		})
	})

	Context("When vendor address is already verified", func() {
		It("Should return verified message", func() {
			seller := test_helper.CreateSeller(ctx, &models.Seller{})
			vendorAddress1 := test_helper.CreateVendorAddress(ctx,
				&models.VendorAddress{
					SellerID:           int(seller.ID),
					VerificationStatus: utils.Verified,
				})
			param := vapb.VerifyAddressParams{
				Id: vendorAddress1.UUID,
			}
			res, err := new(services.VendorAddressService).VerifyAddress(ctx, &param)
			Expect(res.Status).To(Equal("success"))
			Expect(res.Message).To(Equal("vendor address already verified"))
			Expect(err).To(BeNil())
		})
	})
})
