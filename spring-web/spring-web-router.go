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

// Router 路由分组
type Router struct {
	mapping  WebMapping
	basePath string
	filters  []Filter
}

// NewRouter Router 的构造函数，不依赖具体的 WebMapping 对象
func NewRouter(basePath string, filters ...Filter) *Router {
	return &Router{
		mapping:  NewDefaultWebMapping(),
		basePath: basePath,
		filters:  filters,
	}
}

// Request 注册任意 HTTP 方法处理函数
func (r *Router) Request(method uint32, path string, fn interface{}, filters ...Filter) *Mapper {
	filters = append(r.filters, filters...)
	return r.mapping.Request(method, r.basePath+path, fn, filters...)
}

// Deprecated: 推荐使用 Get* 系列函数进行编译检查
func (r *Router) GET(path string, fn interface{}, filters ...Filter) *Mapper {
	return r.Request(MethodGet, path, fn, filters...)
}

// HandleGet 注册 GET 方法处理函数
func (r *Router) HandleGet(path string, fn Handler, filters ...Filter) *Mapper {
	return r.Request(MethodGet, path, fn, filters...)
}

// GetMapping 注册 GET 方法处理函数
func (r *Router) GetMapping(path string, fn HandlerFunc, filters ...Filter) *Mapper {
	return r.Request(MethodGet, path, FUNC(fn), filters...)
}

// GetBinding 注册 GET 方法处理函数
func (r *Router) GetBinding(path string, fn interface{}, filters ...Filter) *Mapper {
	return r.Request(MethodGet, path, BIND(fn), filters...)
}

// Deprecated: 推荐使用 Post* 系列函数进行编译检查
func (r *Router) POST(path string, fn interface{}, filters ...Filter) *Mapper {
	return r.Request(MethodPost, path, fn, filters...)
}

// HandlePost 注册 POST 方法处理函数
func (r *Router) HandlePost(path string, fn Handler, filters ...Filter) *Mapper {
	return r.Request(MethodPost, path, fn, filters...)
}

// PostMapping 注册 POST 方法处理函数
func (r *Router) PostMapping(path string, fn HandlerFunc, filters ...Filter) *Mapper {
	return r.Request(MethodPost, path, FUNC(fn), filters...)
}

// PostBinding 注册 POST 方法处理函数
func (r *Router) PostBinding(path string, fn interface{}, filters ...Filter) *Mapper {
	return r.Request(MethodPost, path, BIND(fn), filters...)
}

// PATCH 注册 PATCH 方法处理函数
func (r *Router) PATCH(path string, fn interface{}, filters ...Filter) *Mapper {
	return r.Request(MethodPatch, path, fn, filters...)
}

// PUT 注册 PUT 方法处理函数
func (r *Router) PUT(path string, fn interface{}, filters ...Filter) *Mapper {
	return r.Request(MethodPut, path, fn, filters...)
}

// DELETE 注册 DELETE 方法处理函数
func (r *Router) DELETE(path string, fn interface{}, filters ...Filter) *Mapper {
	return r.Request(MethodDelete, path, fn, filters...)
}

// HEAD 注册 HEAD 方法处理函数
func (r *Router) HEAD(path string, fn interface{}, filters ...Filter) *Mapper {
	return r.Request(MethodHead, path, fn, filters...)
}

// OPTIONS 注册 OPTIONS 方法处理函数
func (r *Router) OPTIONS(path string, fn interface{}, filters ...Filter) *Mapper {
	return r.Request(MethodOptions, path, fn, filters...)
}
