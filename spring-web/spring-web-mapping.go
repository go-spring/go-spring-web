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
	"errors"
	"reflect"
)

// HandlerType Handler 的反射类型
var HandlerType = reflect.TypeOf((*Handler)(nil)).Elem()

// FnHandlerType fnHandler 的反射类型
var FnHandlerType = reflect.TypeOf((*fnHandler)(nil)).Elem()

// WebMapping 路由表，Spring-Web 使用的路由规则和 echo 完全相同，并对 gin 做了适配。
type WebMapping interface {
	// Mappers 返回映射器列表
	Mappers() map[string]*Mapper

	// AddMapper 添加一个 Mapper
	AddMapper(m *Mapper) *Mapper

	// Route 返回和 Mapping 绑定的路由分组
	Route(basePath string, filters ...Filter) *Router

	// Request 注册任意 HTTP 方法处理函数
	Request(method uint32, path string, fn interface{}, filters ...Filter) *Mapper

	// GET 注册 GET 方法处理函数
	GET(path string, fn interface{}, filters ...Filter) *Mapper

	// HandleGet 注册 GET 方法处理函数
	HandleGet(path string, fn Handler, filters ...Filter) *Mapper

	// GetMapping 注册 GET 方法处理函数
	GetMapping(path string, fn HandlerFunc, filters ...Filter) *Mapper

	// GetBinding 注册 GET 方法处理函数
	GetBinding(path string, fn interface{}, filters ...Filter) *Mapper

	// POST 注册 POST 方法处理函数
	POST(path string, fn interface{}, filters ...Filter) *Mapper

	// HandlePost 注册 POST 方法处理函数
	HandlePost(path string, fn Handler, filters ...Filter) *Mapper

	// PostMapping 注册 POST 方法处理函数
	PostMapping(path string, fn HandlerFunc, filters ...Filter) *Mapper

	// PostBinding 注册 POST 方法处理函数
	PostBinding(path string, fn interface{}, filters ...Filter) *Mapper

	// PATCH 注册 PATCH 方法处理函数
	PATCH(path string, fn interface{}, filters ...Filter) *Mapper

	// PUT 注册 PUT 方法处理函数
	PUT(path string, fn interface{}, filters ...Filter) *Mapper

	// DELETE 注册 DELETE 方法处理函数
	DELETE(path string, fn interface{}, filters ...Filter) *Mapper

	// HEAD 注册 HEAD 方法处理函数
	HEAD(path string, fn interface{}, filters ...Filter) *Mapper

	// OPTIONS 注册 OPTIONS 方法处理函数
	OPTIONS(path string, fn interface{}, filters ...Filter) *Mapper
}

// defaultWebMapping 路由表的默认实现
type defaultWebMapping struct {
	mappers map[string]*Mapper
}

// NewDefaultWebMapping defaultWebMapping 的构造函数
func NewDefaultWebMapping() *defaultWebMapping {
	return &defaultWebMapping{
		mappers: make(map[string]*Mapper),
	}
}

// Mappers 返回映射器列表
func (w *defaultWebMapping) Mappers() map[string]*Mapper {
	return w.mappers
}

// AddMapper 添加一个 Mapper
func (w *defaultWebMapping) AddMapper(m *Mapper) *Mapper {
	w.mappers[m.Key()] = m
	return m
}

// Route 返回和 Mapping 绑定的路由分组
func (w *defaultWebMapping) Route(basePath string, filters ...Filter) *Router {
	return &Router{mapping: w, basePath: basePath, filters: filters}
}

// Request 注册任意 HTTP 方法处理函数
func (w *defaultWebMapping) Request(method uint32, path string, fn interface{}, filters ...Filter) *Mapper {
	var v reflect.Value

	fnType := reflect.TypeOf(fn)
	fnValue := reflect.ValueOf(fn)

	if fnType.AssignableTo(FnHandlerType) { // 标准形式
		v = fnValue.Convert(FnHandlerType)
		v = v.Convert(HandlerType)

	} else if fnType.AssignableTo(HandlerType) {
		v = fnValue.Convert(HandlerType)

	} else {
		panic(errors.New("error func type " + fnType.String()))
	}

	h := v.Interface().(Handler)
	m := NewMapper(method, path, h, filters)
	w.mappers[m.Key()] = m
	return m
}

// GET 注册 GET 方法处理函数
func (w *defaultWebMapping) GET(path string, fn interface{}, filters ...Filter) *Mapper {
	return w.Request(MethodGet, path, fn, filters...)
}

// HandleGet 注册 GET 方法处理函数
func (w *defaultWebMapping) HandleGet(path string, fn Handler, filters ...Filter) *Mapper {
	return w.Request(MethodGet, path, fn, filters...)
}

// GetMapping 注册 GET 方法处理函数
func (w *defaultWebMapping) GetMapping(path string, fn HandlerFunc, filters ...Filter) *Mapper {
	return w.Request(MethodGet, path, FUNC(fn), filters...)
}

// GetBinding 注册 GET 方法处理函数
func (w *defaultWebMapping) GetBinding(path string, fn interface{}, filters ...Filter) *Mapper {
	return w.Request(MethodGet, path, BIND(fn), filters...)
}

// POST 注册 POST 方法处理函数
func (w *defaultWebMapping) POST(path string, fn interface{}, filters ...Filter) *Mapper {
	return w.Request(MethodPost, path, fn, filters...)
}

// HandlePost 注册 POST 方法处理函数
func (w *defaultWebMapping) HandlePost(path string, fn Handler, filters ...Filter) *Mapper {
	return w.Request(MethodPost, path, fn, filters...)
}

// PostMapping 注册 POST 方法处理函数
func (w *defaultWebMapping) PostMapping(path string, fn HandlerFunc, filters ...Filter) *Mapper {
	return w.Request(MethodPost, path, FUNC(fn), filters...)
}

// PostBinding 注册 POST 方法处理函数
func (w *defaultWebMapping) PostBinding(path string, fn interface{}, filters ...Filter) *Mapper {
	return w.Request(MethodPost, path, BIND(fn), filters...)
}

// PATCH 注册 PATCH 方法处理函数
func (w *defaultWebMapping) PATCH(path string, fn interface{}, filters ...Filter) *Mapper {
	return w.Request(MethodPatch, path, fn, filters...)
}

// PUT 注册 PUT 方法处理函数
func (w *defaultWebMapping) PUT(path string, fn interface{}, filters ...Filter) *Mapper {
	return w.Request(MethodPut, path, fn, filters...)
}

// DELETE 注册 DELETE 方法处理函数
func (w *defaultWebMapping) DELETE(path string, fn interface{}, filters ...Filter) *Mapper {
	return w.Request(MethodDelete, path, fn, filters...)
}

// HEAD 注册 HEAD 方法处理函数
func (w *defaultWebMapping) HEAD(path string, fn interface{}, filters ...Filter) *Mapper {
	return w.Request(MethodHead, path, fn, filters...)
}

// OPTIONS 注册 OPTIONS 方法处理函数
func (w *defaultWebMapping) OPTIONS(path string, fn interface{}, filters ...Filter) *Mapper {
	return w.Request(MethodOptions, path, fn, filters...)
}
