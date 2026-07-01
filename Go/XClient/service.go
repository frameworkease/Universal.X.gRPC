// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XClient

import (
	"reflect"
	"sync"

	"google.golang.org/grpc"
)

var services sync.Map

func Service[T any](factory func(grpc.ClientConnInterface) T) T {
	if factory == nil {
		panic("XClient.Service: factory cannot be nil.")
	}

	hptr := uintptr(reflect.ValueOf(factory).Pointer())
	if service, ok := services.Load(hptr); ok {
		return service.(T)
	}

	typ := reflect.TypeFor[T]()
	configsMu.RLock()
	cfg := configs[typ]
	configsMu.RUnlock()
	interceptorsMu.RLock()
	typedInterceptors := interceptors[typ]
	interceptorsMu.RUnlock()

	var conn *grpc.ClientConn
	if (cfg == nil || (cfg.gateway == "" && len(cfg.options) == 0)) && len(typedInterceptors) == 0 {
		conn = Channel[any]() // 使用共享通道
	} else {
		conn = Channel[T]() // 使用专用通道
	}

	service := factory(conn)
	if actual, loaded := services.LoadOrStore(hptr, service); loaded {
		return actual.(T)
	}
	return service
}
