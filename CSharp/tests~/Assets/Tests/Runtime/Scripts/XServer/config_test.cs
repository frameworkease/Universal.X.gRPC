// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#if !UNITY_5_3_OR_NEWER
using FrameworkEase.Universal.X.gRPC;
using NUnit.Framework;

public partial class TestXServer
{
    [Test]
    public void Configure()
    {
        XServer.Configure(50051, 50052);
        Assert.That(XServer.Configuration.ServePort, Is.EqualTo(50051), "ServePort 配置值应当为 50051。");
        Assert.That(XServer.Configuration.ProbePort, Is.EqualTo(50052), "ProbePort 配置值应当为 50052。");
    }
}
#endif
