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
	"fmt"

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
func NewContainer() *Container {
	c := &Container{
		BaseWebContainer: SpringWeb.NewBaseWebContainer(),
	}
	return c
}

// Start 启动 Web 容器，非阻塞
func (c *Container) Start() {

	c.PreStart()

	// 使用默认的 echo 容器
	if c.echoServer == nil {
		e := echo.New()
		c.echoServer = e
		e.HideBanner = true
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

		path := SpringWeb.PathConvert(mapper.Path())
		filters := append(cFilters, mapper.Filters()...)
		handler := HandlerWrapper(mapper.Handler(), filters)

		for _, method := range SpringWeb.GetMethod(mapper.Method()) {
			c.echoServer.Add(method, path, handler)
		}
	}

	// 启动 echo 容器
	go func() {
		address := fmt.Sprintf("%s:%d", c.GetIP(), c.GetPort())
		var err error
		if c.EnableSSL() {
			err = c.echoServer.StartTLS(address, c.GetCertFile(), c.GetKeyFile())
		} else {
			err = c.echoServer.Start(address)
		}
		SpringLogger.Infof("exit http server on %s return %v", address, err)
	}()
}

// Stop 停止 Web 容器，阻塞
func (c *Container) Stop(ctx context.Context) {
	err := c.echoServer.Shutdown(ctx)
	address := fmt.Sprintf("%s:%d", c.GetIP(), c.GetPort())
	SpringLogger.Infof("shutdown http server on %s return %v", address, err)
}

// HandlerWrapper Web 处理函数包装器
func HandlerWrapper(fn SpringWeb.Handler, filters []SpringWeb.Filter) echo.HandlerFunc {
	return func(echoCtx echo.Context) error {
		webCtx := NewContext(fn, echoCtx)
		SpringWeb.InvokeHandler(webCtx, fn, filters)
		return nil
	}
}

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

// echoFilter 封装 Echo 中间件
type echoFilter struct {
	fn echo.MiddlewareFunc
}

func (f *echoFilter) Invoke(ctx SpringWeb.WebContext, chain *SpringWeb.FilterChain) {

	h := f.fn(func(echoCtx echo.Context) error {
		chain.Next(ctx)
		return nil
	})

	if err := h(EchoContext(ctx)); err != nil {
		panic(err)
	}
}

// Filter Web Echo 中间件适配器
func Filter(fn echo.MiddlewareFunc) SpringWeb.Filter {
	return &echoFilter{fn: fn}
}
