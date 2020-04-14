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
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/go-openapi/spec"
	"github.com/swaggo/swag"
)

// doc 全局 swagger 对象
var doc = NewSwagger()

func init() {
	swag.Register(swag.Name, doc)
}

// Swagger 返回全局的 swagger 对象
func Swagger() *swagger {
	return doc
}

// swagger 封装 spec.Swagger 对象，提供流式调用
type swagger struct {
	spec.Swagger
}

// NewSwagger swagger 的构造函数
func NewSwagger() *swagger {
	return &swagger{
		Swagger: spec.Swagger{
			SwaggerProps: spec.SwaggerProps{
				Swagger: "2.0",
				Info: &spec.Info{
					InfoProps: spec.InfoProps{
						Contact: &spec.ContactInfo{},
						License: &spec.License{},
					},
				},
				Paths: &spec.Paths{
					Paths: make(map[string]spec.PathItem),
				},
				Definitions:         make(map[string]spec.Schema),
				SecurityDefinitions: map[string]*spec.SecurityScheme{},
			},
		},
	}
}

// ReadDoc 获取应用的 Swagger 描述内容
func (s *swagger) ReadDoc() string {
	if b, err := s.MarshalJSON(); err == nil {
		return string(b)
	} else {
		panic(err)
	}
}

// WithID 设置应用 ID
func (s *swagger) WithID(id string) *swagger {
	s.ID = id
	return s
}

// WithConsumes 设置消费协议
func (s *swagger) WithConsumes(consumes ...string) *swagger {
	s.Consumes = consumes
	return s
}

// WithProduces 设置生产协议
func (s *swagger) WithProduces(produces ...string) *swagger {
	s.Produces = produces
	return s
}

// WithSchemes 设置服务协议
func (s *swagger) WithSchemes(schemes ...string) *swagger {
	s.Schemes = schemes
	return s
}

// WithDescription 设置服务描述
func (s *swagger) WithDescription(description string) *swagger {
	s.Info.Description = description
	return s
}

// WithTitle 设置服务名称
func (s *swagger) WithTitle(title string) *swagger {
	s.Info.Title = title
	return s
}

// WithTermsOfService 设置服务条款地址
func (s *swagger) WithTermsOfService(termsOfService string) *swagger {
	s.Info.TermsOfService = termsOfService
	return s
}

// WithContact 设置作者的名字、主页地址、邮箱
func (s *swagger) WithContact(name string, url string, email string) *swagger {
	c := new(spec.ContactInfo)
	c.Name = name
	c.URL = url
	c.Email = email
	s.Info.Contact = c
	return s
}

// WithLicense 设置开源协议的名称、地址
func (s *swagger) WithLicense(name string, url string) *swagger {
	l := new(spec.License)
	l.Name = name
	l.URL = url
	s.Info.License = l
	return s
}

// WithVersion 设置 API 版本号
func (s *swagger) WithVersion(version string) *swagger {
	s.Info.Version = version
	return s
}

// WithHost 设置可用服务器地址
func (s *swagger) WithHost(host string) *swagger {
	s.Host = host
	return s
}

// WithBasePath 设置 API 路径的前缀
func (s *swagger) WithBasePath(basePath string) *swagger {
	s.BasePath = basePath
	return s
}

// WithTags 添加标签
func (s *swagger) WithTags(tags ...spec.Tag) *swagger {
	s.Swagger.Tags = tags
	return s
}

// AddPath 添加一个路由
func (s *swagger) AddPath(path string, method uint32, op *Operation,
	parameters ...spec.Parameter) *swagger {

	path = strings.TrimPrefix(path, doc.BasePath)
	path = strings.TrimRight(path, "/")
	pathItem, ok := s.Paths.Paths[path]

	if !ok {
		pathItem = spec.PathItem{
			PathItemProps: spec.PathItemProps{
				Parameters: parameters,
			},
		}
	}

	for _, m := range GetMethod(method) {
		switch m {
		case http.MethodGet:
			pathItem.Get = op.operation
		case http.MethodPost:
			pathItem.Post = op.operation
		case http.MethodPut:
			pathItem.Put = op.operation
		case http.MethodDelete:
			pathItem.Delete = op.operation
		case http.MethodOptions:
			pathItem.Options = op.operation
		case http.MethodHead:
			pathItem.Head = op.operation
		case http.MethodPatch:
			pathItem.Patch = op.operation
		}
	}

	s.Paths.Paths[path] = pathItem
	return s
}

// AddDefinition 添加一个定义
func (s *swagger) AddDefinition(name string, schema *spec.Schema) *swagger {
	s.Definitions[name] = *schema
	return s
}

type DefinitionField struct {
	Description string
	Example     interface{}
	Enums       []interface{}
}

// BindDefinitions 绑定一个定义
func (s *swagger) BindDefinitions(i ...interface{}) *swagger {
	m := map[string]DefinitionField{}
	for _, v := range i {
		s.BindDefinitionWithTags(v, m)
	}
	return s
}

// BindDefinitionWithTags 绑定一个定义
func (s *swagger) BindDefinitionWithTags(i interface{}, attachFields map[string]DefinitionField) *swagger {

	it := reflect.TypeOf(i)
	if it.Kind() == reflect.Ptr {
		it = it.Elem()
	}

	objSchema := new(spec.Schema).Typed("object", "")
	for i := 0; i < it.NumField(); i++ {
		f := it.Field(i)

		// 处理 XML 标签
		var xmlTag []string
		if tag, ok := f.Tag.Lookup("xml"); ok {
			xmlTag = strings.Split(tag, ",")
			if f.Type == reflect.TypeOf(xml.Name{}) {
				objSchema.WithXMLName(xmlTag[0])
				continue
			}
		}

		propName := f.Name

		// 处理 JSON 标签
		var jsonTag []string
		if tag, ok := f.Tag.Lookup("json"); ok {
			jsonTag = strings.Split(tag, ",")
			propName = jsonTag[0]
		}

		var propSchema *spec.Schema
		switch k := f.Type.Kind(); k {
		case reflect.Bool:
			propSchema = spec.BoolProperty()
		case reflect.Int8:
			propSchema = spec.Int8Property()
		case reflect.Int16:
			propSchema = spec.Int16Property()
		case reflect.Int32:
			propSchema = spec.Int32Property()
		case reflect.Int64:
			propSchema = spec.Int64Property()
		case reflect.String:
			propSchema = spec.StringProperty()
		case reflect.Struct:
			if f.Type == reflect.TypeOf(time.Time{}) {
				propSchema = spec.DateTimeProperty()
			} else {
				panic(fmt.Errorf("unsupported swagger type %s", f.Type))
			}
		case reflect.Ptr:
			if et := f.Type.Elem(); et.Kind() == reflect.Struct {
				propSchema = spec.RefSchema("#/definitions/" + et.Name())
			} else {
				panic(fmt.Errorf("unsupported swagger type %s", f.Type))
			}
		case reflect.Slice:
			{
				et := f.Type.Elem()

				var items *spec.Schema
				switch k := et.Kind(); k {
				case reflect.Bool:
					items = spec.BoolProperty()
				case reflect.Int8:
					items = spec.Int8Property()
				case reflect.Int16:
					items = spec.Int16Property()
				case reflect.Int32:
					items = spec.Int32Property()
				case reflect.Int64:
					items = spec.Int64Property()
				case reflect.String:
					items = spec.StringProperty()
				case reflect.Struct:
					items = spec.RefSchema("#/definitions/" + et.Name())
				default:
					panic(fmt.Errorf("unsupported swagger type %s", f.Type))
				}

				if len(xmlTag) > 0 {
					items.WithXMLName(xmlTag[0])
				}

				propSchema = spec.ArrayProperty(items)
			}
		default:
			panic(fmt.Errorf("unsupported swagger type %s", f.Type))
		}

		if len(xmlTag) > 1 {
			for _, v := range xmlTag {
				if v == "wrapped" {
					propSchema.AsWrappedXML()
					break
				}
			}
		}

		required := true

		for _, v := range jsonTag {
			if v == "omitempty" {
				required = false
				break
			}
		}

		if required {
			objSchema.AddRequired(propName)
		}

		if attachField, ok := attachFields[propName]; ok {
			if len(attachField.Enums) > 0 {
				propSchema.WithEnum(attachField.Enums...)
			}
			if attachField.Description != "" {
				propSchema.WithDescription(attachField.Description)
			}
			if attachField.Example != "" {
				propSchema.WithExample(attachField.Example)
			}
		}

		objSchema.SetProperty(propName, *propSchema)
	}

	s.Definitions[it.Name()] = *objSchema
	return s
}

// AddBasicSecurityDefinition 添加 Basic 方式认证
func (s *swagger) AddBasicSecurityDefinition() *swagger {
	s.Swagger.SecurityDefinitions["BasicAuth"] = spec.BasicAuth()
	return s
}

// AddApiKeySecurityDefinition 添加 ApiKey 方式认证
func (s *swagger) AddApiKeySecurityDefinition(name string, in string) *swagger {
	if name == "" {
		name = "ApiKeyAuth"
	}
	s.Swagger.SecurityDefinitions[name] = spec.APIKeyAuth(name, in)
	return s
}

// AddOauth2ApplicationSecurityDefinition 添加 OAuth2 Application 方式认证
func (s *swagger) AddOauth2ApplicationSecurityDefinition(name string, tokenUrl string, scopes map[string]string) *swagger {
	if name == "" {
		name = "OAuth2Application"
	}
	securityScheme := spec.OAuth2Application(tokenUrl)
	return s.securitySchemeWithScopes(name, securityScheme, scopes)
}

// AddOauth2ImplicitSecurityDefinition 添加 OAuth2 Implicit 方式认证
func (s *swagger) AddOauth2ImplicitSecurityDefinition(name string, authorizationUrl string, scopes map[string]string) *swagger {
	if name == "" {
		name = "OAuth2Implicit"
	}
	securityScheme := spec.OAuth2Implicit(authorizationUrl)
	return s.securitySchemeWithScopes(name, securityScheme, scopes)
}

// AddOauth2PasswordSecurityDefinition 添加 OAuth2 Password 方式认证
func (s *swagger) AddOauth2PasswordSecurityDefinition(name string, tokenUrl string, scopes map[string]string) *swagger {
	if name == "" {
		name = "OAuth2Password"
	}
	securityScheme := spec.OAuth2Password(tokenUrl)
	return s.securitySchemeWithScopes(name, securityScheme, scopes)
}

// AddOauth2AccessCodeSecurityDefinition 添加 OAuth2 AccessCode 方式认证
func (s *swagger) AddOauth2AccessCodeSecurityDefinition(name string, authorizationUrl string, tokenUrl string, scopes map[string]string) *swagger {
	if name == "" {
		name = "OAuth2AccessCode"
	}
	securityScheme := spec.OAuth2AccessToken(authorizationUrl, tokenUrl)
	return s.securitySchemeWithScopes(name, securityScheme, scopes)
}

func (s *swagger) securitySchemeWithScopes(name string, scheme *spec.SecurityScheme, scopes map[string]string) *swagger {
	securityScheme := scheme
	for scope, description := range scopes {
		securityScheme.AddScope(scope, description)
	}
	s.Swagger.SecurityDefinitions[name] = securityScheme
	return s
}

// WithExternalDocs
func (s *swagger) WithExternalDocs(externalDocs *spec.ExternalDocumentation) *swagger {
	s.Swagger.ExternalDocs = externalDocs
	return s
}

type bindParam struct {
	param       interface{}
	description string
}

// Operation 封装 *spec.Operation 对象，提供更多功能
type Operation struct {
	operation *spec.Operation
	bindParam *bindParam
}

// NewOperation creates a new operation instance.
func NewOperation(id string) *Operation {
	return &Operation{operation: spec.NewOperation(id)}
}

// WithID sets the ID property on this operation, allows for chaining.
func (o *Operation) WithID(id string) *Operation {
	o.operation.WithID(id)
	return o
}

// WithDescription sets the description on this operation, allows for chaining
func (o *Operation) WithDescription(description string) *Operation {
	o.operation.WithDescription(description)
	return o
}

// WithSummary sets the summary on this operation, allows for chaining
func (o *Operation) WithSummary(summary string) *Operation {
	o.operation.WithSummary(summary)
	return o
}

// WithExternalDocs sets/removes the external docs for/from this operation.
func (o *Operation) WithExternalDocs(description, url string) *Operation {
	o.operation.WithExternalDocs(description, url)
	return o
}

// Deprecate marks the operation as deprecated
func (o *Operation) Deprecate() *Operation {
	o.operation.Deprecate()
	return o
}

// Undeprecate marks the operation as not deprecated
func (o *Operation) Undeprecate() *Operation {
	o.operation.Undeprecate()
	return o
}

// WithConsumes adds media types for incoming body values
func (o *Operation) WithConsumes(mediaTypes ...string) *Operation {
	o.operation.WithConsumes(mediaTypes...)
	return o
}

// WithProduces adds media types for outgoing body values
func (o *Operation) WithProduces(mediaTypes ...string) *Operation {
	o.operation.WithProduces(mediaTypes...)
	return o
}

// WithTags adds tags for this operation
func (o *Operation) WithTags(tags ...string) *Operation {
	o.operation.WithTags(tags...)
	return o
}

// SetSchemes 设置服务协议
func (o *Operation) WithSchemes(schemes ...string) *Operation {
	o.operation.Schemes = schemes
	return o
}

// AddParam adds a parameter to this operation
func (o *Operation) AddParam(param *spec.Parameter) *Operation {
	o.operation.AddParam(param)
	return o
}

// RemoveParam removes a parameter from the operation
func (o *Operation) RemoveParam(name, in string) *Operation {
	o.operation.RemoveParam(name, in)
	return o
}

// SecuredWith adds a security scope to this operation.
func (o *Operation) SecuredWith(name string, scopes ...string) *Operation {
	o.operation.SecuredWith(name, scopes...)
	return o
}

// WithDefaultResponse adds a default response to the operation.
func (o *Operation) WithDefaultResponse(response *spec.Response) *Operation {
	o.operation.WithDefaultResponse(response)
	return o
}

// RespondsWith adds a status code response to the operation.
func (o *Operation) RespondsWith(code int, response *spec.Response) *Operation {
	o.operation.RespondsWith(code, response)
	return o
}

// Bind 绑定请求参数
func (o *Operation) BindParam(i interface{}, description string) *Operation {
	o.bindParam = &bindParam{param: i, description: description}
	return o
}

// parseBind 解析绑定的请求参数
func (o *Operation) parseBind() error {
	if o.bindParam != nil && o.bindParam.param != nil {
		t := reflect.TypeOf(o.bindParam.param)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.Struct {
			schema := spec.RefSchema("#/definitions/" + t.Name())
			param := BodyParam("body", schema).
				WithDescription(o.bindParam.description).
				AsRequired()
			o.AddParam(param)
		}
	}
	return nil
}

// HeaderParam creates a header parameter, this is always required by default
func HeaderParam(name string, typ, format string) *spec.Parameter {
	param := spec.HeaderParam(name)
	param.Typed(typ, format)
	return param
}

// PathParam creates a path parameter, this is always required
func PathParam(name string, typ, format string) *spec.Parameter {
	param := spec.PathParam(name)
	param.Typed(typ, format)
	return param
}

// BodyParam creates a body parameter
func BodyParam(name string, schema *spec.Schema) *spec.Parameter {
	return &spec.Parameter{ParamProps: spec.ParamProps{Name: name, In: "body", Schema: schema}}
}

// NewResponse creates a new response instance
func NewResponse(description string) *spec.Response {
	resp := new(spec.Response)
	resp.Description = description
	return resp
}

// NewBindResponse creates a new response instance
func NewBindResponse(i interface{}, description string) *spec.Response {
	resp := new(spec.Response)
	resp.Description = description

	t := reflect.TypeOf(i)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	slice := false
	if t.Kind() == reflect.Slice {
		slice = true
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		schema := spec.RefSchema("#/definitions/" + t.Name())
		if slice {
			resp.WithSchema(spec.ArrayProperty(schema))
		} else {
			resp.WithSchema(schema)
		}
	}

	return resp
}
