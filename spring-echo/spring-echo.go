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

package SpringEcho

import (
	"context"

	"github.com/go-spring/go-spring-parent/spring-logger"
	"github.com/go-spring/go-spring-parent/spring-utils"
	"github.com/go-spring/go-spring-web/spring-web"
	"github.com/labstack/echo"
)

// Container 适配 echo 的 Web 容器
type Container struct {
	*SpringWeb.BaseWebContainer
	echoServer *echo.Echo
}

// NewContainer Container 的构造函数
func NewContainer(config SpringWeb.ContainerConfig) *Container {
	c := &Container{
		BaseWebContainer: SpringWeb.NewBaseWebContainer(config),
	}
	return c
}

// Deprecated: Filter 机制可完美替代中间件机制，不再需要定制化
func (c *Container) SetEchoServer(e *echo.Echo) {
	c.echoServer = e
}

// Start 启动 Web 容器，非阻塞
func (c *Container) Start() {

	c.PreStart()

	// 使用默认的 echo 容器
	if c.echoServer == nil {
		e := echo.New()
		e.HideBanner = true
		c.echoServer = e
	}

	var cFilters []SpringWeb.Filter

	if f := c.GetLoggerFilter(); f != nil {
		cFilters = append(cFilters, f)
	}

	if f := c.GetRecoveryFilter(); f != nil {
		cFilters = append(cFilters, f)
	}

	cFilters = append(cFilters, c.GetFilters()...)

	// 映射 Web 处理函数
	for _, mapper := range c.Mappers() {
		c.PrintMapper(mapper)

		path, wildCardName := SpringWeb.ToPathStyle(mapper.Path(), SpringWeb.EchoPathStyle)
		filters := append(cFilters, mapper.Filters()...)
		handler := HandlerWrapper(mapper.Handler(), wildCardName, filters)

		for _, method := range SpringWeb.GetMethod(mapper.Method()) {
			c.echoServer.Add(method, path, handler)
		}
	}

	// 启动 echo 容器
	go func() {
		// TODO 应用 ReadTimeout 和 WriteTimeout。

		var err error
		if cfg := c.Config(); cfg.EnableSSL {
			err = c.echoServer.StartTLS(c.Address(), cfg.CertFile, cfg.KeyFile)
		} else {
			err = c.echoServer.Start(c.Address())
		}
		SpringLogger.Infof("exit echo server on %s return %s", c.Address(), SpringUtils.ToString(err))
	}()
}

// Stop 停止 Web 容器，阻塞
func (c *Container) Stop(ctx context.Context) {
	err := c.echoServer.Shutdown(ctx)
	SpringLogger.Infof("shutdown echo server on %s return %s", c.Address(), SpringUtils.ToString(err))
}

// HandlerWrapper Web 处理函数包装器
func HandlerWrapper(fn SpringWeb.Handler, wildCardName string, filters []SpringWeb.Filter) echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		webCtx := NewContext(fn, wildCardName, echoCtx)
		SpringWeb.InvokeHandler(webCtx, fn, filters)
		return nil
	}
}

/////////////////// handler //////////////////////

// echoHandler 封装 Echo 处理函数
type echoHandler echo.HandlerFunc

func (e echoHandler) Invoke(ctx SpringWeb.WebContext) {
	if err := e(EchoContext(ctx)); err != nil {
		panic(err)
	}
}

func (e echoHandler) FileLine() (file string, line int, fnName string) {
	return SpringUtils.FileLine(e)
}

// Echo Web Echo 适配函数
func Echo(fn echo.HandlerFunc) SpringWeb.Handler {
	return echoHandler(fn)
}

/////////////////// filter //////////////////////

// echoFilter 封装 Echo 中间件
type echoFilter echo.MiddlewareFunc

func (filter echoFilter) Invoke(ctx SpringWeb.WebContext, chain SpringWeb.FilterChain) {

	h := filter(func(echoCtx echo.Context) error {
		chain.Next(ctx)
		return nil
	})

	if err := h(EchoContext(ctx)); err != nil {
		panic(err)
	}
}

// Filter Web Echo 中间件适配器
func Filter(fn echo.MiddlewareFunc) SpringWeb.Filter {
	return echoFilter(fn)
}
