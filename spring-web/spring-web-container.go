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

// Handler Web 处理函数
type Handler func(WebContext)

// WebContainer Web 容器
type WebContainer interface {
	// WebMapping 路由表
	WebMapping

	// GetIP 返回监听的 IP
	GetIP() string

	// SetIP 设置监听的 IP
	SetIP(ip string)

	// GetPort 返回监听的 Port
	GetPort() int

	// SetPort 设置监听的 Port
	SetPort(port int)

	// EnableSSL 返回是否启用 SSL
	EnableSSL() bool

	// SetEnableSSL 设置是否启用 SSL
	SetEnableSSL(enable bool)

	// GetKeyFile 返回 KeyFile 的路径
	GetKeyFile() string

	// SetKeyFile 设置 KeyFile 的路径
	SetKeyFile(keyFile string)

	// GetCertFile 返回 CertFile 的路径
	GetCertFile() string

	// SetCertFile 设置 CertFile 的路径
	SetCertFile(certFile string)

	// Start 启动 Web 容器，非阻塞
	Start()

	// Stop 停止 Web 容器，阻塞
	Stop(ctx context.Context)
}

// BaseWebContainer WebContainer 的通用部分
type BaseWebContainer struct {
	WebMapping

	ip   string // 监听的 IP
	port int    // 监听的端口

	enableSSL bool
	keyFile   string
	certFile  string
}

// NewBaseWebContainer BaseWebContainer 的构造函数
func NewBaseWebContainer() *BaseWebContainer {
	return &BaseWebContainer{
		WebMapping: NewDefaultWebMapping(),
	}
}

// GetIP 返回监听的 IP
func (c *BaseWebContainer) GetIP() string {
	return c.ip
}

// SetIP 设置监听的 IP
func (c *BaseWebContainer) SetIP(ip string) {
	c.ip = ip
}

// GetPort 返回监听的 Port
func (c *BaseWebContainer) GetPort() int {
	return c.port
}

// SetPort 设置监听的 Port
func (c *BaseWebContainer) SetPort(port int) {
	c.port = port
}

// EnableSSL 返回是否启用 SSL
func (c *BaseWebContainer) EnableSSL() bool {
	return c.enableSSL
}

// SetEnableSSL 设置是否启用 SSL
func (c *BaseWebContainer) SetEnableSSL(enable bool) {
	c.enableSSL = enable
}

// GetKeyFile 返回 KeyFile 的路径
func (c *BaseWebContainer) GetKeyFile() string {
	return c.keyFile
}

// SetKeyFile 设置 KeyFile 的路径
func (c *BaseWebContainer) SetKeyFile(keyFile string) {
	c.keyFile = keyFile
}

// GetCertFile 返回 CertFile 的路径
func (c *BaseWebContainer) GetCertFile() string {
	return c.certFile
}

// SetCertFile 设置 CertFile 的路径
func (c *BaseWebContainer) SetCertFile(certFile string) {
	c.certFile = certFile
}

/////////////////////////////////////////////////////////

// InvokeHandler 执行 Web 处理函数
func InvokeHandler(ctx WebContext, fn Handler, filters []Filter) {
	if len(filters) > 0 {
		filters = append(filters, HandlerFilter(fn))
		chain := NewFilterChain(filters)
		chain.Next(ctx)
	} else {
		fn(ctx)
	}
}

/////////////////////////////////////////////////////////

// 定义 WebContainer 的工厂函数
type Factory func() WebContainer

// 保存 WebContainer 的工厂函数
var WebContainerFactory Factory

// 注册 WebContainer 的工厂函数
func RegisterWebContainerFactory(fn Factory) {
	WebContainerFactory = fn
}
