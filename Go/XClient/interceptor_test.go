// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XClient

import (
	"context"
	"reflect"
	"testing"

	"github.com/frameworkease/Universal.X.gRPC/Go/internal/Proto/Greet"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestIntercept(t *testing.T) {
	// fooInterceptor 是用于测试的拦截器。
	fooInterceptor := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	// barInterceptor 是用于测试的拦截器。
	barInterceptor := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	// resetInterceptors 清理测试的拦截器。
	resetInterceptors := func() {
		interceptorsMu.Lock()
		defer interceptorsMu.Unlock()

		fooInterceptorPtr := reflect.ValueOf(fooInterceptor).Pointer()
		barInterceptorPtr := reflect.ValueOf(barInterceptor).Pointer()

		for typ, list := range interceptors {
			filtered := make(interceptorList, 0)
			for _, entry := range list {
				entryPtr := reflect.ValueOf(entry.Interceptor).Pointer()
				if entryPtr != fooInterceptorPtr && entryPtr != barInterceptorPtr {
					filtered = append(filtered, entry)
				}
			}
			if len(filtered) == 0 {
				delete(interceptors, typ)
			} else {
				interceptors[typ] = filtered
			}
		}
	}

	// getInterceptors 获取指定类型的拦截器列表。
	getInterceptors := func(typ reflect.Type) interceptorList {
		interceptorsMu.RLock()
		defer interceptorsMu.RUnlock()
		// 返回一个副本以避免竞态条件
		result := interceptors[typ]
		return result
	}

	// indexInterceptor 查找拦截器在列表中的索引。
	indexInterceptor := func(list interceptorList, interceptor grpc.UnaryClientInterceptor) int {
		interceptorPtr := reflect.ValueOf(interceptor).Pointer()
		for i, entry := range list {
			if reflect.ValueOf(entry.Interceptor).Pointer() == interceptorPtr {
				return i
			}
		}
		return -1
	}

	tests := []struct {
		name     string
		priority bool
	}{
		{
			name:     "默认优先级",
			priority: false,
		},
		{
			name:     "指定优先级",
			priority: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer resetInterceptors()

			fooPriority := 0
			barPriority := 0
			if tt.priority {
				fooPriority = 1000
				barPriority = 1001
			}

			assert.True(t, Intercept[any](fooInterceptor, fooPriority), "Foo 拦截器应当注册至共享拦截器列表。")
			assert.False(t, Intercept[any](fooInterceptor, fooPriority), "Foo 拦截器只能注册一次至共享拦截器列表。")
			assert.False(t, Intercept[Greet.FooServiceClient](fooInterceptor, fooPriority), "Foo 拦截器只能注册一次至 FooServiceClient 专用拦截器列表。")
			assert.True(t, Intercept[Greet.FooServiceClient](barInterceptor, barPriority), "Bar 拦截器应当注册至 FooServiceClient 专用拦截器列表。")
			assert.True(t, Intercept[Greet.BarServiceClient](barInterceptor, barPriority), "Bar 拦截器应当注册至 BarServiceClient 专用拦截器列表。")
			assert.True(t, Intercept[any](barInterceptor, barPriority), "Bar 拦截器应当注册至共享拦截器列表。")

			anyType := reflect.TypeOf((*any)(nil)).Elem()
			fooServiceType := reflect.TypeOf((*Greet.FooServiceClient)(nil)).Elem()
			barServiceType := reflect.TypeOf((*Greet.BarServiceClient)(nil)).Elem()

			anyInterceptors := getInterceptors(anyType)
			fooInterceptors := getInterceptors(fooServiceType)
			barInterceptors := getInterceptors(barServiceType)

			assert.True(t, indexInterceptor(anyInterceptors, fooInterceptor) >= 0, "共享拦截器列表应当包含 Foo 拦截器。")
			assert.True(t, indexInterceptor(anyInterceptors, barInterceptor) >= 0, "共享拦截器列表应当包含 Bar 拦截器。")
			assert.Less(t, indexInterceptor(anyInterceptors, fooInterceptor), indexInterceptor(anyInterceptors, barInterceptor), "Foo 拦截器应当排在 Bar 拦截器之前。")
			assert.True(t, indexInterceptor(fooInterceptors, barInterceptor) >= 0, "FooServiceClient 专用拦截器列表应当包含 Bar 拦截器。")
			assert.True(t, indexInterceptor(barInterceptors, barInterceptor) >= 0, "BarServiceClient 专用拦截器列表应当包含 Bar 拦截器。")
		})
	}
}
