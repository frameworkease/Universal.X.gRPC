// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#if !UNITY_5_3_OR_NEWER
using FrameworkEase.Universal.X.gRPC;
using FrameworkEase.Universal.X.gRPC.Tests.Proto.Greet;
using NUnit.Framework;

public partial class TestXServer
{
    [Test]
    public void Serve()
    {
        try
        {
            Assert.That(XServer.Serve(typeof(FooService)), Is.True, "Foo 服务应当注册成功。");
            Assert.That(XServer.Serve(typeof(BarService)), Is.True, "Bar 服务应当注册成功。");
            Assert.That(XServer.Serve(typeof(FooService)), Is.False, "Foo 服务只能注册一次。");
            Assert.That(XServer.Serve(typeof(BarService)), Is.False, "Bar 服务只能注册一次。");
            Assert.That(XServer.Services, Does.Contain(typeof(FooService)), "服务列表应当包含 Foo 服务。");
            Assert.That(XServer.Services, Does.Contain(typeof(BarService)), "服务列表应当包含 Bar 服务。");
        }
        finally { XServer.Services.RemoveAll(ele => ele == typeof(FooService) || ele == typeof(BarService)); }
    }
}
#endif
