package appPreference

import (
	"context"
	"strconv"

	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
)

const shouldSendSupplierLog = "should_send_supplier_log"

func GetAppPreferenceBool(ctx context.Context, key string, defaultValue string) bool {
	value := aaaModels.AppPreference.GetValue(
		aaaModels.AppPreference{}, ctx, key, defaultValue).(string)
	boolVal, err := strconv.ParseBool(value)

	if err != nil {
		v, _ := strconv.ParseBool(defaultValue)
		return v
	}

	return boolVal
}

func ShouldSendSupplierLog(ctx context.Context) bool {
	return GetAppPreferenceBool(ctx, shouldSendSupplierLog, "true")
}
