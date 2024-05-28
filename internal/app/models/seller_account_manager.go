package models

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/ss2/internal/app/utils"
	gormio "gorm.io/gorm"
)

type SellerAccountManager struct {
	database.VaccountGorm
	SellerID  uint64
	Role      string
	DeletedAt gormio.DeletedAt
	Priority  int
	Phone     int64
	Name      string
	Email     string
	Seller    *Seller
}

func (sam *SellerAccountManager) Validate(db *gorm.DB) {
	if sam.SellerID == 0 {
		db.AddError(errors.New("SellerID can't be blank")) //nolint:errcheck
	}
	if sam.Role == "" {
		db.AddError(errors.New("Role can't be blank")) //nolint:errcheck
	}
	if sam.Name == "" {
		db.AddError(errors.New("Name can't be blank")) //nolint:errcheck
	}
	if phoneNumber := fmt.Sprint(sam.Phone); !(strings.HasPrefix(phoneNumber, "8801") && len(phoneNumber) == 13) {
		db.AddError(errors.New("Phone Number should have 13 digits")) //nolint:errcheck
	}
}

func GetSellerCodesForSA(ctx context.Context, phone string) ([]uint64, error) {
	sellerCodes := []uint64{}
	intPhone, _ := strconv.ParseInt(phone, 10, 64)
	err := database.DBAPM(ctx).Model(&SellerAccountManager{}).Joins(SamSellerJoinString()).Where(
		&SellerAccountManager{
			Phone: intPhone,
			Role:  utils.SourcingAssociateRole,
		},
	).Pluck("sellers.user_id", &sellerCodes).Error
	return sellerCodes, err
}
