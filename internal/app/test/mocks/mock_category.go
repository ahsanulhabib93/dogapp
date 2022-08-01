package mocks

import (
	"context"

	categoryPb "github.com/voonik/goConnect/api/go/cmt/category"
	"github.com/voonik/ss2/internal/app/helpers"

	mock "github.com/stretchr/testify/mock"
)

type MockCategoryHelper struct {
	mock.Mock
}

func SetCategoryMock() *MockCategoryHelper {
	mock := &MockCategoryHelper{}
	helpers.InjectMockCategoryClientInstance(mock)

	return mock
}

func UnsetCategoryMock() {
	helpers.InjectMockCategoryClientInstance(nil)
}

func (_m *MockCategoryHelper) GetCategoriesData(ctx context.Context, category_ids []uint64) (*categoryPb.CategoryDataList, error) {
	args := _m.Called(ctx, category_ids)
	return args.Get(0).(*categoryPb.CategoryDataList), args.Error(1)
}
