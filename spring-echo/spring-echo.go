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

//
// 适配 echo 的 Web 容器
//
type Container struct {
	SpringWeb.BaseWebContainer
	EchoServers []*echo.Echo
}

//
// 构造函数
//
func NewContainer() *Container {
	c := &Container{
		EchoServers: make([]*echo.Echo, 0),
	}
	c.Init()
	return c
}

//
// 启动 web 容器
//
func (c *Container) Start() {
	for _, port := range c.GetPort() {
		address := fmt.Sprintf("%s:%d", c.GetIP(), port)

		e := echo.New()
		e.Use(middleware.Recover())

		for _, mapper := range c.GetMapper() {
			h := HandlerWrapper(mapper.Handler, mapper.Filters)
			switch mapper.Method {
			case "GET":
				e.GET(mapper.Path, h)
			case "PATCH":
				e.PATCH(mapper.Path, h)
			case "PUT":
				e.PUT(mapper.Path, h)
			case "POST":
				e.POST(mapper.Path, h)
			case "DELETE":
				e.DELETE(mapper.Path, h)
			case "HEAD":
				e.HEAD(mapper.Path, h)
			case "OPTIONS":
				e.OPTIONS(mapper.Path, h)
			}
		}

		c.EchoServers = append(c.EchoServers, e)

		go func() {
			var err error

			if c.EnableSSL() {
				err = e.StartTLS(address, c.GetCertFile(), c.GetKeyFile())
			} else {
				err = e.Start(address)
			}

			if err != nil {
				fmt.Println(err)
			}
		}()
	}
}

//
// 停止 Web 容器
//
func (c *Container) Stop(ctx context.Context) {
	for _, s := range c.EchoServers {
		s.Shutdown(ctx)
	}
}

//
// Web 处理函数包装器
//
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
