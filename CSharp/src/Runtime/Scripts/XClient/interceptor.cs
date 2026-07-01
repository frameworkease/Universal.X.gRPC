// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading;
using Grpc.Core.Interceptors;

namespace FrameworkEase.Universal.X.gRPC
{
    public partial class XClient
    {
        internal class InterceptorEntry : IComparable<InterceptorEntry>
        {
            internal Interceptor Interceptor { get; }
            internal int Priority { get; }
            internal int Increment { get; }

            internal InterceptorEntry(Interceptor interceptor, int priority = 0, int increment = 0)
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
            if (interceptorType == null) throw new ArgumentNullException(nameof(interceptorType));
            if (!interceptorType.IsSubclassOf(typeof(Interceptor))) throw new Exception($"Interceptor type {interceptorType} is not a subclass of {typeof(Interceptor)}");
            lock (Interceptors)
            {
                if (Interceptors.Any(e => e.Interceptor.GetType() == interceptorType)) return false;
                var increment = Interlocked.Increment(ref InterceptorIncrement);
                var entry = new InterceptorEntry((Interceptor)Activator.CreateInstance(interceptorType)!, priority, increment);
                Interceptors.Add(entry);
                Interceptors.Sort();
                return true;
            }
        }
    }

    public partial class XClient<TService>
    {
        internal static readonly List<XClient.InterceptorEntry> Interceptors = new();

        public static bool Intercept(Type interceptorType, int priority = 0)
        {
            if (interceptorType == null) throw new ArgumentNullException(nameof(interceptorType));
            if (!interceptorType.IsSubclassOf(typeof(Interceptor))) throw new Exception($"Interceptor type {interceptorType} is not a subclass of {typeof(Interceptor)}");
            if (XClient.Interceptors.Any(e => e.Interceptor.GetType() == interceptorType)) return false;
            lock (Interceptors)
            {
                if (Interceptors.Any(e => e.Interceptor.GetType() == interceptorType)) return false;
                var increment = Interlocked.Increment(ref XClient.InterceptorIncrement);
                var entry = new XClient.InterceptorEntry((Interceptor)Activator.CreateInstance(interceptorType)!, priority, increment);
                Interceptors.Add(entry);
                Interceptors.Sort();
                return true;
            }
        }
    }
}
