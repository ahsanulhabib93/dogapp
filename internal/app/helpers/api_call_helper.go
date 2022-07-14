package helpers

import (
	"bytes"
	"context"

	"github.com/voonik/goFramework/pkg/rest"
	"google.golang.org/grpc/metadata"
)

type ApiCallHelper struct{}

type ApiCallHelperInterface interface {
	Call(ctx context.Context, url string, headers map[string]string) (*rest.Response, error)
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

func (apiCallHelper *ApiCallHelper) Call(ctx context.Context, url string, headers map[string]string) (*rest.Response, error) {
	reqHeaders, _ := metadata.FromIncomingContext(ctx)
	headers["Authorization"] = reqHeaders["authorization"][0]
	return rest.HTTPRequest(ctx, url, "GET", headers, map[string]interface{}{}, bytes.NewBuffer([]byte{}))
}
