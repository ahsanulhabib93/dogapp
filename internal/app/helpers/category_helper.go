package helpers

import (
	"context"
	"log"
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

func (s *CategoryHelper) GetCategoriesData(ctx context.Context, categoryIds []uint64) (*categoryPb.CategoryDataList, error) {
	return categoryService.Category().GetCategoryData(ctx,
		&categoryPb.CategoryIDData{CategoryIds: categoryIds, RootCategoryFilter: true})
}

func GetParentCategories(ctx context.Context, categoryIds []uint64) []uint64 {
	parentCategory := []uint64{}
	resp, err := getCategoryClient().GetCategoriesData(ctx, categoryIds)
	log.Printf("GetParentCategories: resp = %v err = %v\n", resp, err)

	for _, cat := range resp.Data {
		parentCategory = append(parentCategory, cat.Id)
	}
	return parentCategory
}
