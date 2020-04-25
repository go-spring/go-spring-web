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
	"net/http"
	"reflect"
	"runtime"

	"github.com/go-spring/go-spring-parent/spring-logger"
	"github.com/swaggo/http-swagger"
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

	// GetFilters 返回过滤器列表
	GetFilters() []Filter

	// SetFilters 设置过滤器列表
	SetFilters(filters ...Filter)

	// AddRouter 添加新的路由信息
	AddRouter(router *Router)

	// EnableSwagger 是否启用 Swagger 功能
	EnableSwagger() bool

	// SetEnableSwagger 设置是否启用 Swagger 功能
	SetEnableSwagger(enable bool)

	// Start 启动 Web 容器，非阻塞
	Start()

	// Stop 停止 Web 容器，阻塞
	Stop(ctx context.Context)
}

// BaseWebContainer WebContainer 的通用部分
type BaseWebContainer struct {
	WebMapping

	ip        string // 监听 IP
	port      int    // 监听端口
	enableSSL bool   // 使用 SSL
	keyFile   string
	certFile  string
	filters   []Filter
	enableSwg bool // 是否启用 Swagger 功能
}

// NewBaseWebContainer BaseWebContainer 的构造函数
func NewBaseWebContainer() *BaseWebContainer {
	return &BaseWebContainer{
		WebMapping: NewDefaultWebMapping(),
		enableSwg:  true,
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

// GetFilters 返回过滤器列表
func (c *BaseWebContainer) GetFilters() []Filter {
	return c.filters
}

// SetFilters 设置过滤器列表
func (c *BaseWebContainer) SetFilters(filters ...Filter) {
	c.filters = filters
}

// AddRouter 添加新的路由信息
func (c *BaseWebContainer) AddRouter(router *Router) {
	for _, mapper := range router.mapping.Mappers() {
		c.AddMapper(mapper)
	}
}

// EnableSwagger 是否启用 Swagger 功能
func (c *BaseWebContainer) EnableSwagger() bool {
	return c.enableSwg
}

// SetEnableSwagger 设置是否启用 Swagger 功能
func (c *BaseWebContainer) SetEnableSwagger(enable bool) {
	c.enableSwg = enable
}

// PreStart 执行 Start 之前的准备工作
func (c *BaseWebContainer) PreStart() {

	if c.enableSwg {

		// 注册 path 的 Operation
		for _, mapper := range c.Mappers() {
			if op := mapper.swagger; op != nil {
				if err := op.parseBind(); err != nil {
					panic(err)
				}
				doc.AddPath(mapper.Path(), mapper.Method(), op)
			}
		}

		// 注册 swagger-ui 和 doc.json 接口
		c.GET("/swagger/*", HTTP(httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		)))

		// 注册 redoc 接口
		c.GET("/redoc", ReDoc)
	}
}

// PrintMapper 打印路由注册信息
func (c *BaseWebContainer) PrintMapper(m *Mapper) {
	fnPtr := reflect.ValueOf(m.handler).Pointer()
	fnInfo := runtime.FuncForPC(fnPtr)
	file, line := fnInfo.FileLine(fnPtr)
	SpringLogger.Infof("%v :%d %s -> %s:%d", GetMethod(m.method), c.port, m.path, file, line)
}

/////////////////// Invoke Handler //////////////////////

// InvokeHandler 执行 Web 处理函数
func InvokeHandler(ctx WebContext, fn Handler, filters []Filter) {

	defer func() {
		if err := recover(); err != nil {
			ctx.LogError(err)
			ctx.Status(http.StatusInternalServerError)
		}
	}()

	if len(filters) > 0 {
		filters = append(filters, HandlerFilter(fn))
		chain := NewFilterChain(filters)
		chain.Next(ctx)
	} else {
		fn(ctx)
	}
}

// HTTP Web HTTP 适配函数
func HTTP(fn http.HandlerFunc) Handler {
	return func(webCtx WebContext) {
		fn(webCtx.ResponseWriter(), webCtx.Request())
	}
}
