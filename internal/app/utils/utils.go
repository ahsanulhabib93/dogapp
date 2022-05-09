package utils

import (
	"context"
	"fmt"
	"strconv"
	"time"

	aaaModels "github.com/voonik/goFramework/pkg/aaa/models"
	"github.com/voonik/goFramework/pkg/cloudstorage"
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

func GetObjectName(fileType string, fileName string, fileExtension string) string {
	if fileExtension == "" {
		fileExtension = "png"
	}
	if fileName == "" {
		fileName = fmt.Sprintf("%s-%s.%s", fileType, strconv.FormatInt(time.Now().UnixNano(), 10), fileExtension)
	}
	return fmt.Sprintf("%s/%s/%s", BucketFolder, fileType, fileName)
}

func GetBucketName(ctx context.Context) string {
	bucketName := aaaModels.AppPreference.GetValue(
		aaaModels.AppPreference{}, ctx, "ss2_bucket", "uatvndrs.shopups2.xyz",
	)
	return bucketName.(string)
}

func GetUploadURL(ctx context.Context, bucketName string, filePath string) (string, error) {
	gcs := cloudstorage.GetGCSClient()
	return gcs.GetUploadURL(ctx, bucketName, filePath, time.Now().Add(45*time.Minute))
}
