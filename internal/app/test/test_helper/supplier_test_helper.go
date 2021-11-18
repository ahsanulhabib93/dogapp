package test_helper

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/supplier_service/internal/app/models"
	"github.com/voonik/supplier_service/internal/app/utils"
)

func CreateSupplier(ctx context.Context, supplier *models.Supplier) *models.Supplier {
	id := rand.Intn(100)
	supplier.Name = fmt.Sprintf("Test-%v", id)
	supplier.Email = fmt.Sprintf("test-%v@shopup.org", id)
	supplier.SupplierType = utils.Hlc
	database.DBAPM(ctx).Save(supplier)
	return supplier
}
