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
	s.Info.Contact = &spec.ContactInfo{Name: name, URL: url, Email: email}
	return s
}

// WithLicense 设置开源协议的名称、地址
func (s *swagger) WithLicense(name string, url string) *swagger {
	s.Info.License = &spec.License{Name: name, URL: url}
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

// AddTag 添加一个标签
func (s *swagger) AddTag(tag *spec.Tag) *swagger {
	s.Swagger.Tags = append(s.Swagger.Tags, *tag)
	return s
}

// AddPath 添加一个路由
func (s *swagger) AddPath(path string, method uint32, op *operation,
	parameters ...spec.Parameter) *swagger {

	pathItem := spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Parameters: parameters,
		},
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

// AddBasicSecurityDefinition 添加 Basic 方式认证
func (s *swagger) AddBasicSecurityDefinition() *swagger {
	s.Swagger.SecurityDefinitions["BasicAuth"] = spec.BasicAuth()
	return s
}

// AddApiKeySecurityDefinition 添加 ApiKey 方式认证
func (s *swagger) AddApiKeySecurityDefinition(name string, in string) *swagger {
	s.Swagger.SecurityDefinitions["ApiKeyAuth"] = spec.APIKeyAuth(name, in)
	return s
}

// AddOauth2ApplicationSecurityDefinition 添加 OAuth2 Application 方式认证
func (s *swagger) AddOauth2ApplicationSecurityDefinition(tokenUrl string, scopes map[string]string) *swagger {
	securityScheme := spec.OAuth2Application(tokenUrl)
	return s.securitySchemeWithScopes("OAuth2Application", securityScheme, scopes)
}

// AddOauth2ImplicitSecurityDefinition 添加 OAuth2 Implicit 方式认证
func (s *swagger) AddOauth2ImplicitSecurityDefinition(authorizationUrl string, scopes map[string]string) *swagger {
	securityScheme := spec.OAuth2Implicit(authorizationUrl)
	return s.securitySchemeWithScopes("OAuth2Implicit", securityScheme, scopes)
}

// AddOauth2PasswordSecurityDefinition 添加 OAuth2 Password 方式认证
func (s *swagger) AddOauth2PasswordSecurityDefinition(tokenUrl string, scopes map[string]string) *swagger {
	securityScheme := spec.OAuth2Password(tokenUrl)
	return s.securitySchemeWithScopes("OAuth2Password", securityScheme, scopes)
}

// AddOauth2AccessCodeSecurityDefinition 添加 OAuth2 AccessCode 方式认证
func (s *swagger) AddOauth2AccessCodeSecurityDefinition(authorizationUrl string, tokenUrl string, scopes map[string]string) *swagger {
	securityScheme := spec.OAuth2AccessToken(authorizationUrl, tokenUrl)
	return s.securitySchemeWithScopes("OAuth2AccessCode", securityScheme, scopes)
}

func (s *swagger) securitySchemeWithScopes(name string, scheme *spec.SecurityScheme, scopes map[string]string) *swagger {
	securityScheme := scheme
	for scope, description := range scopes {
		securityScheme.AddScope(scope, description)
	}
	s.Swagger.SecurityDefinitions[name] = securityScheme
	return s
}

// operation 封装 *spec.Operation 对象，提供更多功能
type operation struct {
	operation *spec.Operation
}

// NewOperation creates a new operation instance.
func NewOperation(id string) *operation {
	return &operation{operation: spec.NewOperation(id)}
}

// WithID sets the ID property on this operation, allows for chaining.
func (o *operation) WithID(id string) *operation {
	o.operation.WithID(id)
	return o
}

// WithDescription sets the description on this operation, allows for chaining
func (o *operation) WithDescription(description string) *operation {
	o.operation.WithDescription(description)
	return o
}

// WithSummary sets the summary on this operation, allows for chaining
func (o *operation) WithSummary(summary string) *operation {
	o.operation.WithSummary(summary)
	return o
}

// WithExternalDocs sets/removes the external docs for/from this operation.
func (o *operation) WithExternalDocs(description, url string) *operation {
	o.operation.WithExternalDocs(description, url)
	return o
}

// Deprecate marks the operation as deprecated
func (o *operation) Deprecate() *operation {
	o.operation.Deprecate()
	return o
}

// Undeprecate marks the operation as not deprecated
func (o *operation) Undeprecate() *operation {
	o.operation.Undeprecate()
	return o
}

// WithConsumes adds media types for incoming body values
func (o *operation) WithConsumes(mediaTypes ...string) *operation {
	o.operation.WithConsumes(mediaTypes...)
	return o
}

// WithProduces adds media types for outgoing body values
func (o *operation) WithProduces(mediaTypes ...string) *operation {
	o.operation.WithProduces(mediaTypes...)
	return o
}

// WithTags adds tags for this operation
func (o *operation) WithTags(tags ...string) *operation {
	o.operation.WithTags(tags...)
	return o
}

// SetSchemes 设置服务协议
func (o *operation) WithSchemes(schemes ...string) *operation {
	o.operation.Schemes = schemes
	return o
}

// AddParam adds a parameter to this operation
func (o *operation) AddParam(param *spec.Parameter) *operation {
	o.operation.AddParam(param)
	return o
}

// RemoveParam removes a parameter from the operation
func (o *operation) RemoveParam(name, in string) *operation {
	o.operation.RemoveParam(name, in)
	return o
}

// SecuredWith adds a security scope to this operation.
func (o *operation) SecuredWith(name string, scopes ...string) *operation {
	o.operation.SecuredWith(name, scopes...)
	return o
}

// WithDefaultResponse adds a default response to the operation.
func (o *operation) WithDefaultResponse(response *spec.Response) *operation {
	o.operation.WithDefaultResponse(response)
	return o
}

// RespondsWith adds a status code response to the operation.
func (o *operation) RespondsWith(code int, response *spec.Response) *operation {
	o.operation.RespondsWith(code, response)
	return o
}
