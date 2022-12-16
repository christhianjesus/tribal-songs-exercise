package services

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"
)

type CacheClient interface {
	Get(context.Context, string) *redis.StringCmd
	Set(context.Context, string, interface{}, time.Duration) *redis.StatusCmd
}
