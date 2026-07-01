// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XClient

import (
	"reflect"
	"sync"

	"google.golang.org/grpc"
)

type config struct {
	gateway string
	options []grpc.DialOption
}

var (
	configs   = make(map[reflect.Type]*config)
	configsMu sync.RWMutex
)

func Configure[T any](gateway string, options ...grpc.DialOption) {
	typ := reflect.TypeFor[T]()
	configsMu.Lock()
	defer configsMu.Unlock()
	configs[typ] = &config{gateway: gateway, options: options}
}
