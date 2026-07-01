// Copyright (c) 2026 FrameworkEase Technologies. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package XClient

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/frameworkease/Universal.X.Utility/Go/XLog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	channels   = make(map[reflect.Type]*grpc.ClientConn)
	channelsMu sync.RWMutex
)

func Channel[T any]() *grpc.ClientConn {
	// 先尝试读取已存在的连接
	typ := reflect.TypeFor[T]()
	channelsMu.RLock()
	if conn, ok := channels[typ]; ok {
		channelsMu.RUnlock()
		return conn
	}
	channelsMu.RUnlock()

	// 需要创建新连接，加写锁
	channelsMu.Lock()
	defer channelsMu.Unlock()

	// 双重检查，避免重复创建
	if conn, ok := channels[typ]; ok {
		return conn
	}

	anyType := reflect.TypeFor[any]()

	// 获取配置
	configsMu.RLock()
	cfg := configs[typ]
	anyCfg := configs[anyType]
	configsMu.RUnlock()

	// 获取拦截器，检查是否有特定类型的拦截器
	interceptorsMu.RLock()
	anyInterceptors := interceptors[anyType]
	typedInterceptors := interceptors[typ]
	interceptorsMu.RUnlock()

	// 获取网关地址：专用的配置 > 共享的配置 > 首选项配置（本地 > 远端 > 资产）
	var gateway string
	var source string
	configsMu.RLock()
	if cfg != nil && cfg.gateway != "" {
		gateway = cfg.gateway
		source = "XClient[T].Configure"
	} else if anyCfg != nil && anyCfg.gateway != "" {
		gateway = anyCfg.gateway
		source = "XClient[any].Configure"
	}
	configsMu.RUnlock()

	if gateway == "" {
		panic("XClient.Channel: gateway is not configured. Please set it via XClient.Configure[T] or XClient.Configure[any].")
	}

	// 构建通道选项
	var dialOptions []grpc.DialOption
	if cfg != nil && len(cfg.options) > 0 { // 使用专用的通道选项
		dialOptions = make([]grpc.DialOption, len(cfg.options))
		copy(dialOptions, cfg.options)
	} else if anyCfg != nil && len(anyCfg.options) > 0 { // 使用共享的通道选项
		dialOptions = make([]grpc.DialOption, len(anyCfg.options))
		copy(dialOptions, anyCfg.options)
	} else { // 使用默认的通道选项
		dialOptions = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}

	// 合并拦截器并去重
	interceptorMap := make(map[uintptr]*interceptorEntry)
	for _, entry := range anyInterceptors {
		ptr := reflect.ValueOf(entry.Interceptor).Pointer()
		if _, ok := interceptorMap[ptr]; !ok {
			interceptorMap[ptr] = entry
		}
	}
	for _, entry := range typedInterceptors {
		ptr := reflect.ValueOf(entry.Interceptor).Pointer()
		if _, ok := interceptorMap[ptr]; !ok {
			interceptorMap[ptr] = entry
		}
	}

	// 转换为列表并排序
	mergedInterceptors := make(interceptorList, 0, len(interceptorMap))
	for _, entry := range interceptorMap {
		mergedInterceptors = append(mergedInterceptors, entry)
	}
	sort.Sort(mergedInterceptors)

	// 添加拦截器到通道选项
	if len(mergedInterceptors) > 0 {
		tempInterceptors := make([]grpc.UnaryClientInterceptor, len(mergedInterceptors))
		for i, entry := range mergedInterceptors {
			tempInterceptors[i] = entry.Interceptor
		}
		dialOptions = append(dialOptions, grpc.WithChainUnaryInterceptor(tempInterceptors...))
	}

	// 创建连接通道
	// 处理目标地址：去除 http:// 或 https:// 前缀
	target := strings.TrimPrefix(strings.TrimPrefix(gateway, "https://"), "http://")
	conn, err := grpc.NewClient(target, dialOptions...)
	if err != nil {
		panic(fmt.Sprintf("XClient.Channel: failed to dial %s: %v", gateway, err))
	}
	XLog.Notice("XClient.Channel: created channel to %s (from %s).", gateway, source)

	// 缓存连接通道
	channels[typ] = conn
	return conn
}
