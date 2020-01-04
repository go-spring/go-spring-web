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

// SetEchoServer 设置自定义 echo 容器
func (c *Container) SetEchoServer(e *echo.Echo) {
	c.echoServer = e
}

// Start 启动 Web 容器，非阻塞
func (c *Container) Start() {

	// 使用默认的 echo 容器
	if c.echoServer == nil {
		e := echo.New()
		c.echoServer = e
		e.HideBanner = true
	}

	// 映射 Web 处理函数
	for _, mapper := range c.Mappers() {
		h := HandlerWrapper(mapper.Handler(), mapper.Filters())
		c.echoServer.Add(mapper.Method(), mapper.Path(), h)
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
		fmt.Println("exit http server on", address, "return", err)
	}()
}

// Stop 停止 Web 容器，阻塞
func (c *Container) Stop(ctx context.Context) {
	if err := c.echoServer.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}

// HandlerWrapper Web 处理函数包装器
func HandlerWrapper(fn SpringWeb.Handler, filters []SpringWeb.Filter) echo.HandlerFunc {
	return func(echoCtx echo.Context) error {

		ctx := echoCtx.Request().Context()
		logCtx := SpringLogger.NewDefaultLoggerContext(ctx)

		webCtx := &Context{
			LoggerContext: logCtx,
			echoContext:   echoCtx,
			handlerFunc:   fn,
		}

		SpringWeb.InvokeHandler(webCtx, fn, filters)
		return nil
	}
}
