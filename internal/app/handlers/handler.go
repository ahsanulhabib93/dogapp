package handlers

import (
	"github.com/voonik/ss2/internal/app/services"
)

func GetSupplierInstance() *services.SupplierService {
	return new(services.SupplierService)
}
