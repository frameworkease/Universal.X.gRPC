// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XServer

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/frameworkease/Universal.X.gRPC/Go/internal/Proto/Greet"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GreetServer struct {
	Base
	awaked  bool
	started bool
	stopped bool
}

func (s *GreetServer) Awake() bool {
	Serve[GreetFooService](&Greet.FooService_ServiceDesc)
	Serve[GreetBarService](&Greet.BarService_ServiceDesc)
	Intercept(FooInterceptor)
	Intercept(BarInterceptor)
	s.awaked = true
	return true
}

func (s *GreetServer) Start() {
	s.Base.Start()
	s.started = true
}

func (s *GreetServer) Stop(wait *sync.WaitGroup) {
	s.Base.Stop(wait)
	s.stopped = true
}

type GreetFooService struct {
	Greet.UnimplementedFooServiceServer
}

func (s *GreetFooService) SayHi(ctx context.Context, req *Greet.RequestData) (*Greet.ResponseData, error) {
	return &Greet.ResponseData{Message: "Hello " + req.Name + "!"}, nil
}

type GreetBarService struct {
	Greet.UnimplementedBarServiceServer
}

func (s *GreetBarService) SayHello(ctx context.Context, req *Greet.RequestData) (*Greet.ResponseData, error) {
	return &Greet.ResponseData{Message: "Hi " + req.Name + "!"}, nil
}

var invokedInterceptors []string

var FooInterceptor = func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	invokedInterceptors = append(invokedInterceptors, "FooInterceptor")
	return handler(ctx, req)
}

var BarInterceptor = func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	invokedInterceptors = append(invokedInterceptors, "BarInterceptor")
	return handler(ctx, req)
}

func TestLifecycle(t *testing.T) {
	// 准备数据
	Configure(50051, 50052)
	serveUrl := fmt.Sprintf("127.0.0.1:%d", configuredServePort)
	probeUrl := fmt.Sprintf("http://127.0.0.1:%d", configuredProbePort)
	server := &GreetServer{}

	// 唤醒服务
	server.Awake()
	assert.True(t, server.awaked, "Greet 服务应当被唤醒。")

	// 启动服务
	server.Start()
	assert.True(t, server.started, "Greet 服务应当被启动。")

	// 等待服务
	{
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		client := &http.Client{Timeout: 2 * time.Second}
		for {
			select {
			case <-ctx.Done():
				t.Fatal("Greet 服务启动超时")
			default:
			}
			resp, err := client.Get(probeUrl + "/health")
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

	// 验证 healthz 端点
	{
		client := &http.Client{Timeout: 5 * time.Second}
		healthzResponse, err := client.Get(probeUrl + "/healthz")
		assert.NoError(t, err, "healthz 端点应当可以访问。")
		assert.Equal(t, http.StatusOK, healthzResponse.StatusCode, "healthz 端点应当返回 200 OK。")
		healthzResponse.Body.Close()
	}

	// 验证 metrics 端点
	{
		client := &http.Client{Timeout: 5 * time.Second}
		metricsResponse, err := client.Get(probeUrl + "/metrics")
		assert.NoError(t, err, "metrics 端点应当可以访问。")
		assert.Equal(t, http.StatusOK, metricsResponse.StatusCode, "metrics 端点应当返回成功状态码。")
		metricsContent := make([]byte, 1024)
		n, _ := metricsResponse.Body.Read(metricsContent)
		metricsResponse.Body.Close()
		content := string(metricsContent[:n])
		assert.True(t, strings.Contains(content, "go_gc_duration_seconds"), "metrics 端点内容应当包含指标数据。")
	}

	// 验证 pprof 端点
	{
		client := &http.Client{Timeout: 5 * time.Second}
		pprofResponse, err := client.Get(probeUrl + "/pprof/")
		assert.NoError(t, err, "pprof 端点应当可以访问。")
		assert.Equal(t, http.StatusOK, pprofResponse.StatusCode, "pprof 端点应当返回 200 OK。")
		pprofResponse.Body.Close()

		pprofResponse, err = client.Get(probeUrl + "/pprof/allocs")
		assert.NoError(t, err, "allocs 端点应当可以访问。")
		assert.Equal(t, http.StatusOK, pprofResponse.StatusCode, "allocs 端点应当返回 200 OK。")
		pprofResponse.Body.Close()

		pprofResponse, err = client.Get(probeUrl + "/pprof/block")
		assert.NoError(t, err, "block 端点应当可以访问。")
		assert.Equal(t, http.StatusOK, pprofResponse.StatusCode, "block 端点应当返回 200 OK。")
		pprofResponse.Body.Close()

		pprofResponse, err = client.Get(probeUrl + "/pprof/cmdline")
		assert.NoError(t, err, "cmdline 端点应当可以访问。")
		assert.Equal(t, http.StatusOK, pprofResponse.StatusCode, "cmdline 端点应当返回 200 OK。")
		pprofResponse.Body.Close()

		pprofResponse, err = client.Get(probeUrl + "/pprof/goroutine")
		assert.NoError(t, err, "goroutine 端点应当可以访问。")
		assert.Equal(t, http.StatusOK, pprofResponse.StatusCode, "goroutine 端点应当返回 200 OK。")
		pprofResponse.Body.Close()

		pprofResponse, err = client.Get(probeUrl + "/pprof/heap")
		assert.NoError(t, err, "heap 端点应当可以访问。")
		assert.Equal(t, http.StatusOK, pprofResponse.StatusCode, "heap 端点应当返回 200 OK。")
		pprofResponse.Body.Close()

		pprofResponse, err = client.Get(probeUrl + "/pprof/mutex")
		assert.NoError(t, err, "mutex 端点应当可以访问。")
		assert.Equal(t, http.StatusOK, pprofResponse.StatusCode, "mutex 端点应当返回 200 OK。")
		pprofResponse.Body.Close()

		pprofResponse, err = client.Get(probeUrl + "/pprof/profile?seconds=1")
		assert.NoError(t, err, "profile 端点应当可以访问。")
		if assert.NotNil(t, pprofResponse, "profile 响应不应当为 nil") {
			assert.Equal(t, http.StatusOK, pprofResponse.StatusCode, "profile 端点应当返回 200 OK。")
			pprofResponse.Body.Close()
		}

		pprofResponse, err = client.Get(probeUrl + "/pprof/symbol")
		assert.NoError(t, err, "symbol 端点应当可以访问。")
		if assert.NotNil(t, pprofResponse, "symbol 响应不应当为 nil") {
			assert.Equal(t, http.StatusOK, pprofResponse.StatusCode, "symbol 端点应当返回 200 OK。")
			pprofResponse.Body.Close()
		}

		pprofResponse, err = client.Get(probeUrl + "/pprof/threadcreate")
		assert.NoError(t, err, "threadcreate 端点应当可以访问。")
		assert.Equal(t, http.StatusOK, pprofResponse.StatusCode, "threadcreate 端点应当返回 200 OK。")
		pprofResponse.Body.Close()

		pprofResponse, err = client.Get(probeUrl + "/pprof/trace?seconds=1")
		assert.NoError(t, err, "trace 端点应当可以访问。")
		if assert.NotNil(t, pprofResponse, "trace 响应不应当为 nil") {
			assert.Equal(t, http.StatusOK, pprofResponse.StatusCode, "trace 端点应当返回 200 OK。")
			pprofResponse.Body.Close()
		}
	}

	// 验证 grpc 服务
	{
		conn, err := grpc.NewClient(serveUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
		assert.NoError(t, err, "应当能够创建 gRPC 连接。")
		defer conn.Close()

		{
			client := Greet.NewFooServiceClient(conn)
			invokedInterceptors = nil
			response, err := client.SayHi(context.Background(), &Greet.RequestData{Name: "Alice"})
			assert.NoError(t, err, "Foo 服务应当正常执行。")
			assert.Equal(t, "Hello Alice!", response.Message, "Foo 服务应当返回正确的消息。")
			assert.Equal(t, "FooInterceptor", invokedInterceptors[0], "Foo 拦截器应当被调用。")
			assert.Equal(t, "BarInterceptor", invokedInterceptors[1], "Bar 拦截器应当被调用。")
		}

		{
			client := Greet.NewBarServiceClient(conn)
			invokedInterceptors = nil
			response, err := client.SayHello(context.Background(), &Greet.RequestData{Name: "Bob"})
			assert.NoError(t, err, "Bar 服务应当正常执行。")
			assert.Equal(t, "Hi Bob!", response.Message, "Bar 服务应当返回正确的消息。")
			assert.Equal(t, "FooInterceptor", invokedInterceptors[0], "Foo 拦截器应当被调用。")
			assert.Equal(t, "BarInterceptor", invokedInterceptors[1], "Bar 拦截器应当被调用。")
		}
	}

	// 停止服务
	{
		wait := &sync.WaitGroup{}
		server.Stop(wait)
		wait.Wait()
		assert.True(t, server.stopped, "Greet 服务应当被停止。")
	}
}
