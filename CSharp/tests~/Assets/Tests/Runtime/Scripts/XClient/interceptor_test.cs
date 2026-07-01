// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

using FrameworkEase.Universal.X.gRPC;
using FrameworkEase.Universal.X.gRPC.Tests.Proto.Greet;
using Grpc.Core.Interceptors;
using NUnit.Framework;

public partial class TestXClient
{
    internal class FooInterceptor : Interceptor { }

    internal class BarInterceptor : Interceptor { }

    [TestCase(true)]
    [TestCase(false)]
    public void Intercept(bool priority = false)
    {
        try
        {
            Assert.That(XClient.Intercept(typeof(FooInterceptor), priority ? 1000 : 0), Is.True, "Foo 拦截器应当注册至共享拦截器列表。");
            Assert.That(XClient<FooService.FooServiceClient>.Intercept(typeof(FooInterceptor), priority ? 1000 : 0), Is.False, "Foo 拦截器只能注册一次至 FooServiceClient 专用拦截器列表。");
            Assert.That(XClient<FooService.FooServiceClient>.Intercept(typeof(BarInterceptor), priority ? 1001 : 0), Is.True, "Bar 拦截器应当注册至 FooServiceClient 专用拦截器列表。");
            Assert.That(XClient<BarService.BarServiceClient>.Intercept(typeof(BarInterceptor), priority ? 1001 : 0), Is.True, "Bar 拦截器应当注册至 BarServiceClient 专用拦截器列表。");
            Assert.That(XClient.Intercept(typeof(BarInterceptor), priority ? 1001 : 0), Is.True, "Bar 拦截器应当注册至共享拦截器列表。");

            Assert.That(XClient.Interceptors.Exists(ele => ele.Interceptor.GetType() == typeof(FooInterceptor)), Is.True, "共享拦截器列表应当包含 Foo 拦截器。");
            Assert.That(XClient.Interceptors.Exists(ele => ele.Interceptor.GetType() == typeof(BarInterceptor)), Is.True, "共享拦截器列表应当包含 Bar 拦截器。");
            Assert.That(XClient.Interceptors.FindIndex(ele => ele.Interceptor.GetType() == typeof(FooInterceptor)), Is.LessThan(XClient.Interceptors.FindIndex(ele => ele.Interceptor.GetType() == typeof(BarInterceptor))), "Foo 拦截器应当排在 Bar 拦截器之前。");

            Assert.That(XClient<FooService.FooServiceClient>.Interceptors.Exists(ele => ele.Interceptor.GetType() == typeof(BarInterceptor)), Is.True, "FooServiceClient 专用拦截器列表应当包含 Bar 拦截器。");
            Assert.That(XClient<BarService.BarServiceClient>.Interceptors.Exists(ele => ele.Interceptor.GetType() == typeof(BarInterceptor)), Is.True, "BarServiceClient 专用拦截器列表应当包含 Bar 拦截器。");
        }
        finally
        {
            XClient.Interceptors.RemoveAll(ele => ele.Interceptor.GetType() == typeof(FooInterceptor) || ele.Interceptor.GetType() == typeof(BarInterceptor));
            XClient<FooService.FooServiceClient>.Interceptors.RemoveAll(ele => ele.Interceptor.GetType() == typeof(FooInterceptor) || ele.Interceptor.GetType() == typeof(BarInterceptor));
            XClient<BarService.BarServiceClient>.Interceptors.RemoveAll(ele => ele.Interceptor.GetType() == typeof(FooInterceptor) || ele.Interceptor.GetType() == typeof(BarInterceptor));
        }
    }
}
