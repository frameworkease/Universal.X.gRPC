// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#if !UNITY_5_3_OR_NEWER
using System;
using System.Threading;
using FrameworkEase.Universal.X.Utility;
using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Routing;
using Microsoft.AspNetCore.Server.Kestrel.Core;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Prometheus;

namespace FrameworkEase.Universal.X.gRPC
{
    public partial class XServer
    {
        public abstract class Base : XApp.IBase
        {
            private WebApplicationBuilder builder;
            protected WebApplicationBuilder Builder
            {
                get
                {
                    if (builder == null)
                    {
                        builder = WebApplication.CreateBuilder();
                        builder.Logging.ClearProviders();
                    }
                    return builder;
                }
            }

            private WebApplication application;
            protected WebApplication Application
            {
                get
                {
                    if (application == null)
                    {
                        ArgumentOutOfRangeException.ThrowIfNegativeOrZero(Configuration.ServePort);
                        Builder.WebHost.ConfigureKestrel(options => options.ListenAnyIP(Configuration.ServePort, o => o.Protocols = HttpProtocols.Http2));
                        if (Configuration.ProbePort > 0) Builder.WebHost.ConfigureKestrel(options => options.ListenAnyIP(Configuration.ProbePort, o => o.Protocols = HttpProtocols.Http1));
                        Builder.Services.AddGrpc(options => { foreach (var interceptor in Interceptors) options.Interceptors.Add(interceptor.Interceptor); });
                        application = Builder.Build();
                        foreach (var service in Services)
                        {
                            var method = typeof(GrpcEndpointRouteBuilderExtensions).GetMethod("MapGrpcService", new[] { typeof(IEndpointRouteBuilder) });
                            if (method != null)
                            {
                                var genericMethod = method.MakeGenericMethod(service);
                                _ = genericMethod.Invoke(null, new object[] { Application });
                            }
                        }
                        if (Configuration.ProbePort > 0)
                        {
                            Application.MapGet("/health", (Func<IResult>)(() => Results.Ok()));
                            Application.MapGet("/healthz", (Func<IResult>)(() => Results.Ok()));
                            Application.UseMetricServer();
                        }
                    }
                    return application;
                }
            }

            public virtual bool Awake() { return true; }

            public virtual void Start()
            {
                if (Configuration.ProbePort <= 0) XLog.Notice($"XServer.Start: server has been listening on port {Configuration.ServePort}.");
                else XLog.Notice($"XServer.Start: server has been listening on port {Configuration.ServePort} and probe on port {Configuration.ProbePort}.");
                Application.Start();
            }

            public virtual void Stop(CountdownEvent counter)
            {
                counter.AddCount();
                application?.StopAsync();
                counter.Signal();
            }
        }
    }
}
#endif
