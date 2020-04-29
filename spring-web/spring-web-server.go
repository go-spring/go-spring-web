/*
 * Copyright 2012-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package SpringWeb

import (
	"context"
)

// WebServer 一个 WebServer 包含多个 WebContainer
type WebServer struct {
	filters    []Filter
	Containers []WebContainer

	loggerFilter   Filter // 日志过滤器
	recoveryFilter Filter // 恢复过滤器
}

// NewWebServer WebServer 的构造函数
func NewWebServer() *WebServer {
	return &WebServer{
		Containers: make([]WebContainer, 0),
	}
}

// AddFilter 添加过滤器
func (s *WebServer) AddFilter(filter ...Filter) {
	s.filters = append(s.filters, filter...)
}

// SetLoggerFilter 设置 Logger Filter
func (s *WebServer) SetLoggerFilter(filter Filter) {
	s.loggerFilter = filter
}

// 设置 Recovery Filter
func (s *WebServer) SetRecoveryFilter(filter Filter) {
	s.recoveryFilter = filter
}

// AddWebContainer 添加 WebContainer 实例
func (s *WebServer) AddWebContainer(container WebContainer) {
	s.Containers = append(s.Containers, container)
}

// Start 启动 Web 容器，非阻塞
func (s *WebServer) Start() {
	for _, c := range s.Containers {

		// Container 使用 Server 的日志过滤器
		if s.loggerFilter != nil && c.GetLoggerFilter() == nil {
			c.SetLoggerFilter(s.loggerFilter)
		}

		// Container 使用 Server 的恢复过滤器
		if s.recoveryFilter != nil && c.GetRecoveryFilter() == nil {
			c.SetRecoveryFilter(s.recoveryFilter)
		}

		// 添加 Server 的过滤器给 Container
		filters := append(s.filters, c.GetFilters()...)
		c.setFilters(filters)

		c.Start()
	}
}

// Stop 停止 Web 容器，阻塞
func (s *WebServer) Stop(ctx context.Context) {
	for _, c := range s.Containers {
		c.Stop(ctx)
	}
}
