// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#if !UNITY_5_3_OR_NEWER
using System;
using System.Collections.Generic;

namespace FrameworkEase.Universal.X.gRPC
{
    public partial class XServer
    {
        internal static readonly List<Type> Services = new();

        public static bool Serve(Type serviceType)
        {
            ArgumentNullException.ThrowIfNull(serviceType);
            lock (Services)
            {
                if (!Services.Contains(serviceType))
                {
                    Services.Add(serviceType);
                    return true;
                }
                return false;
            }
        }
    }
}
#endif
