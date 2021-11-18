package handlers

import (
	"github.com/voonik/supplier_service/internal/app/services"
)

func GetSupplierInstance() *services.SupplierService {
	return new(services.SupplierService)
}
