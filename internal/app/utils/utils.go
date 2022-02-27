package utils

import (
	"context"

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
