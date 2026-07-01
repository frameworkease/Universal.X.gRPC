// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XServer

import (
	"reflect"
	"sync"

	"google.golang.org/grpc"
)

var (
	services   = make([]any, 0)
	servicesMu sync.RWMutex
)

func Serve[T any](desc *grpc.ServiceDesc) bool {
	if desc == nil {
		panic("XServer.Serve: desc cannot be nil.")
	}

	typ := reflect.TypeFor[*T]()
	instance := reflect.New(typ.Elem()).Interface()

	servicesMu.Lock()
	defer servicesMu.Unlock()

	for _, service := range services {
		entry := service.([]any)
		if reflect.TypeOf(entry[0]) == typ {
			return false
		}
	}

	services = append(services, []any{instance, desc})
	return true
}
