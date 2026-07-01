// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XClient

import (
	"reflect"
	"testing"

	"github.com/frameworkease/Universal.X.gRPC/Go/internal/Proto/Greet"
	"github.com/stretchr/testify/assert"
)

func TestChannel(t *testing.T) {
	defer func() {
		anyType := reflect.TypeFor[any]()
		fooType := reflect.TypeFor[Greet.FooServiceClient]()

		channelsMu.Lock()
		if conn, ok := channels[anyType]; ok {
			conn.Close()
			delete(channels, anyType)
		}
		if conn, ok := channels[fooType]; ok {
			conn.Close()
			delete(channels, fooType)
		}
		channelsMu.Unlock()
	}()

	t.Run("Shared", func(t *testing.T) {
		Configure[any]("http://127.0.0.1:50051")
		channel1 := Channel[any]()
		channel2 := Channel[any]()

		assert.NotNil(t, channel1, "共享 Channel 应当被创建。")
		assert.Equal(t, channel1, channel2, "共享 Channel 应当是单例的。")
	})

	t.Run("Generic", func(t *testing.T) {
		Configure[Greet.FooServiceClient]("http://127.0.0.1:50051")
		channel1 := Channel[Greet.FooServiceClient]()
		channel2 := Channel[Greet.FooServiceClient]()

		assert.NotNil(t, channel1, "专用 Channel 应当被创建。")
		assert.Equal(t, channel1, channel2, "专用 Channel 应当是单例的。")
	})
}
