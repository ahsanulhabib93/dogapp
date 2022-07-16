package helpers

import (
	"bytes"
	"context"

	"github.com/voonik/goFramework/pkg/rest"
)

type ApiCallHelper struct{}

type ApiCallHelperInterface interface {
	Get(ctx context.Context, url string, headers map[string]string) (*rest.Response, error)
}

var apiCall ApiCallHelperInterface

func InjectMockApiCallHelperInstance(mockObj ApiCallHelperInterface) {
	apiCall = mockObj
}

func GetApiCallHelperInstance() ApiCallHelperInterface {
	if apiCall == nil {
		apiCall = new(ApiCallHelper)
	}
	return apiCall
}

func (apiCallHelper *ApiCallHelper) Get(ctx context.Context, url string, headers map[string]string) (*rest.Response, error) {
	return rest.HTTPRequest(ctx, url, "GET", headers, map[string]interface{}{}, bytes.NewBuffer([]byte{}))
}
