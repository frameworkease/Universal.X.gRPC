// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XServer

import (
	"reflect"
	"testing"

	"github.com/frameworkease/Universal.X.gRPC/Go/internal/Proto/Greet"
	"github.com/stretchr/testify/assert"
)

func TestServe(t *testing.T) {
	type GreetFooService struct {
		Greet.UnimplementedFooServiceServer
	}

	type GreetBarService struct {
		Greet.UnimplementedBarServiceServer
	}

	defer func() {
		servicesMu.Lock()
		for i := len(services) - 1; i >= 0; i-- {
			entry := services[i].([]any)
			instance := entry[0]
			typ := reflect.ValueOf(instance).Type()
			if typ == reflect.TypeFor[*GreetFooService]() || typ == reflect.TypeOf((*GreetBarService)(nil)) {
				services = append(services[:i], services[i+1:]...)
			}
		}
		servicesMu.Unlock()
	}()

	assert.True(t, Serve[GreetFooService](&Greet.FooService_ServiceDesc), "Foo 服务应当注册成功。")
	assert.True(t, Serve[GreetBarService](&Greet.BarService_ServiceDesc), "Bar 服务应当注册成功。")
	assert.False(t, Serve[GreetFooService](&Greet.FooService_ServiceDesc), "Foo 服务只能注册一次。")
	assert.False(t, Serve[GreetBarService](&Greet.BarService_ServiceDesc), "Bar 服务只能注册一次。")

	servicesMu.RLock()
	hasFooService := false
	hasBarService := false
	for _, service := range services {
		entry := service.([]any)
		instance := entry[0]
		typ := reflect.ValueOf(instance).Type()
		if typ == reflect.TypeFor[*GreetFooService]() {
			hasFooService = true
		}
		if typ == reflect.TypeFor[*GreetBarService]() {
			hasBarService = true
		}
	}
	servicesMu.RUnlock()

	assert.True(t, hasFooService, "服务列表应当包含 Foo 服务。")
	assert.True(t, hasBarService, "服务列表应当包含 Bar 服务。")
}
