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
	ginEngine  *gin.Engine
}

// NewContainer Container 的构造函数
func NewContainer() *Container {
	c := &Container{
		BaseWebContainer: SpringWeb.NewBaseWebContainer(),
	}
	return c
}

// SetGinEngine 设置自定义 gin 引擎
func (c *Container) SetGinEngine(e *gin.Engine) {
	c.ginEngine = e
}

// Start 启动 Web 容器，非阻塞
func (c *Container) Start() {
	address := fmt.Sprintf("%s:%d", c.GetIP(), c.GetPort())

	c.PreStart()

	// 使用默认的 gin 引擎
	if c.ginEngine == nil {
		c.ginEngine = gin.New()
	}

	for _, mapper := range c.Mappers() {
		filters := append(c.GetFilters(), mapper.Filters()...)
		handler := HandlerWrapper(mapper.Path(), mapper.Handler(), filters)
		for _, method := range SpringWeb.GetMethod(mapper.Method()) {
			c.ginEngine.Handle(method, mapper.Path(), handler)
		}
	}

	c.httpServer = &http.Server{
		Addr:    address,
		Handler: c.ginEngine,
	}

	go func() {
		SpringLogger.Info("⇨ http server started on", address)
		var err error
		if c.EnableSSL() {
			err = c.httpServer.ListenAndServeTLS(c.GetCertFile(), c.GetKeyFile())
		} else {
			err = c.httpServer.ListenAndServe()
		}
		SpringLogger.Infof("exit http server on %s return %v", address, err)
	}()
}

// Stop 停止 Web 容器，阻塞
func (c *Container) Stop(ctx context.Context) {
	err := c.httpServer.Shutdown(ctx)
	address := fmt.Sprintf("%s:%d", c.GetIP(), c.GetPort())
	SpringLogger.Infof("shutdown http server on %s return %v", address, err)
}

// HandlerWrapper Web 处理函数包装器
func HandlerWrapper(path string, fn SpringWeb.Handler, filters []SpringWeb.Filter) func(*gin.Context) {
	return func(ginCtx *gin.Context) {

		ctx := ginCtx.Request.Context()
		logCtx := SpringLogger.NewDefaultLoggerContext(ctx)

		webCtx := &Context{
			LoggerContext: logCtx,
			ginContext:    ginCtx,
			handlerPath:   path,
			handlerFunc:   fn,
		}

		SpringWeb.InvokeHandler(webCtx, fn, filters)
	}
}

// Gin Web Gin 适配函数
func Gin(fn gin.HandlerFunc) SpringWeb.Handler {
	return func(webCtx SpringWeb.WebContext) {
		fn(webCtx.NativeContext().(*gin.Context))
	}
}
