// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

using System;
using FrameworkEase.Universal.X.Utility;
using Grpc.Net.Client;

namespace FrameworkEase.Universal.X.gRPC
{
    public partial class XClient
    {
        public static readonly Lazy<GrpcChannel> Channel = new(() =>
        {
            var gateway = Configuration.Gateway;
            var source = "XClient.Configure";
            if (string.IsNullOrEmpty(gateway)) throw new Exception("Gateway is not configured. Please set it via XClient.Configure.");
            var channel = GrpcChannel.ForAddress(gateway, Configuration.Options ??
#if UNITY_5_3_OR_NEWER
            new() { HttpHandler = new Cysharp.Net.Http.YetAnotherHttpHandler() { Http2Only = true } }
#else
            new() { HttpHandler = new System.Net.Http.SocketsHttpHandler { EnableMultipleHttp2Connections = true } }
#endif
            );
            XLog.Notice($"XClient.Channel: created channel to {gateway} (from {source}).");
            return channel;
        });
    }

    public partial class XClient<TService>
    {
        public static readonly Lazy<GrpcChannel> Channel = new(() =>
        {
            var gateway = string.Empty;
            var source = string.Empty;
            if (!string.IsNullOrEmpty(Configuration.Gateway)) { gateway = Configuration.Gateway; source = "XClient<TService>.Configuration"; }
            else if (!string.IsNullOrEmpty(XClient.Configuration.Gateway)) { gateway = XClient.Configuration.Gateway; source = "XClient.Configure"; }
            if (string.IsNullOrEmpty(gateway)) throw new Exception("Gateway is not configured. Please set it via XClient<TService>.Configure or XClient.Configure.");
            var channel = GrpcChannel.ForAddress(gateway, Configuration.Options ?? XClient.Configuration.Options ??
#if UNITY_5_3_OR_NEWER
            new() { HttpHandler = new Cysharp.Net.Http.YetAnotherHttpHandler() { Http2Only = true } }
#else
            new() { HttpHandler = new System.Net.Http.SocketsHttpHandler { EnableMultipleHttp2Connections = true } }
#endif
            );
            XLog.Notice($"XClient<{typeof(TService).Name}>.Channel: created channel to {gateway} (from {source}).");
            return channel;
        });
    }
}
