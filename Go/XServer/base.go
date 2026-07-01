// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XServer

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"

	"github.com/frameworkease/Universal.X.Utility/Go/XLog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

type Base struct {
	grpcServer  *grpc.Server
	probeServer *http.Server
}

func (server *Base) Awake() bool { return true }

func (server *Base) Start() {
	if configuredServePort <= 0 {
		panic("XServer.Start: serve port must be greater than 0.")
	}

	interceptorsMu.RLock()
	interceptorChain := make([]grpc.UnaryServerInterceptor, len(interceptors))
	for i, entry := range interceptors {
		interceptorChain[i] = entry.Interceptor
	}
	interceptorsMu.RUnlock()

	server.grpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(interceptorChain...))

	servicesMu.RLock()
	for _, service := range services {
		entry := service.([]any)
		instance := entry[0]
		desc := entry[1].(*grpc.ServiceDesc)
		server.grpcServer.RegisterService(desc, instance)
	}
	servicesMu.RUnlock()

	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", configuredServePort))
	if err != nil {
		panic(fmt.Sprintf("XServer.Start: failed to listen on serve port %d: %v", configuredServePort, err))
	}
	go func() {
		if err := server.grpcServer.Serve(grpcListener); err != nil {
			panic(fmt.Sprintf("XServer.Start: serve grpc server error: %v", err))
		}
	}()

	if configuredProbePort > 0 {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
		mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) { promhttp.Handler().ServeHTTP(w, r) })
		mux.HandleFunc("/pprof/", pprof.Index)
		mux.HandleFunc("/pprof/allocs", func(w http.ResponseWriter, r *http.Request) { pprof.Handler("allocs").ServeHTTP(w, r) })
		mux.HandleFunc("/pprof/block", func(w http.ResponseWriter, r *http.Request) { pprof.Handler("block").ServeHTTP(w, r) })
		mux.HandleFunc("/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/pprof/goroutine", func(w http.ResponseWriter, r *http.Request) { pprof.Handler("goroutine").ServeHTTP(w, r) })
		mux.HandleFunc("/pprof/heap", func(w http.ResponseWriter, r *http.Request) { pprof.Handler("heap").ServeHTTP(w, r) })
		mux.HandleFunc("/pprof/mutex", func(w http.ResponseWriter, r *http.Request) { pprof.Handler("mutex").ServeHTTP(w, r) })
		mux.HandleFunc("/pprof/profile", pprof.Profile)
		mux.HandleFunc("/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/pprof/threadcreate", func(w http.ResponseWriter, r *http.Request) { pprof.Handler("threadcreate").ServeHTTP(w, r) })
		mux.HandleFunc("/pprof/trace", pprof.Trace)

		probeListener, err := net.Listen("tcp", fmt.Sprintf(":%d", configuredProbePort))
		if err != nil {
			panic(fmt.Sprintf("XServer.Start: failed to listen on probe port %d: %v", configuredProbePort, err))
		}

		server.probeServer = &http.Server{Addr: fmt.Sprintf(":%d", configuredProbePort), Handler: mux}
		go func() {
			if err := server.probeServer.Serve(probeListener); err != nil && err != http.ErrServerClosed {
				panic(fmt.Sprintf("XServer.Start: serve probe server error: %v", err))
			}
		}()
	}

	if configuredProbePort <= 0 {
		XLog.Notice("XServer.Start: server has been listening on port %d.", configuredServePort)
	} else {
		XLog.Notice("XServer.Start: server has been listening on port %d and probe on port %d.", configuredServePort, configuredProbePort)
	}
}

func (server *Base) Stop(wait *sync.WaitGroup) {
	if server.grpcServer != nil {
		server.grpcServer.Stop()
	}

	if server.probeServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.probeServer.Shutdown(ctx); err != nil {
			XLog.Error("XServer.Stop: failed to shutdown probe server: %v", err)
		}
	}
}
