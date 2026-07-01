// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

using FrameworkEase.Universal.X.gRPC;
using FrameworkEase.Universal.X.gRPC.Tests.Proto.Greet;
using Grpc.Net.Client;
using NUnit.Framework;

public partial class TestXClient
{
    [Test]
    public void Configure()
    {
        var gateway = "http://127.0.0.1:50051";
        var options = new GrpcChannelOptions
        {
#if UNITY_5_3_OR_NEWER
            HttpHandler = new Cysharp.Net.Http.YetAnotherHttpHandler() { Http2Only = true }
#else
            HttpHandler = new System.Net.Http.SocketsHttpHandler { EnableMultipleHttp2Connections = true }
#endif
        };

        #region Shared
        {
            XClient.Configure(gateway, options);

            Assert.That(XClient.Configuration.Gateway, Is.EqualTo(gateway), "共享 Gateway 应当被设置。");
            Assert.That(XClient.Configuration.Options, Is.EqualTo(options), "共享 Options 应当被设置。");

            XClient.Configuration.Gateway = null;
            XClient.Configuration.Options = null;
        }
        #endregion

        #region Generic
        {
            XClient<FooService.FooServiceClient>.Configure(gateway, options);

            Assert.That(XClient<FooService.FooServiceClient>.Configuration.Gateway, Is.EqualTo(gateway), "专用 Gateway 应当被设置。");
            Assert.That(XClient<FooService.FooServiceClient>.Configuration.Options, Is.EqualTo(options), "专用 Options 应当被设置。");

            XClient<FooService.FooServiceClient>.Configuration.Gateway = null;
            XClient<FooService.FooServiceClient>.Configuration.Options = null;
        }
        #endregion
    }
}
