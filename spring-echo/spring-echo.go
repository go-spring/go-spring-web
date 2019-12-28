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
	"github.com/labstack/echo/middleware"
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
	address := fmt.Sprintf("%s:%d", c.GetIP(), c.GetPort())

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Recover())

	for _, mapper := range c.Mappers() {
		h := HandlerWrapper(mapper.Handler(), mapper.Filters())
		e.Add(mapper.Method(), mapper.Path(), h)
	}

	c.echoServer = e

	go func() {
		var err error
		if c.EnableSSL() {
			err = e.StartTLS(address, c.GetCertFile(), c.GetKeyFile())
		} else {
			err = e.Start(address)
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
			DefaultLoggerContext: logCtx,
			EchoContext:          echoCtx,
			HandlerFunc:          fn,
		}

		SpringWeb.InvokeHandler(webCtx, fn, filters)
		return nil
	}
}
