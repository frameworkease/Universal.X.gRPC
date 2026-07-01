// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XClient

import (
	"context"
	"net/http"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/frameworkease/Universal.X.gRPC/Go/XServer"
	"github.com/frameworkease/Universal.X.gRPC/Go/internal/Proto/Greet"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

// invokedInterceptors 记录被调用的拦截器。
var invokedInterceptors []string
var invokedInterceptorsMu sync.Mutex

// GreetFooInterceptor 是用于测试的拦截器。
var GreetFooInterceptor = func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	invokedInterceptorsMu.Lock()
	invokedInterceptors = append(invokedInterceptors, "GreetFooInterceptor")
	invokedInterceptorsMu.Unlock()
	return invoker(ctx, method, req, reply, cc, opts...)
}

// GreetBarInterceptor 是用于测试的拦截器。
var GreetBarInterceptor = func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	invokedInterceptorsMu.Lock()
	invokedInterceptors = append(invokedInterceptors, "GreetBarInterceptor")
	invokedInterceptorsMu.Unlock()
	return invoker(ctx, method, req, reply, cc, opts...)
}

// GreetServer 是用于测试的服务器。
type GreetServer struct {
	XServer.Base
}

func (s *GreetServer) Awake() bool {
	XServer.Serve[GreetFooService](&Greet.FooService_ServiceDesc)
	XServer.Serve[GreetBarService](&Greet.BarService_ServiceDesc)
	return true
}

// GreetFooService 是用于测试的 Foo 服务实现。
type GreetFooService struct {
	Greet.UnimplementedFooServiceServer
}

func (s *GreetFooService) SayHi(ctx context.Context, req *Greet.RequestData) (*Greet.ResponseData, error) {
	return &Greet.ResponseData{Message: "Hello " + req.Name + "!"}, nil
}

// GreetBarService 是用于测试的 Bar 服务实现。
type GreetBarService struct {
	Greet.UnimplementedBarServiceServer
}

func (s *GreetBarService) SayHello(ctx context.Context, req *Greet.RequestData) (*Greet.ResponseData, error) {
	return &Greet.ResponseData{Message: "Hi " + req.Name + "!"}, nil
}

func TestService(t *testing.T) {
	// 启动服务
	XServer.Configure(50051, 50052)
	server := &GreetServer{}
	server.Awake()
	server.Start()
	Configure[any]("http://127.0.0.1:50051")

	defer func() {
		// 停止服务
		wait := &sync.WaitGroup{}
		server.Stop(wait)
		wait.Wait()

		// 清理配置和通道
		fooType := reflect.TypeFor[Greet.FooServiceClient]()
		barType := reflect.TypeFor[Greet.BarServiceClient]()
		anyType := reflect.TypeFor[any]()

		channelsMu.Lock()
		if conn, ok := channels[fooType]; ok {
			conn.Close()
			delete(channels, fooType)
		}
		if conn, ok := channels[barType]; ok {
			conn.Close()
			delete(channels, barType)
		}
		if conn, ok := channels[anyType]; ok {
			conn.Close()
			delete(channels, anyType)
		}
		channelsMu.Unlock()

		configsMu.Lock()
		delete(configs, fooType)
		delete(configs, barType)
		delete(configs, anyType)
		configsMu.Unlock()

		// 清理拦截器
		interceptorsMu.Lock()
		delete(interceptors, fooType)
		delete(interceptors, barType)
		delete(interceptors, anyType)
		interceptorsMu.Unlock()

		// 清理实例缓存
		services = sync.Map{}
	}()

	// 等待服务
	{
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		client := &http.Client{Timeout: 2 * time.Second}
		for {
			select {
			case <-ctx.Done():
				t.Fatal("服务超时")
			default:
			}
			resp, err := client.Get("http://127.0.0.1:50052/health")
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				break
			}
			if resp != nil {
				resp.Body.Close()
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	// 测试 FooServiceClient
	{
		// 注册拦截器
		Intercept[Greet.FooServiceClient](GreetFooInterceptor)
		Intercept[Greet.FooServiceClient](GreetBarInterceptor)

		// 测试单例
		client1 := Service(Greet.NewFooServiceClient)
		client2 := Service(Greet.NewFooServiceClient)
		assert.Equal(t, client1, client2, "Foo 服务应当是同一个实例。")

		// 测试服务调用和拦截器
		{
			invokedInterceptorsMu.Lock()
			invokedInterceptors = nil
			invokedInterceptorsMu.Unlock()

			response, err := client1.SayHi(context.Background(), &Greet.RequestData{Name: "Alice"})
			assert.NoError(t, err, "Foo 服务应当正常执行。")
			assert.Equal(t, "Hello Alice!", response.Message, "Foo 服务应当返回正确的消息。")

			invokedInterceptorsMu.Lock()
			assert.Equal(t, "GreetFooInterceptor", invokedInterceptors[0], "Foo 拦截器应当被调用。")
			assert.Equal(t, "GreetBarInterceptor", invokedInterceptors[1], "Bar 拦截器应当被调用。")
			invokedInterceptorsMu.Unlock()
		}
	}

	// 测试 BarServiceClient
	{
		// 注册拦截器
		Intercept[Greet.BarServiceClient](GreetFooInterceptor)
		Intercept[Greet.BarServiceClient](GreetBarInterceptor)

		// 测试单例
		client1 := Service(Greet.NewBarServiceClient)
		client2 := Service(Greet.NewBarServiceClient)
		assert.Equal(t, client1, client2, "Bar 服务应当是同一个实例。")

		// 测试服务调用和拦截器
		{
			invokedInterceptorsMu.Lock()
			invokedInterceptors = nil
			invokedInterceptorsMu.Unlock()

			response, err := client1.SayHello(context.Background(), &Greet.RequestData{Name: "Bob"})
			assert.NoError(t, err, "Bar 服务应当正常执行。")
			assert.Equal(t, "Hi Bob!", response.Message, "Bar 服务应当返回正确的消息。")

			invokedInterceptorsMu.Lock()
			assert.Equal(t, "GreetFooInterceptor", invokedInterceptors[0], "Foo 拦截器应当被调用。")
			assert.Equal(t, "GreetBarInterceptor", invokedInterceptors[1], "Bar 拦截器应当被调用。")
			invokedInterceptorsMu.Unlock()
		}
	}
}
