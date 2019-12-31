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
	"net/http"
)

// Router 路由分组
type Router struct {
	mapping  WebMapping
	basePath string
	filters  []Filter
}

// NewRouter Router 的构造函数
func NewRouter(mapping WebMapping, path string, filters []Filter) *Router {
	return &Router{
		mapping:  mapping,
		basePath: path,
		filters:  filters,
	}
}

// Request 注册任意 HTTP 方法处理函数
func (g *Router) Request(method string, path string, fn Handler) *Router {
	g.mapping.Request(method, g.basePath+path, fn, g.filters...)
	return g
}

// GET 注册 GET 方法处理函数
func (g *Router) GET(path string, fn Handler) *Router {
	return g.Request(http.MethodGet, path, fn)
}

// POST 注册 POST 方法处理函数
func (g *Router) POST(path string, fn Handler) *Router {
	return g.Request(http.MethodPost, path, fn)
}

// PATCH 注册 PATCH 方法处理函数
func (g *Router) PATCH(path string, fn Handler) *Router {
	return g.Request(http.MethodPatch, path, fn)
}

// PUT 注册 PUT 方法处理函数
func (g *Router) PUT(path string, fn Handler) *Router {
	return g.Request(http.MethodPut, path, fn)
}

// DELETE 注册 DELETE 方法处理函数
func (g *Router) DELETE(path string, fn Handler) *Router {
	return g.Request(http.MethodDelete, path, fn)
}

// HEAD 注册 HEAD 方法处理函数
func (g *Router) HEAD(path string, fn Handler) *Router {
	return g.Request(http.MethodHead, path, fn)
}

// OPTIONS 注册 OPTIONS 方法处理函数
func (g *Router) OPTIONS(path string, fn Handler) *Router {
	return g.Request(http.MethodOptions, path, fn)
}
