// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XServer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	Configure(50051, 50052)
	assert.Equal(t, 50051, configuredServePort, "ServePort 配置值应当为 50051。")
	assert.Equal(t, 50052, configuredProbePort, "ProbePort 配置值应当为 50052。")
}
