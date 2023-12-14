// Code generated by "stringer -type=SupplierType"; DO NOT EDIT.

package utils

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[L0-1]
	_ = x[L1-2]
	_ = x[L2-3]
	_ = x[L3-4]
	_ = x[Hlc-5]
	_ = x[Captive-6]
	_ = x[Driver-7]
	_ = x[CashVendor-8]
	_ = x[RedxHubVendor-9]
	_ = x[CreditVendor-10]
	_ = x[HubRent-11]
	_ = x[WarehouseRent-12]
	_ = x[DBHouseRent-13]
	_ = x[OfficeRent-14]
	_ = x[Mws-15]
	_ = x[Buyer-16]
	_ = x[Procurement-17]
	_ = x[InternalEmployee-18]
}

const _SupplierType_name = "L0L1L2L3HlcCaptiveDriverCashVendorRedxHubVendorCreditVendorHubRentWarehouseRentDBHouseRentOfficeRentMwsBuyerProcurementInternalEmployee"

var _SupplierType_index = [...]uint8{0, 2, 4, 6, 8, 11, 18, 24, 34, 47, 59, 66, 79, 90, 100, 103, 108, 119, 135}

func (i SupplierType) String() string {
	i -= 1
	if i >= SupplierType(len(_SupplierType_index)-1) {
		return "SupplierType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _SupplierType_name[_SupplierType_index[i]:_SupplierType_index[i+1]]
}
