package payment_account_detail_service_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	paywellPb "github.com/voonik/goConnect/api/go/paywell_token/payment_gateway"
	paymentpb "github.com/voonik/goConnect/api/go/ss2/payment_account_detail"
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

var _ = Describe("AddPaymentAccountDetail", func() {
	var ctx context.Context
	var apiHelperInstance *mocks.APIHelperInterface

	BeforeEach(func() {
		test_utils.GetContext(&ctx)
		aaaModels.CreateAppPreferenceServiceInterface()
	})
	AfterEach(func() {
		helpers.InjectMockAPIHelperInstance(nil)
	})

	Context("Add", func() {
		It("Should create payment account detail and return success", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Payment Account Detail Added Successfully"))

			paymentAccounts := []*models.PaymentAccountDetail{{}}
			database.DBAPM(ctx).Model(supplier).Association("PaymentAccountDetails").Find(&paymentAccounts)
			Expect(len(paymentAccounts)).To(Equal(1))
			paymentAccount := paymentAccounts[0]

			Expect(paymentAccount.AccountType).To(Equal(utils.Bank))
			Expect(paymentAccount.AccountSubType).To(Equal(utils.Savings))
			Expect(paymentAccount.AccountName).To(Equal(param.AccountName))
			Expect(paymentAccount.AccountNumber).To(Equal(param.AccountNumber))
			Expect(paymentAccount.BankID).To(Equal(param.BankId))
			Expect(paymentAccount.BranchName).To(Equal(param.BranchName))
			Expect(paymentAccount.RoutingNumber).To(Equal(param.RoutingNumber))
			Expect(paymentAccount.IsDefault).To(Equal(true))

			database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplier.ID)
			Expect(supplier.Status).To(Equal(models.SupplierStatusPending))
		})

		It("Should create payment account detail, prepaid card and return success", func() {
			apiHelperInstance = new(mocks.APIHelperInterface)
			helpers.InjectMockAPIHelperInstance(apiHelperInstance)
			apiHelperInstance.On("CreatePaywellCard", ctx, &paywellPb.CreateCardRequest{UniqueId: "SS2-PAD-1", CardInfo: "11003388", ExpiryMonth: "01", ExpiryYear: "2025"}).Return(&paywellPb.CreateCardResponse{IsError: false, Message: "Successfully created", Token: "sample_token_1", MaskedNumber: "masked_number_11003388"}, nil)
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.PrepaidCard),
				AccountSubType: uint64(utils.UCBL),
				AccountName:    "AccountName",
				AccountNumber:  "11003388",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
				ExtraDetails: &paymentpb.ExtraDetails{
					EmployeeId: uint64(1234),
					ClientId:   uint64(123),
					ExpiryDate: "2025-01-01",
				},
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Payment Account Detail Added Successfully"))

			paymentAccounts := []*models.PaymentAccountDetail{{}}
			database.DBAPM(ctx).Model(supplier).Association("PaymentAccountDetails").Find(&paymentAccounts)
			Expect(len(paymentAccounts)).To(Equal(1))
			paymentAccount := paymentAccounts[0]

			Expect(paymentAccount.AccountType).To(Equal(utils.PrepaidCard))
			Expect(paymentAccount.AccountSubType).To(Equal(utils.UCBL))
			Expect(paymentAccount.AccountName).To(Equal(param.AccountName))
			Expect(paymentAccount.AccountNumber).To(Equal("masked_number_11003388"))
			Expect(paymentAccount.BankID).To(Equal(param.BankId))
			Expect(paymentAccount.BranchName).To(Equal(param.BranchName))
			Expect(paymentAccount.RoutingNumber).To(Equal(param.RoutingNumber))
			Expect(paymentAccount.IsDefault).To(Equal(true))

			extraDetails := models.PaymentAccountDetailExtraDetails{}
			utils.CopyStructAtoB(paymentAccount.ExtraDetails, &extraDetails)
			Expect(extraDetails.ExpiryDate).To(Equal("2025-01-01"))
			Expect(extraDetails.Token).To(Equal("sample_token_1"))
			Expect(extraDetails.ClientId).To(Equal(uint64(123)))
			Expect(extraDetails.EmployeeId).To(Equal(uint64(1234)))
			Expect(extraDetails.UniqueId).To(Equal("SS2-PAD-1"))

			database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplier.ID)
			Expect(supplier.Status).To(Equal(models.SupplierStatusPending))
		})

		It("Should return error if user in blocked state", func() {
			test_utils.SetPermission(&ctx, []string{})

			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusBlocked})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Change Not Allowed"))
		})
	})

	Context("While adding default payment account", func() {
		It("Should return success response and other default payment should be updated as non-default", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{Status: models.SupplierStatusVerified})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: true})
			test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier.ID, IsDefault: false})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				BankId:         bank.ID,
				BranchName:     "branch name",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Payment Account Detail Added Successfully"))

			paymentAccounts := []*models.PaymentAccountDetail{{}}
			database.DBAPM(ctx).Model(supplier).Association("PaymentAccountDetails").Find(&paymentAccounts)
			Expect(len(paymentAccounts)).To(Equal(3))

			paymentAccount1 := paymentAccounts[0]
			Expect(paymentAccount1.IsDefault).To(Equal(false))
			paymentAccount2 := paymentAccounts[1]
			Expect(paymentAccount2.IsDefault).To(Equal(false))
			paymentAccount3 := paymentAccounts[2]
			Expect(paymentAccount3.IsDefault).To(Equal(true))

			database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplier.ID)
			Expect(supplier.Status).To(Equal(models.SupplierStatusPending))
		})
	})

	Context("While adding payment account detail for invalid Supplier ID", func() {
		It("Should return error response", func() {
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:  1000,
				AccountName: "AccountName",
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Supplier Not Found"))
		})
	})

	Context("While adding payment account detail without account name", func() {
		It("Should return error response", func() {
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountNumber:  "AccountNumber",
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Current),
				BankId:         bank.ID,
				BranchName:     "branch name",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Payment Account Detail: account_name can't be blank"))
		})
	})

	Context("While adding payment account detail without account number and account type", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:  supplier.ID,
				AccountName: "AccountName",
				IsDefault:   true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Payment Account Detail: account_number is required; account_type can't be blank"))
		})
	})

	Context("While adding non-default payment account detail first time", func() {
		It("Should return error response", func() {
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Current),
				BankId:         bank.ID,
				BranchName:     "BranchName",
				IsDefault:      false,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Payment Account Detail: Default Payment Account is required"))
		})
	})

	Context("While adding with invalid account subtype", func() {
		It("Should return error response", func() {
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Bkash),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				BankId:         bank.ID,
				BranchName:     "branch name",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Payment Account Detail: Invalid Account SubType"))
		})
	})

	Context("While adding with invalid bank ID", func() {
		It("Should return error response", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				BankId:         1000,
				BranchName:     "BranchName",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Payment Account Detail: Invalid Bank Name"))
		})
	})

	Context("While adding bank type payment account", func() {
		It("Should return error response for empty bank id", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			// bank := test_helper.CreateBank(ctx, &models.Bank{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Payment Account Detail: For Bank account type BankID and BranchName needed"))
		})
		It("Should return error response for empty branch name", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNumber",
				BankId:         bank.ID,
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Payment Account Detail: For Bank account type BankID and BranchName needed"))
		})
	})

	Context("While adding payment account detail with existing account number with AppPreference", func() {
		It("Should return error response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			_ = test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier1.ID, AccountType: utils.Bank, AccountNumber: "AccountNum", IsDefault: true})
			aaaModels.InjectMockAppPreferenceServiceInstance(mocks.GetAppPreferenceMock(map[string]interface{}{
				"enabled_account_number_validation": true,
			}))
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier2.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNum",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
			}

			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Payment Account Detail: Provided bank account number already exists"))
		})
	})

	Context("While adding payment account detail with existing account number without AppPreference", func() {
		It("Should return error response", func() {
			supplier1 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			supplier2 := test_helper.CreateSupplier(ctx, &models.Supplier{})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			_ = test_helper.CreatePaymentAccountDetail(ctx, &models.PaymentAccountDetail{SupplierID: supplier1.ID, AccountType: utils.Bank, AccountNumber: "AccountNum", IsDefault: true})
			aaaModels.InjectMockAppPreferenceServiceInstance(mocks.GetAppPreferenceMock(map[string]interface{}{
				"enabled_account_number_validation": false,
			}))
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier2.ID,
				AccountType:    uint64(utils.Bank),
				AccountSubType: uint64(utils.Savings),
				AccountName:    "AccountName",
				AccountNumber:  "AccountNum",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
			}

			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)

			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Payment Account Detail Added Successfully"))
		})
	})
	Context("Extra Details Validations", func() {
		It("Returns error on older expiry date", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.PrepaidCard),
				AccountSubType: uint64(utils.UCBL),
				AccountName:    "AccountName",
				AccountNumber:  "11003388",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
				ExtraDetails: &paymentpb.ExtraDetails{
					EmployeeId: uint64(1234),
					ClientId:   uint64(123),
					ExpiryDate: "2022-01-01",
				},
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Cannot Create Payment Account: Cannot set older date as expiry date"))

			paymentAccounts := []*models.PaymentAccountDetail{{}}
			database.DBAPM(ctx).Model(supplier).Association("PaymentAccountDetails").Find(&paymentAccounts)
			Expect(len(paymentAccounts)).To(Equal(0))

			database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplier.ID)
			Expect(supplier.Status).To(Equal(models.SupplierStatusPending))
		})

		It("Returns error on invalid expiry date", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.PrepaidCard),
				AccountSubType: uint64(utils.UCBL),
				AccountName:    "AccountName",
				AccountNumber:  "11003388",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
				ExtraDetails: &paymentpb.ExtraDetails{
					EmployeeId: uint64(1234),
					ClientId:   uint64(123),
					ExpiryDate: "ABCD",
				},
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Cannot Create Payment Account: Invalid Date"))

			paymentAccounts := []*models.PaymentAccountDetail{{}}
			database.DBAPM(ctx).Model(supplier).Association("PaymentAccountDetails").Find(&paymentAccounts)
			Expect(len(paymentAccounts)).To(Equal(0))

			database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplier.ID)
			Expect(supplier.Status).To(Equal(models.SupplierStatusPending))
		})

		It("Returns error on invalid expiry date", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.PrepaidCard),
				AccountSubType: uint64(utils.UCBL),
				AccountName:    "AccountName",
				AccountNumber:  "11003388",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
				ExtraDetails: &paymentpb.ExtraDetails{
					ClientId:   uint64(123),
					ExpiryDate: "2025-01-01",
				},
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Cannot Create Payment Account: Employee ID is mandatory"))

			paymentAccounts := []*models.PaymentAccountDetail{{}}
			database.DBAPM(ctx).Model(supplier).Association("PaymentAccountDetails").Find(&paymentAccounts)
			Expect(len(paymentAccounts)).To(Equal(0))

			database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplier.ID)
			Expect(supplier.Status).To(Equal(models.SupplierStatusPending))
		})

		It("Should not add payment account detail, prepaid card for paywell api failure and return failure response", func() {
			apiHelperInstance = new(mocks.APIHelperInterface)
			helpers.InjectMockAPIHelperInstance(apiHelperInstance)
			apiHelperInstance.On("CreatePaywellCard", ctx, &paywellPb.CreateCardRequest{UniqueId: "SS2-PAD-1", CardInfo: "11003388", ExpiryMonth: "01", ExpiryYear: "2025"}).Return(&paywellPb.CreateCardResponse{IsError: true, Message: "Mocked Error Message", Token: "", MaskedNumber: ""}, nil)
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			bank := test_helper.CreateBank(ctx, &models.Bank{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:     supplier.ID,
				AccountType:    uint64(utils.PrepaidCard),
				AccountSubType: uint64(utils.UCBL),
				AccountName:    "AccountName",
				AccountNumber:  "11003388",
				BankId:         bank.ID,
				BranchName:     "BranchName",
				RoutingNumber:  "RoutingNumber",
				IsDefault:      true,
				ExtraDetails: &paymentpb.ExtraDetails{
					EmployeeId: uint64(1234),
					ClientId:   uint64(123),
					ExpiryDate: "2025-01-01",
				},
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Cannot Create Payment Account: Failed to create Paywell Card"))

			paymentAccounts := []*models.PaymentAccountDetail{{}}
			database.DBAPM(ctx).Model(supplier).Association("PaymentAccountDetails").Find(&paymentAccounts)
			Expect(len(paymentAccounts)).To(Equal(0))

			paymentAccount2 := models.PaymentAccountDetail{}
			database.DBAPM(ctx).Model(&models.PaymentAccountDetail{}).Unscoped().Where("supplier_id = ?", supplier.ID).First(&paymentAccount2)
			Expect(paymentAccount2.DeletedAt.Valid).To(Equal(true))
			Expect(paymentAccount2.DeletedAt.Time).NotTo(BeNil())

			database.DBAPM(ctx).Model(&models.Supplier{}).First(&supplier, supplier.ID)
			Expect(supplier.Status).To(Equal(models.SupplierStatusPending))
		})
	})

	Context("Cheque", func() {
		It("Should fail if account name is not passed", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:  supplier.ID,
				AccountType: uint64(utils.Cheque),
				IsDefault:   true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(false))
			Expect(res.Message).To(Equal("Error while creating Payment Account Detail: account_name can't be blank"))
		})

		It("Should add cheque payment account detail", func() {
			supplier := test_helper.CreateSupplier(ctx, &models.Supplier{})
			param := paymentpb.PaymentAccountDetailParam{
				SupplierId:  supplier.ID,
				AccountType: uint64(utils.Cheque),
				AccountName: "Payee Name ABC",
				IsDefault:   true,
			}
			res, err := new(services.PaymentAccountDetailService).Add(ctx, &param)
			Expect(err).To(BeNil())
			Expect(res.Success).To(Equal(true))
			Expect(res.Message).To(Equal("Payment Account Detail Added Successfully"))
		})
	})
})
