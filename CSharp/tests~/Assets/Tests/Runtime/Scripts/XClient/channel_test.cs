// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

using FrameworkEase.Universal.X.gRPC;
using FrameworkEase.Universal.X.gRPC.Tests.Proto.Greet;
using NUnit.Framework;

public partial class TestXClient
{
    [TestCase("Shared")]
    [TestCase("Generic")]
    public void Channel(string name)
    {
        switch (name)
        {
            case "Shared":
                {
                    XClient.Configure("http://127.0.0.1:50051");
                    var channel1 = XClient.Channel.Value;
                    var channel2 = XClient.Channel.Value;

                    Assert.That(channel1, Is.Not.Null, "共享 Channel 应当被创建。");
                    Assert.That(channel1, Is.EqualTo(channel2), "共享 Channel 应当是单例的。");
                }
                break;
            case "Generic":
                {
                    XClient<FooService.FooServiceClient>.Configure("http://127.0.0.1:50051");
                    var channel1 = XClient<FooService.FooServiceClient>.Channel.Value;
                    var channel2 = XClient<FooService.FooServiceClient>.Channel.Value;

                    Assert.That(channel1, Is.Not.Null, "专用 Channel 应当被创建。");
                    Assert.That(channel1, Is.EqualTo(channel2), "专用 Channel 应当是单例的。");
                }
                break;
        }
    }
}
