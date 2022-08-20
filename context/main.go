package context

import (
	"context"
	"time"
)

// 执行超时，取消执行
func GetContext(i ...time.Duration) func() (context.Context, context.CancelFunc) {
	return func() (context.Context, context.CancelFunc) {
		var t time.Duration = 5
		if i != nil && len(i) > 0 && i[0] > 0 {
			t = i[0]
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*t)
		return ctx, cancel
	}
}
