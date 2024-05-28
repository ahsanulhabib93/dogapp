package models

func SamSellerJoinString() string {
	return "inner join sellers on sellers.id = seller_account_managers.seller_id and sellers.vaccount_id = seller_account_managers.vaccount_id"
}
