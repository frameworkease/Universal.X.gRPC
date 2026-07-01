// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XClient

import (
	"reflect"
	"sort"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
)

type interceptorEntry struct {
	Interceptor grpc.UnaryClientInterceptor
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
	interceptors         = make(map[reflect.Type]interceptorList)
	interceptorsMu       sync.RWMutex
	interceptorIncrement int64
)

func Intercept[T any](interceptor grpc.UnaryClientInterceptor, priority ...int) bool {
	if interceptor == nil {
		panic("XClient.Intercept: interceptor cannot be nil.")
	}
	interceptorPtr := reflect.ValueOf(interceptor).Pointer()
	typ := reflect.TypeFor[T]()
	anyType := reflect.TypeFor[any]()

	interceptorsMu.Lock()
	defer interceptorsMu.Unlock()

	if typ != anyType {
		anyInterceptors := interceptors[anyType]
		for _, entry := range anyInterceptors {
			if reflect.ValueOf(entry.Interceptor).Pointer() == interceptorPtr {
				return false
			}
		}
	}

	typedInterceptors := interceptors[typ]
	for _, entry := range typedInterceptors {
		if reflect.ValueOf(entry.Interceptor).Pointer() == interceptorPtr {
			return false
		}
	}

	increment := atomic.AddInt64(&interceptorIncrement, 1)
	priorityValue := 0
	if len(priority) > 0 {
		priorityValue = priority[0]
	}
	entry := &interceptorEntry{
		Interceptor: interceptor,
		Priority:    priorityValue,
		Increment:   increment,
	}
	typedInterceptors = append(typedInterceptors, entry)
	sort.Sort(typedInterceptors)
	interceptors[typ] = typedInterceptors
	return true
}
