package helpers

import (
	"context"
	"reflect"

	categoryPb "github.com/voonik/goConnect/api/go/cmt/category"
	categoryService "github.com/voonik/goConnect/cmt/category"
)

type CategoryHelper struct{}

var categoryClient CategoryClientInterface

func InjectMockCategoryClientInstance(mockObj CategoryClientInterface) {
	categoryClient = mockObj
}

type CategoryClientInterface interface {
	GetCategoriesData(ctx context.Context, category_ids []uint64) (*categoryPb.CategoryDataList, error)
}

func getCategoryClient() CategoryClientInterface {
	if categoryClient == nil || reflect.ValueOf(categoryClient).IsNil() {
		return new(CategoryHelper)
	}
	return categoryClient
}

func (s *CategoryHelper) GetCategoriesData(ctx context.Context, category_ids []uint64) (*categoryPb.CategoryDataList, error) {
	return categoryService.Category().GetCategoryData(ctx,
		&categoryPb.CategoryIDData{CategoryIds: category_ids, RootCategoryFilter: true})
}

func GetParentCategories(ctx context.Context, category_ids []uint64) []uint64 {
	parent_category := []uint64{}
	resp, _ := getCategoryClient().GetCategoriesData(ctx, category_ids)

	for _, cat := range resp.Data {
		parent_category = append(parent_category, cat.Id)
	}
	return parent_category
}
