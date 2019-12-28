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

package SpringGin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-spring/go-spring-parent/spring-logger"
	"github.com/go-spring/go-spring-web/spring-web"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

// Container 适配 gin 的 Web 容器
type Container struct {
	*SpringWeb.BaseWebContainer
	httpServer *http.Server
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

	ginEngine := gin.New()
	ginEngine.Use(gin.Logger(), gin.Recovery())

	for _, mapper := range c.Mappers() {
		h := HandlerWrapper(mapper.Path(), mapper.Handler(), mapper.Filters())
		ginEngine.Handle(mapper.Method(), mapper.Path(), h)
	}

	c.httpServer = &http.Server{
		Addr:    address,
		Handler: ginEngine,
	}

	go func() {
		fmt.Printf("⇨ http server started on %s\n", address)

		var err error
		if c.EnableSSL() {
			err = c.httpServer.ListenAndServeTLS(c.GetCertFile(), c.GetKeyFile())
		} else {
			err = c.httpServer.ListenAndServe()
		}
		fmt.Println("exit http server on", address, "return", err)
	}()
}

// Stop 停止 Web 容器，阻塞
func (c *Container) Stop(ctx context.Context) {
	if err := c.httpServer.Shutdown(ctx); err != nil {
		fmt.Println(err)
	}
}

// HandlerWrapper Web 处理函数包装器
func HandlerWrapper(path string, fn SpringWeb.Handler, filters []SpringWeb.Filter) func(*gin.Context) {
	return func(ginCtx *gin.Context) {

		ctx := ginCtx.Request.Context()
		logCtx := SpringLogger.NewDefaultLoggerContext(ctx)

		webCtx := &Context{
			DefaultLoggerContext: logCtx,
			GinContext:           ginCtx,
			HandlerPath:          path,
			HandlerFunc:          fn,
		}

		SpringWeb.InvokeHandler(webCtx, fn, filters)
	}
}
