// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#if !UNITY_5_3_OR_NEWER
namespace FrameworkEase.Universal.X.gRPC
{
    public abstract partial class XServer
    {
        internal class Configuration
        {
            internal static int ServePort;
            internal static int ProbePort;
        }

        public static void Configure(int servePort, int probePort) { Configuration.ServePort = servePort; Configuration.ProbePort = probePort; }
    }
}
#endif
