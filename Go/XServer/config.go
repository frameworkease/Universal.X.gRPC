// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XServer

var (
	configuredServePort int
	configuredProbePort int
)

func Configure(servePort, probePort int) {
	configuredServePort = servePort
	configuredProbePort = probePort
}
