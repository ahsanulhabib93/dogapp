package script

import (
	"context"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/models"
	"log"
	"strconv"
	"strings"
)

//Excel Format
// 0 - name
// 1 - business_name
// 2 - address1
// 3 - phone
// 4 - email
// 5 - bank_name
// 6 - bank_id
// 7 - account_name
// 7 - account_number
// 8 - branch_name
// 9 - routing_number

const (
	supplierFileName  = "internal/app/helpers/scripts/procurement_vendors_29_10_23.xlsx"
	supplierSheetName = "Sheet1"
	serviceType       = 6
	serviceLevel      = 17
)

func AddSuppliersFromExcel(ctx context.Context) {

	f, err := excelize.OpenFile(supplierFileName)
	if err != nil {
		log.Fatalf("failed to open excel file: %v", err)
	}

	rows := f.GetRows(supplierSheetName)
	if err != nil {
		log.Fatalf("failed to get rows from excel file: %v", err)
	}
	phoneNumbers := make(map[string]bool)
	for _, row := range rows[1:] {
		phone := row[3]
		if _, exists := phoneNumbers[phone]; exists {
			log.Printf("Duplicate phone number found: %s", phone)
			return
		}

		phoneNumbers[phone] = true

		row = trimRowSpaces(row)
		if err := validateSupplierData(ctx, row); err != nil {
			log.Fatalf("data validation failed: %v", err)
			return
		}
	}
	log.Printf("Data validation successful")

	tx := database.DBAPM(ctx).Begin()
	for _, row := range rows[1:] {
		row = trimRowSpaces(row)
		supplierId, err := addSupplierToDB(tx, row)
		if err != nil {
			tx.Rollback()
			return
		}

		if err := addSupplierAddressToDB(tx, row, supplierId); err != nil {
			tx.Rollback()
			return
		}

		if err := addPaymentDetailToDB(tx, row, supplierId); err != nil {
			tx.Rollback()
			return
		}

		if err := addPartnerServiceMappingsToDB(tx, row, supplierId); err != nil {
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	log.Printf("Inserted the records successfully")
}

func validateSupplierData(ctx context.Context, row []string) error {
	db := database.DBAPM(ctx)
	var count int64
	db.Model(&models.Supplier{}).Where("phone = ?", row[3]).Count(&count)
	if count > 0 {
		return errors.New(fmt.Sprintf("phone number %s already exists", row[3]))
	}

	bankID := uint64(atoui(row[6]))
	db.Model(&models.Bank{}).Where("id = ?", bankID).Count(&count)
	if count == 0 {
		return errors.New(fmt.Sprintf("provided BankID-%d does not exist", bankID))
	}

	return nil
}

func addSupplierToDB(tx *gorm.DB, row []string) (uint64, error) {
	supplier := models.Supplier{
		Name:            row[0],
		Status:          "Verified",
		Email:           row[4],
		Phone:           row[3],
		BusinessName:    row[1],
		IsPhoneVerified: &[]bool{true}[0],
	}
	if err := tx.Model(&models.Supplier{}).Save(&supplier).Error; err != nil {
		return 0, err
	}
	return supplier.ID, nil
}

func addSupplierAddressToDB(tx *gorm.DB, row []string, supplierId uint64) error {
	address := models.SupplierAddress{
		Firstname:  row[0],
		Address1:   row[2],
		City:       "Dhaka",
		Country:    "Bangladesh",
		Phone:      row[3],
		IsDefault:  true,
		SupplierID: supplierId,
	}
	return tx.Model(&models.SupplierAddress{}).Save(&address).Error
}

func addPaymentDetailToDB(tx *gorm.DB, row []string, supplierId uint64) error {
	payment := models.PaymentAccountDetail{
		AccountType:    1,
		AccountSubType: 2,
		AccountName:    row[7],
		AccountNumber:  row[8],
		BankID:         uint64(atoui(row[6])),
		BranchName:     row[9],
		RoutingNumber:  row[10],
		IsDefault:      true,
		SupplierID:     supplierId,
	}
	return tx.Model(&models.PaymentAccountDetail{}).Save(&payment).Error
}

func addPartnerServiceMappingsToDB(tx *gorm.DB, row []string, supplierId uint64) error {
	partnerServiceMapping := models.PartnerServiceMapping{
		SupplierId:   supplierId,
		ServiceType:  serviceType,
		ServiceLevel: serviceLevel,
		Active:       true,
	}
	return tx.Model(&models.PartnerServiceMapping{}).Save(&partnerServiceMapping).Error
}

func atoui(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func trimRowSpaces(row []string) []string {
	for i, r := range row {
		row[i] = strings.TrimSpace(r)
	}
	return row
}
