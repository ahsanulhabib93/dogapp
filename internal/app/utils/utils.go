package utils

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/misc"
)

const (
	DEFAULT_PAGE     = uint64(1)
	DEFAULT_PER_PAGE = uint64(5000)
)

func GetCurrentUserID(ctx context.Context) *uint64 {
	threadUser := misc.ExtractThreadObject(ctx).UserData
	if threadUser != nil && threadUser.GetUserId() != 0 {
		return &threadUser.UserId
	}

	return nil
}

func GetVaccount(ctx context.Context) int64 {
	return misc.ExtractThreadObject(ctx).VaccountId
}

func GetPortalId(ctx context.Context) int64 {
	return misc.ExtractThreadObject(ctx).PortalId
}

func GetCurrentActId(ctx context.Context) int64 {
	return misc.ExtractThreadObject(ctx).CurrentActId
}

func GetXRequestId(ctx context.Context) string {
	return misc.ExtractThreadObject(ctx).XRequestId
}

func GetCurrentUserPermissions(ctx context.Context) []string {
	threadUser := misc.ExtractThreadObject(ctx).UserData
	if threadUser != nil {
		return threadUser.Permissions
	}

	return []string{}
}

func Int64Min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

func Int64Max(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func IsInclude(list []string, value string) bool {
	for _, v := range list {
		if strings.Trim(v, " ") == strings.Trim(value, " ") {
			return true
		}
	}
	return false
}

func SliceDifference(sliceA, sliceB interface{}) (interface{}, error) {
	sliceAReflect := reflect.ValueOf(sliceA)
	sliceBReflect := reflect.ValueOf(sliceB)
	if sliceAReflect.Kind() == reflect.Slice && sliceBReflect.Kind() == reflect.Slice {
		if sliceAReflect.Type() == sliceBReflect.Type() {
			switch sliceA.(type) {
			case []uint64, []string:
				var diff []interface{}
				for i := 0; i < sliceAReflect.Len(); i++ {
					a := sliceAReflect.Index(i).Interface()
					found := false
					for j := 0; j < sliceBReflect.Len(); j++ {
						b := sliceBReflect.Index(j).Interface()
						if isEqual(a, b) {
							found = true
							break
						}
					}
					if !found {
						diff = append(diff, a)
					}
				}
				return typedSlice(diff), nil
			}
		} else {
			return nil, fmt.Errorf("arguments should be of same type")
		}
	} else {
		return nil, fmt.Errorf("arguments should be slices")
	}
	return nil, nil
}

func isEqual(a, b interface{}) bool {
	aReflect := reflect.ValueOf(a)
	bReflect := reflect.ValueOf(b)
	if aReflect.Type() == bReflect.Type() {
		switch a.(type) {
		case uint64:
			return a.(uint64) == b.(uint64)
		case string:
			return a.(string) == b.(string)
		default:
			return false
		}
	} else {
		return false
	}
}

func typedSlice(slice interface{}) interface{} {
	sliceReflect := reflect.ValueOf(slice)
	if sliceReflect.IsValid() && !sliceReflect.IsZero() && sliceReflect.Kind() == reflect.Slice {
		switch sliceReflect.Index(0).Interface().(type) {
		case uint64:
			var result []uint64
			for i := 0; i < sliceReflect.Len(); i++ {
				result = append(result, sliceReflect.Index(i).Interface().(uint64))
			}
			return result
		case string:
			var result []string
			for i := 0; i < sliceReflect.Len(); i++ {
				result = append(result, sliceReflect.Index(i).Interface().(string))
			}
			return result
		default:
			return nil
		}
	} else {
		return nil
	}
}

func GetBucketName(ctx context.Context) string {
	bucketName := aaaModels.AppPreference.GetValue(
		aaaModels.AppPreference{}, ctx, "ss2_bucket", "uatvndrs.shopups2.xyz",
	)
	return bucketName.(string)
}

func IsEmptyStr(s string) bool {
	return strings.TrimSpace(s) == EmptyString
}

func Includes(array interface{}, item interface{}) bool {
	arrRef := reflect.ValueOf(array)
	for i := Zero; i < arrRef.Len(); i++ {
		if arrRef.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

func GetCommonElements(arr1 []string, arr2 []string) []string {
	elementMap := make(map[string]bool)
	commonElements := []string{}

	for _, element := range arr1 {
		elementMap[element] = true
	}

	for _, element := range arr2 {
		if elementMap[element] {
			commonElements = append(commonElements, element)
		}
	}

	return commonElements
}
