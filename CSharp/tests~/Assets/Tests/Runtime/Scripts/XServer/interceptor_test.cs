// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#if !UNITY_5_3_OR_NEWER
using FrameworkEase.Universal.X.gRPC;
using Grpc.Core.Interceptors;
using NUnit.Framework;

public partial class TestXServer
{
    internal class FooInterceptor : Interceptor { }

    internal class BarInterceptor : Interceptor { }

    [TestCase(true)]
    [TestCase(false)]
    public void Intercept(bool priority = false)
    {
        try
        {
            Assert.That(XServer.Intercept(typeof(FooInterceptor), priority ? 1000 : 0), Is.True, "Foo 拦截器应当注册成功。");
            Assert.That(XServer.Intercept(typeof(BarInterceptor), priority ? 1001 : 0), Is.True, "Bar 拦截器应当注册成功。");
            Assert.That(XServer.Intercept(typeof(FooInterceptor)), Is.False, "Foo 拦截器只能注册一次。");
            Assert.That(XServer.Intercept(typeof(BarInterceptor)), Is.False, "Bar 拦截器只能注册一次。");
            Assert.That(XServer.Interceptors.Exists(ele => ele.Interceptor == typeof(FooInterceptor)), Is.True, "拦截器列表应当包含 Foo 拦截器。");
            Assert.That(XServer.Interceptors.Exists(ele => ele.Interceptor == typeof(BarInterceptor)), Is.True, "拦截器列表应当包含 Bar 拦截器。");
            Assert.That(XServer.Interceptors.FindIndex(ele => ele.Interceptor == typeof(FooInterceptor)), Is.LessThan(XServer.Interceptors.FindIndex(ele => ele.Interceptor == typeof(BarInterceptor))), "Foo 拦截器应当排在 Bar 拦截器之前。");
        }
        finally { XServer.Interceptors.RemoveAll(ele => ele.Interceptor == typeof(FooInterceptor) || ele.Interceptor == typeof(BarInterceptor)); }
    }
}
#endif
