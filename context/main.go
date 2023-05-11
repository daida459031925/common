package context

import (
	"context"
	"github.com/daida459031925/common/time"
	"net/http"
	"sync"
)

var baseCont baseContext
var once sync.Once // 使用 sync.Once 保证只执行一次初始化操作

const (
	XRequestId requestType = "x-request-id"
	RequestId  requestType = "request_id"
	defTimeOut             = 5
)

type requestType string

type baseContext struct {
	cont context.Context
}

func GetBaseContext() *baseContext {
	once.Do(func() {
		baseCont = baseContext{
			context.Background(),
		}
	})
	return &baseCont
}

// 从 HTTP 标头中提取请求 ID
func (b *baseContext) getRequestIdFromHeader(r *http.Request) string {
	return r.Header.Get(string(XRequestId))
}

// NewSubContext 获取自定义的 Context 对象
func (b *baseContext) NewSubContext(r *http.Request) context.Context {
	return context.WithValue(b.cont, RequestId, b.getRequestIdFromHeader(r))
}

// GetContext 执行超时，取消执行
func (b *baseContext) GetContext(i ...int) (context.Context, context.CancelFunc) {
	return b.GetRequestContext(nil, i...)
}

// GetRequestContext 执行超时，取消执行
func (b *baseContext) GetRequestContext(r *http.Request, i ...int) (context.Context, context.CancelFunc) {
	var t = defTimeOut
	if i != nil && len(i) > 0 && i[0] > 0 {
		t = i[0]
	}

	var cont context.Context

	if r != nil {
		cont = b.NewSubContext(r)
	} else {
		cont = b.cont
	}

	ctx, cancel := context.WithTimeout(cont, time.GetSecond(t))
	return ctx, cancel
}
