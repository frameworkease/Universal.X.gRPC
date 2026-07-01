// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

#if !UNITY_5_3_OR_NEWER
using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;

namespace FrameworkEase.Universal.X.gRPC
{
    public partial class XServer
    {
        internal class InterceptorEntry : IComparable<InterceptorEntry>
        {
            internal Type Interceptor { get; }
            internal int Priority { get; }
            internal int Increment { get; }

            internal InterceptorEntry(Type interceptor, int priority = 0, int increment = 0)
            {
                Interceptor = interceptor;
                Priority = priority;
                Increment = increment;
            }

            public int CompareTo(InterceptorEntry other)
            {
                if (other == null) return 1;
                var priorityComparison = Priority.CompareTo(other.Priority);
                if (priorityComparison != 0) return priorityComparison;
                return Increment.CompareTo(other.Increment);
            }
        }

        internal static readonly List<InterceptorEntry> Interceptors = new();
        internal static int InterceptorIncrement = 0;

        public static bool Intercept(Type interceptorType, int priority = 0)
        {
            ArgumentNullException.ThrowIfNull(interceptorType);
            lock (Interceptors)
            {
                if (Interceptors.Any(e => e.Interceptor == interceptorType)) return false;
                var increment = Interlocked.Increment(ref InterceptorIncrement);
                var entry = new InterceptorEntry(interceptorType, priority, increment);
                Interceptors.Add(entry);
                Interceptors.Sort();
                return true;
            }
        }
    }
}
#endif
