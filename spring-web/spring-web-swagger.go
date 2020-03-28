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
	"strconv"
	"strings"

	"github.com/go-spring/go-spring-parent/spring-utils"
	"github.com/swaggo/swag"
)

var doc = &swaggerDoc{newSwagger()}

func init() {
	swag.Register(swag.Name, doc)
}

// Swagger 返回全局的 swagger 对象
func Swagger() *swagger {
	return doc.swagger
}

// swaggerDoc 实现 swag.Swagger 接口
type swaggerDoc struct {
	swagger *swagger
}

// ReadDoc 获取应用的 Swagger 描述内容
func (s *swaggerDoc) ReadDoc() string {
	return SpringUtils.ToJson(s.swagger)
}

// contact 联系人信息
type contact struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Email string `json:"email"`
}

// license 版权信息
type license struct {
	Name string `json:"name"` // Required.
	Url  string `json:"url"`
}

type path map[string]*Operation

// property 数据类型的属性
type property struct {
	Type   string `json:"type"`
	Format string `json:"format,omitempty"`
}

// definition 数据类型定义
type definition struct {
	Type       string              `json:"type"`
	Properties map[string]property `json:"properties,omitempty"`
}

// NewDefinition definition 的构造函数
func NewDefinition(typ string) *definition {
	return &definition{Type: typ, Properties: make(map[string]property)}
}

// SetProperty 设置数据类型的属性
func (d *definition) SetProperty(name string, typ string, format ...string) *definition {
	p := property{Type: typ}
	if len(format) > 0 {
		p.Format = format[0]
	}
	d.Properties[name] = p
	return d
}

// swagger Swagger V2.0 文档
type swagger struct {
	Swagger string `json:"swagger"` // Swagger 版本号
	Info    struct {
		Description    string   `json:"description"`
		Title          string   `json:"title"` // Required.
		TermsOfService string   `json:"termsOfService"`
		Contact        *contact `json:"contact,omitempty"`
		License        *license `json:"license,omitempty"`
		Version        string   `json:"version"` // Required.
	} `json:"info"`
	Schemes     []string               `json:"schemes"`
	Host        string                 `json:"host"`
	BasePath    string                 `json:"basePath"`
	Paths       map[string]path        `json:"paths"`
	Definitions map[string]*definition `json:"definitions,omitempty"`
}

// newSwagger swagger 的构造函数
func newSwagger() *swagger {
	return &swagger{
		Swagger:     "2.0",
		Paths:       make(map[string]path),
		Definitions: make(map[string]*definition),
	}
}

// SetDescription 设置服务描述
func (s *swagger) SetDescription(description string) *swagger {
	s.Info.Description = description
	return s
}

// SetTitle 设置服务名称
func (s *swagger) SetTitle(title string) *swagger {
	s.Info.Title = title
	return s
}

// SetTermsOfService 设置服务条款地址
func (s *swagger) SetTermsOfService(termsOfService string) *swagger {
	s.Info.TermsOfService = termsOfService
	return s
}

// SetContact 设置作者的名字、主页地址、邮箱
func (s *swagger) SetContact(name string, url string, email string) *swagger {
	s.Info.Contact = &contact{Name: name, Url: url, Email: email}
	return s
}

// SetLicense 设置开源协议的名称、地址
func (s *swagger) SetLicense(name string, url string) *swagger {
	s.Info.License = &license{Name: name, Url: url}
	return s
}

// SetVersion 设置 API 版本号
func (s *swagger) SetVersion(version string) *swagger {
	s.Info.Version = version
	return s
}

// SetSchemes 设置服务协议
func (s *swagger) SetSchemes(schemes ...string) *swagger {
	s.Schemes = schemes
	return s
}

// SetHost 设置可用服务器地址
func (s *swagger) SetHost(host string) *swagger {
	s.Host = host
	return s
}

// SetBasePath 设置 API 路径的前缀
func (s *swagger) SetBasePath(basePath string) *swagger {
	s.BasePath = basePath
	return s
}

// Path 添加一个路由
func (s *swagger) Path(path string, method uint32, op *Operation) *swagger {
	pathOperation := make(map[string]*Operation)
	for _, m := range GetMethod(method) {
		pathOperation[strings.ToLower(m)] = op
	}
	s.Paths[path] = pathOperation
	return s
}

// Definition 添加一个数据类型
func (s *swagger) Definition(name string, d *definition) *swagger {
	s.Definitions[name] = d
	return s
}

// schema 描述(自定义)数据类型
type schema struct {
	Type string `json:"type"`
	Ref  string `json:"$ref,omitempty"`
}

// response 方法返回参数的描述
type response struct {
	Description string `json:"description"`
	Schema      schema `json:"schema"`
}

// newResponse response 的构造函数
func newResponse(description string, typ string, ref string) *response {
	r := &response{}
	r.Description = description
	r.Schema.Type = typ
	r.Schema.Ref = ref
	return r
}

// parameter 方法接收参数的描述
type parameter struct {
	Name        string  `json:"name"`
	Type        string  `json:"type,omitempty"`
	Description string  `json:"description"`
	In          string  `json:"in"`
	Required    bool    `json:"required"`
	Schema      *schema `json:"schema,omitempty"`
}

// newParameter parameter 的构造函数
func newParameter(name string, typ string, description string, in string, required bool) parameter {
	return parameter{
		Name:        name,
		Type:        typ,
		Description: description,
		In:          in,
		Required:    required,
	}
}

// newSchemaParameter parameter 的构造函数
func newSchemaParameter(name string, typ string, ref string, description string, in string, required bool) parameter {
	return parameter{
		Name:        name,
		Description: description,
		In:          in,
		Required:    required,
		Schema:      &schema{Type: typ, Ref: ref},
	}
}

// Operation 描述一个接口方法
type Operation struct {
	Summary     string               `json:"summary"`
	Description string               `json:"description"`
	OperationId string               `json:"operationId"`
	Consumes    []string             `json:"consumes,omitempty"`
	Produces    []string             `json:"produces,omitempty"`
	Parameters  []parameter          `json:"parameters,omitempty"`
	Responses   map[string]*response `json:"responses,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
}

// NewOperation Operation 的构造函数
func NewOperation() *Operation {
	return &Operation{Responses: make(map[string]*response)}
}

// SetSummary 设置方法简介
func (s *Operation) SetSummary(summary string) *Operation {
	s.Summary = summary
	return s
}

// SetDescription 设置方法的详细描述
func (s *Operation) SetDescription(description string) *Operation {
	s.Description = description
	return s
}

// SetOperationId 设置方法的唯一ID
func (s *Operation) SetOperationId(operationId string) *Operation {
	s.OperationId = operationId
	return s
}

// SetTags 设置方法的标签
func (s *Operation) SetTags(tags ...string) *Operation {
	s.Tags = tags
	return s
}

// SetConsumes 设置方法接收参数的类型
func (s *Operation) SetConsumes(consumes ...string) *Operation {
	s.Consumes = consumes
	return s
}

// SetProduces 设置方法返回参数的类型
func (s *Operation) SetProduces(produces ...string) *Operation {
	s.Produces = produces
	return s
}

// BindParam 从 bind 对象中获取接收参数
func (s *Operation) BindParam(name string, obj interface{}) {

}

// Param 添加一个接收参数
func (s *Operation) Param(name string, typ string, description string, in string, required bool) *Operation {
	p := newParameter(name, typ, description, in, required)
	s.Parameters = append(s.Parameters, p)
	return s
}

// Query 添加一个 query 类型的接收参数
func (s *Operation) Query(name string, typ string, description string, required bool) *Operation {
	return s.Param(name, typ, description, "query", required)
}

// Path 添加一个 path 类型的接收参数
func (s *Operation) Path(name string, typ string, description string, required bool) *Operation {
	return s.Param(name, typ, description, "path", required)
}

// Header 添加一个 header 类型的接收参数
func (s *Operation) Header(name string, typ string, description string, required bool) *Operation {
	return s.Param(name, typ, description, "header", required)
}

// Body 添加一个 body 类型的接收参数
func (s *Operation) Body(name string, typ string, description string, required bool) *Operation {
	return s.Param(name, typ, description, "body", required)
}

// Object 添加一个 body 类型的接收参数
func (s *Operation) Object(name string, typ string, ref string, description string, required bool) *Operation {
	p := newSchemaParameter(name, typ, ref, description, "body", required)
	s.Parameters = append(s.Parameters, p)
	return s
}

// FormData 添加一个 formData 类型的接收参数
func (s *Operation) FormData(name string, typ string, description string, required bool) *Operation {
	return s.Param(name, typ, description, "formData", required)
}

// Success 设置成功返回值
func (s *Operation) Success(code int, description string, typ string, ref string) *Operation {
	s.Responses[strconv.Itoa(code)] = newResponse(description, typ, ref)
	return s
}

// Failure 添加一个失败返回值
func (s *Operation) Failure(code int, description string, typ string, ref string) *Operation {
	s.Responses[strconv.Itoa(code)] = newResponse(description, typ, ref)
	return s
}
