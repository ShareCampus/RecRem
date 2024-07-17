package middlewares

import (
	"net/http"
	"recrem/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

var bucket *ratelimit.Bucket

// InitBucket 初始化 Bucket
func InitBucket(fillInternal time.Duration, capacity int64) {
	bucket = ratelimit.NewBucket(fillInternal, capacity)
}

// Limiter 限流中间件
func Limiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			ctx.JSON(http.StatusForbidden, utils.Result{
				Msg:  "您的访问过于频繁！",
				Data: nil,
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
