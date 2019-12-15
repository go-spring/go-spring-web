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

//
// 一个 WebServer 包含多个 WebContainer
//
type WebServer struct {
	Containers []WebContainer
}

func NewWebServer() *WebServer {
	return &WebServer{
		Containers: make([]WebContainer, 0),
	}
}

//
// 启动 Web 容器，非阻塞
//
func (s *WebServer) Start() {
	for _, c := range s.Containers {
		c.Start()
	}
}

//
// 停止 Web 容器
//
func (s *WebServer) Stop(ctx context.Context) {
	for _, c := range s.Containers {
		c.Stop(ctx)
	}
}

//
// 添加 WebContainer 实例
//
func (s *WebServer) AddWebContainer(container WebContainer) {
	s.Containers = append(s.Containers, container)
}
