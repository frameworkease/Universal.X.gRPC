// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#if !UNITY_5_3_OR_NEWER
using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using FrameworkEase.Universal.X.gRPC;
using FrameworkEase.Universal.X.gRPC.Tests.Proto.Greet;
using Grpc.Core;
using Grpc.Core.Interceptors;
using Grpc.Net.Client;
using NUnit.Framework;

public partial class TestXServer
{
    internal class GreetServer : XServer.Base
    {
        internal bool awaked;
        internal bool started;
        internal bool stopped;

        public override bool Awake()
        {
            XServer.Serve(typeof(GreetFooService));
            XServer.Serve(typeof(GreetBarService));
            XServer.Intercept(typeof(GreetFooInterceptor));
            XServer.Intercept(typeof(GreetBarInterceptor));
            awaked = true;
            return true;
        }

        public override void Start()
        {
            base.Start();
            started = true;
        }

        public override void Stop(CountdownEvent counter)
        {
            base.Stop(counter);
            stopped = true;
        }
    }

    internal class GreetFooService : FooService.FooServiceBase
    {
        public override Task<ResponseData> SayHi(RequestData request, ServerCallContext context)
        {
            return Task.FromResult(new ResponseData { Message = $"Hello {request.Name}!" });
        }
    }

    internal class GreetBarService : BarService.BarServiceBase
    {
        public override Task<ResponseData> SayHello(RequestData request, ServerCallContext context)
        {
            return Task.FromResult(new ResponseData { Message = $"Hi {request.Name}!" });
        }
    }

    internal static readonly List<Interceptor> InvokedInterceptors = new();

    internal class GreetFooInterceptor : Interceptor
    {
        public override Task<TResponse> UnaryServerHandler<TRequest, TResponse>(TRequest request, ServerCallContext context, UnaryServerMethod<TRequest, TResponse> continuation)
        {
            InvokedInterceptors.Add(this);
            return base.UnaryServerHandler(request, context, continuation);
        }
    }

    internal class GreetBarInterceptor : Interceptor
    {
        public override Task<TResponse> UnaryServerHandler<TRequest, TResponse>(TRequest request, ServerCallContext context, UnaryServerMethod<TRequest, TResponse> continuation)
        {
            InvokedInterceptors.Add(this);
            return base.UnaryServerHandler(request, context, continuation);
        }
    }

    [Test]
    public async Task Lifecycle()
    {
        #region 准备数据
        XServer.Configure(51051, 51052);
        var serveUrl = $"http://127.0.0.1:{XServer.Configuration.ServePort}";
        var probeUrl = $"http://127.0.0.1:{XServer.Configuration.ProbePort}";
        var server = new GreetServer();
        #endregion

        #region 唤醒服务
        {
            server.Awake();
            Assert.That(server.awaked, Is.True, "Greet 服务应当被唤醒。");
        }
        #endregion

        #region 启动服务
        {
            server.Start();
            Assert.That(server.started, Is.True, "Greet 服务应当被启动。");

            #region 等待服务
            {
                using var httpClient = new HttpClient { Timeout = TimeSpan.FromSeconds(2) };
                var cts = new CancellationTokenSource(TimeSpan.FromSeconds(20));
                while (!cts.IsCancellationRequested)
                {
                    try
                    {
                        var response = await httpClient.GetAsync($"{probeUrl}/health", cts.Token);
                        if (response.IsSuccessStatusCode) break;
                    }
                    catch { }
                    await Task.Delay(100, cts.Token).ConfigureAwait(false);
                }
            }
            #endregion

            #region 验证 healthz 端点
            {
                using var client = new HttpClient();
                var healthzResponse = await client.GetAsync($"{probeUrl}/healthz");
                Assert.That(healthzResponse.IsSuccessStatusCode, Is.True, "healthz 端点应当返回成功状态码。");
                Assert.That(healthzResponse.StatusCode, Is.EqualTo(System.Net.HttpStatusCode.OK), "healthz 端点应当返回 200 OK。");
            }
            #endregion

            #region 验证 metrics 端点
            {
                using var client = new HttpClient();
                var metricsResponse = await client.GetAsync($"{probeUrl}/metrics");
                Assert.That(metricsResponse.IsSuccessStatusCode, Is.True, "metrics 端点应当返回成功状态码。");
                var metricsContent = await metricsResponse.Content.ReadAsStringAsync();
                Assert.That(metricsContent, Does.Contain("dotnet_collection_count_total"), "metrics 端点内容应当包含 'dotnet_collection_count_total'。");
            }
            #endregion

            #region 验证 grpc 服务
            {
                using var channel = GrpcChannel.ForAddress(serveUrl);

                {
                    InvokedInterceptors.Clear();
                    var client = new FooService.FooServiceClient(channel);
                    var response = client.SayHi(new RequestData { Name = "Alice" });
                    Assert.That(response.Message, Is.EqualTo("Hello Alice!"), "Foo 服务应当返回正确的消息。");
                    Assert.That(InvokedInterceptors[0], Is.InstanceOf<GreetFooInterceptor>(), "Foo 拦截器应当被调用。");
                    Assert.That(InvokedInterceptors[1], Is.InstanceOf<GreetBarInterceptor>(), "Bar 拦截器应当被调用。");
                }

                {
                    InvokedInterceptors.Clear();
                    var client = new BarService.BarServiceClient(channel);
                    var response = client.SayHello(new RequestData { Name = "Bob" });
                    Assert.That(response.Message, Is.EqualTo("Hi Bob!"), "Bar 服务应当返回正确的消息。");
                    Assert.That(InvokedInterceptors[0], Is.InstanceOf<GreetFooInterceptor>(), "Foo 拦截器应当被调用。");
                    Assert.That(InvokedInterceptors[1], Is.InstanceOf<GreetBarInterceptor>(), "Bar 拦截器应当被调用。");
                }
            }
            #endregion
            #endregion

            #region 停止服务
            {
                var counter = new CountdownEvent(1);
                server.Stop(counter);
                counter.Signal();
                counter.Wait(TimeSpan.FromSeconds(5));
                Assert.That(server.stopped, Is.True, "Greet 服务应当被停止。");
            }
            #endregion
        }
    }
}
#endif
