// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

using System;
using System.Collections.Generic;
using System.Linq;
using Grpc.Core.Interceptors;

namespace FrameworkEase.Universal.X.gRPC
{
    public partial class XClient<TService>
    {
        internal static readonly Lazy<TService> Instance = new(() =>
        {
            var channel = (!string.IsNullOrEmpty(Configuration.Gateway) || Configuration.Options != null) ? Channel.Value : XClient.Channel.Value;

            var interceptors = new List<XClient.InterceptorEntry>(XClient.Interceptors);
            interceptors.AddRange(Interceptors);
            interceptors = interceptors.GroupBy(e => e.Interceptor.GetType()).Select(g => g.First()).ToList();
            interceptors.Sort();

            var invoker = channel.Intercept(interceptors.Select(e => e.Interceptor).ToArray());
            return (TService)Activator.CreateInstance(typeof(TService), invoker)!;
        });

        public static TService Service => Instance.Value;
    }
}
