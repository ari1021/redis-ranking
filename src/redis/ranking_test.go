package redis_test

import (
	"context"
	"testing"

	"github.com/arkuchy/redis-ranking/src/db"
	"github.com/arkuchy/redis-ranking/src/redis"
)

const limit int = 100

var ctx context.Context

func init() {
	db.InitializeConnection()
	ctx = context.TODO()
}

func BenchmarkGetRankings(b *testing.B) {
	for i := 0; i < b.N; i++ {
		redis.GetRankings(ctx, limit)
	}

}
