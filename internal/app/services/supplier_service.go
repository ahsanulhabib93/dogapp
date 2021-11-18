package services

import (
	"context"

	supplierpb "github.com/voonik/goConnect/api/go/supplier_service/supplier"
	"github.com/voonik/goFramework/pkg/database"
	"github.com/voonik/supplier_service/internal/app/models"
)

type SupplierService struct{}

// proto
type ListResponse struct {
	Data []models.Supplier
}

type ListParams struct {
	Page    uint64
	PerPage uint64
}

func (ss *SupplierService) List(ctx context.Context, params *supplierpb.ListParams) (*supplierpb.ListResponse, error) {
	resp := supplierpb.ListResponse{}
	database.DBAPM(ctx).Model(&models.Supplier{}).Scan(&resp.Data)
	return &resp, nil
}
