package supplier_service_test

import (
	"context"
	"sort"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	supplierpb "github.com/voonik/goConnect/api/go/ss2/supplier"
	test_utils "github.com/voonik/goFramework/pkg/unit_test_helper"
	"github.com/voonik/ss2/internal/app/models"
	"github.com/voonik/ss2/internal/app/services"
	"github.com/voonik/ss2/internal/app/test/test_helper"
	"github.com/voonik/ss2/internal/app/utils"
)

var _ = Describe("GetSupplier", func() {
	var ctx context.Context

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
	})

	Context("Get valid supplier", func() {
		Context("PaymentAccountDetails Doen't Have Warehouses Mapped", func() {
			It("Should Respond with success", func() {
				isPhoneVerified := true
				supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
					SupplierType:             utils.Hlc,
					IsPhoneVerified:          &isPhoneVerified,
					NidNumber:                "123456789",
					NidFrontImageUrl:         "abc.com",
					AgreementUrl:             "abc.com",
					SupplierCategoryMappings: []models.SupplierCategoryMapping{{CategoryID: 1}, {CategoryID: 2}},
					SupplierOpcMappings:      []models.SupplierOpcMapping{{ProcessingCenterID: 3}, {ProcessingCenterID: 4}},
					GuarantorImageUrl:        "abc.xyz",
					GuarantorNidNumber:       "0987654321",
				})
				supplierAddress := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID})
				paymentDetails := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})
				resp, err := new(services.SupplierService).Get(ctx, &supplierpb.GetSupplierParam{Id: supplier.ID})

				Expect(err).To(BeNil())
				Expect(resp.Success).To(Equal(true))
				Expect(resp.Data.Email).To(Equal(supplier.Email))
				Expect(resp.Data.Name).To(Equal(supplier.Name))
				Expect(resp.Data.Phone).To(Equal(supplier.Phone))
				Expect(resp.Data.AlternatePhone).To(Equal(supplier.AlternatePhone))
				Expect(resp.Data.BusinessName).To(Equal(supplier.BusinessName))
				Expect(resp.Data.ShopImageUrl).To(Equal(supplier.ShopImageURL))
				Expect(resp.Data.Reason).To(Equal(supplier.Reason))
				Expect(resp.Data.IsPhoneVerified).To(Equal(true))
				Expect(resp.Data.CategoryIds).To(Equal([]uint64{1, 2}))
				Expect(resp.Data.OpcIds).To(Equal([]uint64{3, 4}))
				Expect(resp.Data.SupplierType).To(Equal(uint64(utils.Hlc)))
				Expect(resp.Data.NidNumber).To(Equal(supplier.NidNumber))
				Expect(resp.Data.NidFrontImageUrl).To(Equal(supplier.NidFrontImageUrl))
				Expect(resp.Data.AgreementUrl).To(Equal(supplier.AgreementUrl))
				Expect(resp.Data.GuarantorNidNumber).To(Equal(supplier.GuarantorNidNumber))
				Expect(resp.Data.GuarantorImageUrl).To(Equal(supplier.GuarantorImageUrl))
				Expect(resp.Data.Status).To(Equal(string(models.SupplierStatusPending)))

				Expect(len(resp.Data.SupplierAddresses)).To(Equal(1))
				Expect(resp.Data.SupplierAddresses[0].Firstname).To(Equal(supplierAddress.Firstname))
				Expect(resp.Data.SupplierAddresses[0].Address1).To(Equal(supplierAddress.Address1))
				Expect(resp.Data.SupplierAddresses[0].Phone).To(Equal(supplierAddress.Phone))

				Expect(len(resp.Data.PaymentAccountDetails)).To(Equal(1))
				Expect(resp.Data.PaymentAccountDetails[0].Id).To(Equal(paymentDetails.ID))
				Expect(resp.Data.PaymentAccountDetails[0].AccountName).To(Equal(paymentDetails.AccountName))
				Expect(resp.Data.PaymentAccountDetails[0].AccountNumber).To(Equal(paymentDetails.AccountNumber))
				Expect(resp.Data.PaymentAccountDetails[0].BankId).To(Equal(paymentDetails.BankID))
				Expect(resp.Data.PaymentAccountDetails[0].BranchName).To(Equal(paymentDetails.BranchName))
				Expect(resp.Data.PaymentAccountDetails[0].BankName).To(Equal(resp.Data.PaymentAccountDetails[0].BankName))
				Expect(resp.Data.PaymentAccountDetails[0].RoutingNumber).To(Equal(paymentDetails.RoutingNumber))
				Expect(resp.Data.PaymentAccountDetails[0].IsDefault).To(Equal(paymentDetails.IsDefault))
				Expect(resp.Data.PaymentAccountDetails[0].AccountType).To(Equal(uint64(paymentDetails.AccountType)))
				Expect(resp.Data.PaymentAccountDetails[0].AccountSubType).To(Equal(uint64(paymentDetails.AccountSubType)))
			})

			It("Should Respond with only non-deleted mapping", func() {
				deletedAt := time.Now()
				supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
					SupplierCategoryMappings: []models.SupplierCategoryMapping{
						{CategoryID: 1, DeletedAt: &deletedAt},
						{CategoryID: 2},
					},
					SupplierType: utils.Hlc,
					SupplierOpcMappings: []models.SupplierOpcMapping{
						{ProcessingCenterID: 2},
						{ProcessingCenterID: 3, DeletedAt: &deletedAt},
						{ProcessingCenterID: 4},
					},
				})

				resp, err := new(services.SupplierService).Get(ctx, &supplierpb.GetSupplierParam{Id: supplier.ID})

				Expect(err).To(BeNil())
				Expect(resp.Success).To(Equal(true))
				Expect(resp.Data.Email).To(Equal(supplier.Email))
				Expect(resp.Data.Name).To(Equal(supplier.Name))
				Expect(resp.Data.CategoryIds).To(Equal([]uint64{2}))
				Expect(resp.Data.OpcIds).To(Equal([]uint64{2, 4}))
				Expect(resp.Data.SupplierType).To(Equal(uint64(utils.Hlc)))
			})
		})
		Context("PaymentAccountDetails Have Warehouses Mapped", func() {
			It("Should Return only Given Warehouse Mapped PaymentAccountDetails", func() {
				isPhoneVerified := true
				supplier := test_helper.CreateSupplier(ctx, &models.Supplier{
					SupplierType:             utils.Hlc,
					IsPhoneVerified:          &isPhoneVerified,
					NidNumber:                "123456789",
					NidFrontImageUrl:         "abc.com",
					AgreementUrl:             "abc.com",
					SupplierCategoryMappings: []models.SupplierCategoryMapping{{CategoryID: 1}, {CategoryID: 2}},
					SupplierOpcMappings:      []models.SupplierOpcMapping{{ProcessingCenterID: 3}, {ProcessingCenterID: 4}},
					GuarantorImageUrl:        "abc.xyz",
					GuarantorNidNumber:       "0987654321",
				})
				supplierAddress := test_helper.CreateSupplierAddress(ctx, &models.SupplierAddress{SupplierID: supplier.ID})
				paymentDetail1 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})
				paymentDetail2 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: false})
				paymentDetail3 := test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: false})
				test_helper.CreatePaymentAccountDetailWarehouseMappings(ctx, &models.PaymentAccountDetailWarehouseMapping{WarehouseID: 10, PaymentAccountDetailID: paymentDetail1.ID})
				test_helper.CreatePaymentAccountDetailWarehouseMappings(ctx, &models.PaymentAccountDetailWarehouseMapping{WarehouseID: 10, PaymentAccountDetailID: paymentDetail2.ID})
				test_helper.CreatePaymentAccountDetailWarehouseMappings(ctx, &models.PaymentAccountDetailWarehouseMapping{WarehouseID: 11, PaymentAccountDetailID: paymentDetail3.ID})

				resp, err := new(services.SupplierService).Get(ctx, &supplierpb.GetSupplierParam{
					Id:          supplier.ID,
					WarehouseId: 10,
				})

				Expect(err).To(BeNil())
				Expect(resp.Success).To(Equal(true))
				Expect(resp.Data.Email).To(Equal(supplier.Email))
				Expect(resp.Data.Name).To(Equal(supplier.Name))
				Expect(resp.Data.Phone).To(Equal(supplier.Phone))
				Expect(resp.Data.AlternatePhone).To(Equal(supplier.AlternatePhone))
				Expect(resp.Data.BusinessName).To(Equal(supplier.BusinessName))
				Expect(resp.Data.ShopImageUrl).To(Equal(supplier.ShopImageURL))
				Expect(resp.Data.Reason).To(Equal(supplier.Reason))
				Expect(resp.Data.IsPhoneVerified).To(Equal(true))
				Expect(resp.Data.CategoryIds).To(Equal([]uint64{1, 2}))
				Expect(resp.Data.OpcIds).To(Equal([]uint64{3, 4}))
				Expect(resp.Data.SupplierType).To(Equal(uint64(utils.Hlc)))
				Expect(resp.Data.NidNumber).To(Equal(supplier.NidNumber))
				Expect(resp.Data.NidFrontImageUrl).To(Equal(supplier.NidFrontImageUrl))
				Expect(resp.Data.AgreementUrl).To(Equal(supplier.AgreementUrl))
				Expect(resp.Data.GuarantorNidNumber).To(Equal(supplier.GuarantorNidNumber))
				Expect(resp.Data.GuarantorImageUrl).To(Equal(supplier.GuarantorImageUrl))
				Expect(resp.Data.Status).To(Equal(string(models.SupplierStatusPending)))

				Expect(len(resp.Data.SupplierAddresses)).To(Equal(1))
				Expect(resp.Data.SupplierAddresses[0].Firstname).To(Equal(supplierAddress.Firstname))
				Expect(resp.Data.SupplierAddresses[0].Address1).To(Equal(supplierAddress.Address1))
				Expect(resp.Data.SupplierAddresses[0].Phone).To(Equal(supplierAddress.Phone))

				Expect(len(resp.Data.PaymentAccountDetails)).To(Equal(2))
				// sorting in asc order
				sort.Slice(resp.Data.PaymentAccountDetails, func(i, j int) bool {
					return resp.Data.PaymentAccountDetails[i].Id < resp.Data.PaymentAccountDetails[j].Id
				})
				Expect(resp.Data.PaymentAccountDetails[0].Id).To(Equal(paymentDetail1.ID))
				Expect(resp.Data.PaymentAccountDetails[0].AccountName).To(Equal(paymentDetail1.AccountName))
				Expect(resp.Data.PaymentAccountDetails[0].AccountNumber).To(Equal(paymentDetail1.AccountNumber))
				Expect(resp.Data.PaymentAccountDetails[0].BankId).To(Equal(paymentDetail1.BankID))
				Expect(resp.Data.PaymentAccountDetails[0].BranchName).To(Equal(paymentDetail1.BranchName))
				Expect(resp.Data.PaymentAccountDetails[0].BankName).To(Equal(resp.Data.PaymentAccountDetails[0].BankName))
				Expect(resp.Data.PaymentAccountDetails[0].RoutingNumber).To(Equal(paymentDetail1.RoutingNumber))
				Expect(resp.Data.PaymentAccountDetails[0].IsDefault).To(Equal(paymentDetail1.IsDefault))
				Expect(resp.Data.PaymentAccountDetails[0].AccountType).To(Equal(uint64(paymentDetail1.AccountType)))
				Expect(resp.Data.PaymentAccountDetails[0].AccountSubType).To(Equal(uint64(paymentDetail1.AccountSubType)))

				Expect(resp.Data.PaymentAccountDetails[1].Id).To(Equal(paymentDetail2.ID))
				Expect(resp.Data.PaymentAccountDetails[1].AccountName).To(Equal(paymentDetail2.AccountName))
				Expect(resp.Data.PaymentAccountDetails[1].AccountNumber).To(Equal(paymentDetail2.AccountNumber))
				Expect(resp.Data.PaymentAccountDetails[1].BankId).To(Equal(paymentDetail2.BankID))
				Expect(resp.Data.PaymentAccountDetails[1].BranchName).To(Equal(paymentDetail2.BranchName))
				Expect(resp.Data.PaymentAccountDetails[1].BankName).To(Equal(resp.Data.PaymentAccountDetails[1].BankName))
				Expect(resp.Data.PaymentAccountDetails[1].RoutingNumber).To(Equal(paymentDetail2.RoutingNumber))
				Expect(resp.Data.PaymentAccountDetails[1].IsDefault).To(Equal(paymentDetail2.IsDefault))
				Expect(resp.Data.PaymentAccountDetails[1].AccountType).To(Equal(uint64(paymentDetail2.AccountType)))
				Expect(resp.Data.PaymentAccountDetails[1].AccountSubType).To(Equal(uint64(paymentDetail2.AccountSubType)))
			})
		})
	})

	Context("Invalid supplier", func() {
		It("Should Respond with error", func() {
			res, err := new(services.SupplierService).Get(ctx, &supplierpb.GetSupplierParam{Id: 10})
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
		})
	})
})
