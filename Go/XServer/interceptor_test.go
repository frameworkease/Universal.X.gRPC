// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XServer

import (
	"context"
	"reflect"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestIntercept(t *testing.T) {
	tests := []struct {
		name      string
		priority1 int
		priority2 int
	}{
		{
			name:      "默认优先级",
			priority1: 0,
			priority2: 0,
		},
		{
			name:      "指定优先级",
			priority1: 1000,
			priority2: 1001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fooInterceptor := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
				return handler(ctx, req)
			}
			barInterceptor := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
				return handler(ctx, req)
			}
			fooInterceptorPointer := reflect.ValueOf(fooInterceptor).Pointer()
			barInterceptorPointer := reflect.ValueOf(barInterceptor).Pointer()

			defer func() {
				interceptorsMu.Lock()
				interceptors = slices.DeleteFunc(interceptors, func(ele *interceptorEntry) bool {
					p := reflect.ValueOf(ele.Interceptor).Pointer()
					return p == fooInterceptorPointer || p == barInterceptorPointer
				})
				interceptorsMu.Unlock()
			}()

			assert.True(t, Intercept(fooInterceptor, tt.priority1), "Foo 拦截器应当注册成功。")
			assert.True(t, Intercept(barInterceptor, tt.priority2), "Bar 拦截器应当注册成功。")
			assert.False(t, Intercept(fooInterceptor), "Foo 拦截器只能注册一次。")
			assert.False(t, Intercept(barInterceptor), "Bar 拦截器只能注册一次。")

			interceptorsMu.RLock()
			hasFooInterceptor := false
			hasBarInterceptor := false
			fooInterceptorIndex := -1
			barInterceptorIndex := -1
			for i, entry := range interceptors {
				interceptorPointer := reflect.ValueOf(entry.Interceptor).Pointer()
				if interceptorPointer == fooInterceptorPointer {
					hasFooInterceptor = true
					fooInterceptorIndex = i
				}
				if interceptorPointer == barInterceptorPointer {
					hasBarInterceptor = true
					barInterceptorIndex = i
				}
			}
			interceptorsMu.RUnlock()

			assert.True(t, hasFooInterceptor, "拦截器列表应当包含 Foo 拦截器。")
			assert.True(t, hasBarInterceptor, "拦截器列表应当包含 Bar 拦截器。")
			assert.Less(t, fooInterceptorIndex, barInterceptorIndex, "Foo 拦截器应当排在 Bar 拦截器之前。")
		})
	}
}
