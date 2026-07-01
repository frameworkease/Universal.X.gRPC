// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#if !UNITY_5_3_OR_NEWER
using FrameworkEase.Universal.X.gRPC;
using FrameworkEase.Universal.X.Utility;

XServer.Configure(50051, 50052);
await XApp.Run(new TestXServer.GreetServer());
#endif
