// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XServer

import (
	"reflect"
	"sort"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
)

type interceptorEntry struct {
	Interceptor grpc.UnaryServerInterceptor
	Priority    int
	Increment   int64
}

type interceptorList []*interceptorEntry

func (l interceptorList) Len() int { return len(l) }

func (l interceptorList) Less(i, j int) bool {
	if l[i].Priority != l[j].Priority {
		return l[i].Priority < l[j].Priority
	}
	return l[i].Increment < l[j].Increment
}

func (l interceptorList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

var (
	interceptors         = interceptorList{}
	interceptorsMu       sync.RWMutex
	interceptorIncrement atomic.Int64
)

func Intercept(interceptor grpc.UnaryServerInterceptor, priority ...int) bool {
	if interceptor == nil {
		panic("XServer.Intercept: interceptor cannot be nil.")
	}
	interceptorsMu.Lock()
	defer interceptorsMu.Unlock()

	interceptorPtr := reflect.ValueOf(interceptor).Pointer()
	for _, entry := range interceptors {
		if reflect.ValueOf(entry.Interceptor).Pointer() == interceptorPtr {
			return false
		}
	}

	increment := interceptorIncrement.Add(1)
	priorityValue := 0
	if len(priority) > 0 {
		priorityValue = priority[0]
	}
	entry := &interceptorEntry{
		Interceptor: interceptor,
		Priority:    priorityValue,
		Increment:   increment,
	}
	interceptors = append(interceptors, entry)
	sort.Sort(interceptors)
	return true
}
