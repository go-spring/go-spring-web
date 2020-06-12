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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-spring/go-spring-parent/spring-logger"
	"github.com/go-spring/go-spring-parent/spring-utils"
	"github.com/go-spring/go-spring-web/spring-web"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	binding.Validator = SpringWeb.NewBuiltInValidator()
}

// Container 适配 gin 的 Web 容器
type Container struct {
	*SpringWeb.BaseWebContainer
	httpServer *http.Server
	ginEngine  *gin.Engine
}

// NewContainer Container 的构造函数
func NewContainer(config SpringWeb.ContainerConfig) *Container {
	c := &Container{
		BaseWebContainer: SpringWeb.NewBaseWebContainer(config),
	}
	return c
}

// Deprecated: Filter 机制可完美替代中间件机制，不再需要定制化
func (c *Container) SetGinEngine(e *gin.Engine) {
	c.ginEngine = e
}

// Start 启动 Web 容器，非阻塞
func (c *Container) Start() {

	c.PreStart()

	// 使用默认的 gin 引擎
	if c.ginEngine == nil {
		c.ginEngine = gin.New()
	}

	var cFilters []SpringWeb.Filter

	if f := c.GetLoggerFilter(); f != nil {
		cFilters = append(cFilters, f)
	}

	if f := c.GetRecoveryFilter(); f != nil {
		cFilters = append(cFilters, f)
	}

	cFilters = append(cFilters, c.GetFilters()...)

	// TODO 添加全局过滤器

	// 映射 Web 处理函数
	for _, mapper := range c.Mappers() {
		c.PrintMapper(mapper)

		path, wildCardName := SpringWeb.ToPathStyle(mapper.Path(), SpringWeb.GinPathStyle)
		filters := append(cFilters, mapper.Filters()...)
		handlers := HandlerWrapper(mapper.Path(), mapper.Handler(), wildCardName, filters)

		for _, method := range SpringWeb.GetMethod(mapper.Method()) {
			c.ginEngine.Handle(method, path, handlers...)
		}
	}

	go func() {
		var err error
		cfg := c.Config()

		c.httpServer = &http.Server{
			Addr:         c.Address(),
			Handler:      c.ginEngine,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		}

		SpringLogger.Info("⇨ http server started on ", c.Address())

		if cfg.EnableSSL {
			err = c.httpServer.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = c.httpServer.ListenAndServe()
		}
		SpringLogger.Infof("exit gin server on %s return %s", c.Address(), SpringUtils.ToString(err))
	}()
}

// Stop 停止 Web 容器，阻塞
func (c *Container) Stop(ctx context.Context) {
	err := c.httpServer.Shutdown(ctx)
	SpringLogger.Infof("shutdown gin server on %s return %s", c.Address(), SpringUtils.ToString(err))
}

// HandlerWrapper Web 处理函数包装器
func HandlerWrapper(path string, fn SpringWeb.Handler, wildCardName string, filters []SpringWeb.Filter) []gin.HandlerFunc {
	var handlers []gin.HandlerFunc

	// 建立 WebContext 和 GinContext 之间的关联
	handlers = append(handlers, func(ginCtx *gin.Context) {
		NewContext(path, fn, wildCardName, ginCtx)
	})

	// 封装过滤器
	for _, filter := range filters {
		f := filter // 避免延迟绑定
		handlers = append(handlers, func(ginCtx *gin.Context) {
			f.Invoke(WebContext(ginCtx), &ginFilterChain{ginCtx})
		})
	}

	// 封装 Web 处理函数
	handlers = append(handlers, func(ginCtx *gin.Context) {
		fn.Invoke(WebContext(ginCtx))
	})

	return handlers
}

/////////////////// handler //////////////////////

// ginHandler 封装 Gin 处理函数
type ginHandler gin.HandlerFunc

func (g ginHandler) Invoke(ctx SpringWeb.WebContext) {
	g(GinContext(ctx))
}

func (g ginHandler) FileLine() (file string, line int, fnName string) {
	return SpringUtils.FileLine(g)
}

// Gin Web Gin 适配函数
func Gin(fn gin.HandlerFunc) SpringWeb.Handler {
	return ginHandler(fn)
}

/////////////////// filter //////////////////////

// ginFilter 封装 Gin 中间件
type ginFilter gin.HandlerFunc

func (filter ginFilter) Invoke(ctx SpringWeb.WebContext, _ SpringWeb.FilterChain) {
	filter(GinContext(ctx))
}

// Filter Web Gin 中间件适配器
func Filter(fn gin.HandlerFunc) SpringWeb.Filter {
	return ginFilter(fn)
}

// ginFilterChain gin 适配的过滤器链条
type ginFilterChain struct {
	ginCtx *gin.Context
}

// Next 内部调用 gin.Context 对象的 Next 函数驱动链条向后执行
func (chain *ginFilterChain) Next(_ SpringWeb.WebContext) {
	chain.ginCtx.Next()
}
