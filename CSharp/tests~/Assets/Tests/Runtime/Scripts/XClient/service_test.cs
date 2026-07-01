// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Net.Http;
using System.Threading;
using System.Threading.Tasks;
using FrameworkEase.Universal.X.gRPC;
using FrameworkEase.Universal.X.gRPC.Tests.Proto.Greet;
using Grpc.Core;
using Grpc.Core.Interceptors;
using NUnit.Framework;

public partial class TestXClient
{
    internal static readonly List<Interceptor> InvokedInterceptors = new();

    internal class GreetFooInterceptor : Interceptor
    {
        public override TResponse BlockingUnaryCall<TRequest, TResponse>(TRequest request, ClientInterceptorContext<TRequest, TResponse> context, BlockingUnaryCallContinuation<TRequest, TResponse> continuation)
        {
            InvokedInterceptors.Add(this);
            return base.BlockingUnaryCall(request, context, continuation);
        }

        public override AsyncUnaryCall<TResponse> AsyncUnaryCall<TRequest, TResponse>(TRequest request, ClientInterceptorContext<TRequest, TResponse> context, AsyncUnaryCallContinuation<TRequest, TResponse> continuation)
        {
            InvokedInterceptors.Add(this);
            return base.AsyncUnaryCall(request, context, continuation);
        }
    }

    internal class GreetBarInterceptor : Interceptor
    {
        public override TResponse BlockingUnaryCall<TRequest, TResponse>(TRequest request, ClientInterceptorContext<TRequest, TResponse> context, BlockingUnaryCallContinuation<TRequest, TResponse> continuation)
        {
            InvokedInterceptors.Add(this);
            return base.BlockingUnaryCall(request, context, continuation);
        }

        public override AsyncUnaryCall<TResponse> AsyncUnaryCall<TRequest, TResponse>(TRequest request, ClientInterceptorContext<TRequest, TResponse> context, AsyncUnaryCallContinuation<TRequest, TResponse> continuation)
        {
            InvokedInterceptors.Add(this);
            return base.AsyncUnaryCall(request, context, continuation);
        }
    }

    [Test]
    public async Task Service()
    {
        #region 启动服务
#if UNITY_5_3_OR_NEWER
        var dotnetVersion = $"net8.0";
        var workingDir = System.IO.Path.GetFullPath("./Assets/Tests/Runtime/Scripts/");
        Process.Start(new ProcessStartInfo { FileName = "dotnet", Arguments = "build", WorkingDirectory = System.IO.Path.GetFullPath("./Assets/Tests/Runtime/Scripts/") }).WaitForExit();
#endif
        var server = Process.Start(new ProcessStartInfo
        {
            FileName = "dotnet",
            Arguments = "FrameworkEase.Universal.X.gRPC.Tests.dll",
#if UNITY_5_3_OR_NEWER
            WorkingDirectory = System.IO.Path.GetFullPath("./Assets/Tests/Runtime/Scripts/bin~/Debug/net8.0")
#endif
        });
        Assert.That(server, Is.Not.Null, "服务器进程应当被启动。");
        XClient.Configure("http://127.0.0.1:50051");
        #endregion

        try
        {
            #region 等待服务
            {
                var probeUrl = "http://127.0.0.1:50052/health";
                using var httpClient = new HttpClient { Timeout = TimeSpan.FromSeconds(2) };
                var cts = new CancellationTokenSource(TimeSpan.FromSeconds(20));
                while (!cts.IsCancellationRequested)
                {
                    try
                    {
                        var response = await httpClient.GetAsync(probeUrl, cts.Token);
                        if (response.IsSuccessStatusCode) break;
                    }
                    catch { }
                    await Task.Delay(100, cts.Token).ConfigureAwait(false);
                }
            }
            #endregion

            #region 验证请求
            {
                XClient<FooService.FooServiceClient>.Intercept(typeof(GreetFooInterceptor));
                XClient<FooService.FooServiceClient>.Intercept(typeof(GreetBarInterceptor));

                var service1 = XClient<FooService.FooServiceClient>.Service;
                var service2 = XClient<FooService.FooServiceClient>.Service;
                Assert.That(service1, Is.EqualTo(service2), "Foo 服务应当是同一个实例。");

                {
                    InvokedInterceptors.Clear();
                    var response = service1.SayHi(new RequestData { Name = "Alice" });
                    Assert.That(response.Message, Is.EqualTo("Hello Alice!"), "Foo 服务应当返回正确的消息。");
                    Assert.That(InvokedInterceptors[0], Is.InstanceOf<GreetFooInterceptor>(), "Foo 拦截器应当被调用。");
                    Assert.That(InvokedInterceptors[1], Is.InstanceOf<GreetBarInterceptor>(), "Bar 拦截器应当被调用。");
                }

                {
                    InvokedInterceptors.Clear();
                    var response = await service1.SayHiAsync(new RequestData { Name = "Alice" });
                    Assert.That(response.Message, Is.EqualTo("Hello Alice!"), "Foo 服务应当返回正确的消息。");
                    Assert.That(InvokedInterceptors[0], Is.InstanceOf<GreetFooInterceptor>(), "Foo 拦截器应当被调用。");
                    Assert.That(InvokedInterceptors[1], Is.InstanceOf<GreetBarInterceptor>(), "Bar 拦截器应当被调用。");
                }
            }

            {
                XClient<BarService.BarServiceClient>.Intercept(typeof(GreetFooInterceptor));
                XClient<BarService.BarServiceClient>.Intercept(typeof(GreetBarInterceptor));

                var service1 = XClient<BarService.BarServiceClient>.Service;
                var service2 = XClient<BarService.BarServiceClient>.Service;
                Assert.That(service1, Is.EqualTo(service2), "Bar 服务应当是同一个实例。");

                {
                    InvokedInterceptors.Clear();
                    var response = service1.SayHello(new RequestData { Name = "Bob" });
                    Assert.That(response.Message, Is.EqualTo("Hi Bob!"), "Bar 服务应当返回正确的消息。");
                    Assert.That(InvokedInterceptors[0], Is.InstanceOf<GreetFooInterceptor>(), "Foo 拦截器应当被调用。");
                    Assert.That(InvokedInterceptors[1], Is.InstanceOf<GreetBarInterceptor>(), "Bar 拦截器应当被调用。");
                }

                {
                    InvokedInterceptors.Clear();
                    var response = await service1.SayHelloAsync(new RequestData { Name = "Bob" });
                    Assert.That(response.Message, Is.EqualTo("Hi Bob!"), "Bar 服务应当返回正确的消息。");
                    Assert.That(InvokedInterceptors[0], Is.InstanceOf<GreetFooInterceptor>(), "Foo 拦截器应当被调用。");
                    Assert.That(InvokedInterceptors[1], Is.InstanceOf<GreetBarInterceptor>(), "Bar 拦截器应当被调用。");
                }
            }
            #endregion
        }
        catch { throw; }
        finally
        {
            #region 停止服务
            if (server != null && !server.HasExited)
            {
                try { server.Kill(); server.WaitForExit(); }
                finally { server.Dispose(); }
            }
            #endregion
        }
    }
}
