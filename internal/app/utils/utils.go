package utils

import (
	"context"

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

func GetBucketName(ctx context.Context) string {
	bucketName := aaaModels.AppPreference.GetValue(
		aaaModels.AppPreference{}, ctx, "ss2_bucket", "uatvndrs.shopups2.xyz",
	)
	return bucketName.(string)
}
