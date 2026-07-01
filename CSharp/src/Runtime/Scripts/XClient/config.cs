// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

using Grpc.Net.Client;

namespace FrameworkEase.Universal.X.gRPC
{
    public partial class XClient
    {
        internal class Configuration
        {
            internal static string Gateway { get; set; }
            internal static GrpcChannelOptions Options { get; set; }

            public static void Configure(string gateway = "", GrpcChannelOptions options = null) { Gateway = gateway; Options = options; }
        }

        public static void Configure(string gateway = "", GrpcChannelOptions options = null) { Configuration.Gateway = gateway; Configuration.Options = options; }
    }

    public partial class XClient<TService>
    {
        internal class Configuration
        {
            internal static string Gateway { get; set; }
            internal static GrpcChannelOptions Options { get; set; }
        }

        public static void Configure(string gateway = "", GrpcChannelOptions options = null) { Configuration.Gateway = gateway; Configuration.Options = options; }
    }
}
