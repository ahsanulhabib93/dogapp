package helpers

import (
	"reflect"

	categroyPb "github.com/voonik/goConnect/api/go/cmt/category"
	categoryService "github.com/voonik/goConnect/cmt/category"
)

type CategoryHelper struct{}

var categoryClient CategoryClientInterface

func InjectMockCategoryClientInstance(mockObj CategoryClientInterface) {
	categoryClient = mockObj
}

type CategoryClientInterface interface {
	GetCategoriesData(category_ids []uint64) (*categroyPb.CategoryDataList, error)
}

func getCategoryClient() CategoryClientInterface {
	if categoryClient == nil || reflect.ValueOf(categoryClient).IsNil() {
		return new(CategoryHelper)
	}
	return categoryClient
}

func (s *CategoryHelper) GetCategoriesData(category_ids []uint64) (*categroyPb.CategoryDataList, error) {
	return categoryService.Category().GetCategoryData(&categroyPb.CategoryIDList{CategoryIds: category_ids})
}

func GetParentCategories(category_ids []uint64) []uint64 {
	parent_category := []uint64{}
	resp, _ := getCategoryClient().GetCategoriesData(category_ids)

	for _, cat := range resp {
		parent_category = append(parent_category, cat.Data.Id)
	}
	return parent_category
}
