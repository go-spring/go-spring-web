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

// Mapper 路由映射器
type Mapper struct {
	method  string   // 方法
	path    string   // 路径
	handler Handler  // 处理函数
	filters []Filter // 过滤器列表
}

// NewMapper Mapper 的构造函数
func NewMapper(method string, path string, fn Handler, filters []Filter) *Mapper {
	return &Mapper{
		method:  method,
		path:    path,
		handler: fn,
		filters: filters,
	}
}

// Key 返回 Mapper 的标识符
func (m *Mapper) Key() string {
	return m.method + m.path
}

// Method 返回 Mapper 的方法
func (m *Mapper) Method() string {
	return m.method
}

// Path 返回 Mapper 的路径
func (m *Mapper) Path() string {
	return m.path
}

// Handler 返回 Mapper 的处理函数
func (m *Mapper) Handler() Handler {
	return m.handler
}

// Filters 返回 Mapper 的过滤器列表
func (m *Mapper) Filters() []Filter {
	return m.filters
}

// Filters 设置 Mapper 的过滤器列表
func (m *Mapper) SetFilters(filters []Filter) *Mapper {
	m.filters = filters
	return m
}
