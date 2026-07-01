// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XClient

import (
	"reflect"
	"testing"

	"github.com/frameworkease/Universal.X.gRPC/Go/internal/Proto/Greet"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestConfigure(t *testing.T) {
	defer func() {
		anyType := reflect.TypeFor[any]()
		fooType := reflect.TypeFor[Greet.FooServiceClient]()

		configsMu.Lock()
		delete(configs, anyType)
		delete(configs, fooType)
		configsMu.Unlock()
	}()

	gateway := "http://127.0.0.1:50051"
	options := grpc.WithTransportCredentials(insecure.NewCredentials())

	tests := []struct {
		name      string
		configure func(string, grpc.DialOption)
		typ       reflect.Type
	}{
		{
			name:      "Shared",
			configure: func(g string, o grpc.DialOption) { Configure[any](g, o) },
			typ:       reflect.TypeFor[any](),
		},
		{
			name:      "Generic",
			configure: func(g string, o grpc.DialOption) { Configure[Greet.FooServiceClient](g, o) },
			typ:       reflect.TypeFor[Greet.FooServiceClient](),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				configsMu.Lock()
				delete(configs, tt.typ)
				configsMu.Unlock()
			}()
			tt.configure(gateway, options)

			loadedConfig, ok := configs[tt.typ]
			assert.True(t, ok, "配置应当存在。")
			assert.Equal(t, gateway, loadedConfig.gateway, "网关地址应当匹配。")
			assert.Equal(t, options, loadedConfig.options[0], "通道选项应当匹配。")
		})
	}
}
