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
	"time"

	"github.com/go-spring/go-spring-parent/spring-logger"
	"github.com/go-spring/go-spring-parent/spring-utils"
	"github.com/swaggo/http-swagger"
)

// Handler Web 处理接口
type Handler interface {
	Invoke(WebContext)
	FileLine() (file string, line int, fnName string)
}

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

	// setFilters 设置过滤器列表
	setFilters(filters []Filter)

	// AddFilter 添加过滤器
	AddFilter(filter ...Filter)

	// GetLoggerFilter 获取 Logger Filter
	GetLoggerFilter() Filter

	// SetLoggerFilter 设置 Logger Filter
	SetLoggerFilter(filter Filter)

	// GetRecoveryFilter 获取 Recovery Filter
	GetRecoveryFilter() Filter

	// SetRecoveryFilter 设置 Recovery Filter
	SetRecoveryFilter(filter Filter)

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

	loggerFilter   Filter // 日志过滤器
	recoveryFilter Filter // 恢复过滤器
}

// NewBaseWebContainer BaseWebContainer 的构造函数
func NewBaseWebContainer() *BaseWebContainer {
	return &BaseWebContainer{
		WebMapping:     NewDefaultWebMapping(),
		enableSwg:      true,
		loggerFilter:   defaultLoggerFilter,
		recoveryFilter: defaultRecoveryFilter,
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

// setFilters 设置过滤器列表
func (c *BaseWebContainer) setFilters(filters []Filter) {
	c.filters = filters
}

// AddFilter 添加过滤器
func (c *BaseWebContainer) AddFilter(filter ...Filter) {
	c.filters = append(c.filters, filter...)
}

// GetLoggerFilter 获取 Logger Filter
func (c *BaseWebContainer) GetLoggerFilter() Filter {
	return c.loggerFilter
}

// SetLoggerFilter 设置 Logger Filter
func (c *BaseWebContainer) SetLoggerFilter(filter Filter) {
	c.loggerFilter = filter
}

// GetRecoveryFilter 获取 Recovery Filter
func (c *BaseWebContainer) GetRecoveryFilter() Filter {
	return c.recoveryFilter
}

// 设置 Recovery Filter
func (c *BaseWebContainer) SetRecoveryFilter(filter Filter) {
	c.recoveryFilter = filter
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
	file, line, fnName := m.handler.FileLine()
	SpringLogger.Infof("%v :%d %s -> %s:%d %s", GetMethod(m.method), c.port, m.path, file, line, fnName)
}

/////////////////// Invoke Handler //////////////////////

// InvokeHandler 执行 Web 处理函数
func InvokeHandler(ctx WebContext, fn Handler, filters []Filter) {
	if len(filters) > 0 {
		filters = append(filters, HandlerFilter(fn))
		chain := NewDefaultFilterChain(filters)
		chain.Next(ctx)
	} else {
		fn.Invoke(ctx)
	}
}

/////////////////// Web Handlers //////////////////////

// fnHandler 封装 Web 处理函数
type fnHandler func(WebContext)

func (f fnHandler) Invoke(ctx WebContext) {
	if f != nil {
		f(ctx)
	}
}

func (f fnHandler) FileLine() (file string, line int, fnName string) {
	return SpringUtils.FileLine(f)
}

// FUNC 标准 Web 处理函数的辅助函数
func FUNC(fn func(WebContext)) Handler {
	return fnHandler(fn)
}

// httpHandler 标准 Http 处理函数
type httpHandler http.HandlerFunc

func (h httpHandler) Invoke(ctx WebContext) {
	if h != nil {
		h(ctx.ResponseWriter(), ctx.Request())
	}
}

func (h httpHandler) FileLine() (file string, line int, fnName string) {
	return SpringUtils.FileLine(h)
}

// HTTP 标准 Http 处理函数的辅助函数
func HTTP(fn http.HandlerFunc) Handler {
	return httpHandler(fn)
}

/////////////////// Web Filters //////////////////////

var defaultRecoveryFilter = &recoveryFilter{}

// recoveryFilter 恢复过滤器
type recoveryFilter struct{}

func (f *recoveryFilter) Invoke(ctx WebContext, chain FilterChain) {

	defer func() {
		if err := recover(); err != nil {
			ctx.LogError("[PANIC RECOVER] ", err)
			ctx.Status(http.StatusInternalServerError)
		}
	}()

	chain.Next(ctx)
}

var defaultLoggerFilter = &loggerFilter{}

// loggerFilter 日志过滤器
type loggerFilter struct{}

func (f *loggerFilter) Invoke(ctx WebContext, chain FilterChain) {
	start := time.Now()
	chain.Next(ctx)
	ctx.LogInfo("cost: ", time.Since(start))
}
