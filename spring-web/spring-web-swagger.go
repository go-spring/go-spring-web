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

// doc
var doc = &swaggerDoc{newSwagger()}

func init() {
	swag.Register(swag.Name, doc)
}

// Swagger
func Swagger() *swagger {
	return doc.swagger
}

// swaggerDoc
type swaggerDoc struct {
	swagger *swagger
}

// ReadDoc
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
	Name string `json:"name"`
	Url  string `json:"url"`
}

// path
type path map[string]*operation

// property
type property struct {
	Type   string `json:"type"`
	Format string `json:"format,omitempty"`
}

// definition
type definition struct {
	Type       string              `json:"type"`
	Properties map[string]property `json:"properties,omitempty"`
}

// NewDefinition
func NewDefinition(typ string) *definition {
	return &definition{Type: typ, Properties: make(map[string]property)}
}

// SetProperty
func (d *definition) SetProperty(name string, typ string, format ...string) *definition {
	p := property{Type: typ}
	if len(format) > 0 {
		p.Format = format[0]
	}
	d.Properties[name] = p
	return d
}

// swagger
type swagger struct {
	Swagger string `json:"swagger"`
	Info    struct {
		Description    string   `json:"description"`
		Title          string   `json:"title"`
		TermsOfService string   `json:"termsOfService"`
		Contact        *contact `json:"contact,omitempty"`
		License        *license `json:"license,omitempty"`
		Version        string   `json:"version"`
	} `json:"info"`
	Host        string                 `json:"host"`
	BasePath    string                 `json:"basePath"`
	Paths       map[string]path        `json:"paths"`
	Definitions map[string]*definition `json:"definitions,omitempty"`
}

// newSwagger
func newSwagger() *swagger {
	return &swagger{
		Swagger:     "2.0",
		Paths:       make(map[string]path),
		Definitions: make(map[string]*definition),
	}
}

// SetDescription
func (s *swagger) SetDescription(description string) *swagger {
	s.Info.Description = description
	return s
}

// SetTitle
func (s *swagger) SetTitle(title string) *swagger {
	s.Info.Title = title
	return s
}

// SetTermsOfService
func (s *swagger) SetTermsOfService(termsOfService string) *swagger {
	s.Info.TermsOfService = termsOfService
	return s
}

// SetContact
func (s *swagger) SetContact(name string, url string, email string) *swagger {
	s.Info.Contact = &contact{Name: name, Url: url, Email: email}
	return s
}

// SetLicense
func (s *swagger) SetLicense(name string, url string) *swagger {
	s.Info.License = &license{Name: name, Url: url}
	return s
}

// SetVersion
func (s *swagger) SetVersion(version string) *swagger {
	s.Info.Version = version
	return s
}

// SetHost
func (s *swagger) SetHost(host string) *swagger {
	s.Host = host
	return s
}

// SetBasePath
func (s *swagger) SetBasePath(basePath string) *swagger {
	s.BasePath = basePath
	return s
}

// Path
func (s *swagger) Path(path string, method uint32, op *operation) *swagger {
	pathOperation := make(map[string]*operation)
	for _, m := range GetMethod(method) {
		pathOperation[strings.ToLower(m)] = op
	}
	s.Paths[path] = pathOperation
	return s
}

// Definition
func (s *swagger) Definition(name string, d *definition) *swagger {
	s.Definitions[name] = d
	return s
}

// response
type response struct {
	Description string `json:"description"`
	Schema      struct {
		Type string `json:"type"`
		Ref  string `json:"$ref,omitempty"`
	} `json:"schema"`
}

// newResponse
func newResponse(description string, typ string, ref string) *response {
	r := &response{}
	r.Description = description
	r.Schema.Type = typ
	r.Schema.Ref = ref
	return r
}

// parameter
type parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	In          string `json:"in"`
	Required    bool   `json:"required"`
}

// newParameter
func newParameter(name string, typ string, description string, in string, required bool) parameter {
	return parameter{Name: name, Type: typ, Description: description, In: in, Required: required}
}

// operation
type operation struct {
	Summary     string               `json:"summary"`
	Description string               `json:"description"`
	OperationId string               `json:"operationId"`
	Consumes    []string             `json:"consumes,omitempty"`
	Produces    []string             `json:"produces,omitempty"`
	Parameters  []parameter          `json:"parameters,omitempty"`
	Responses   map[string]*response `json:"responses,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
}

// NewOperation
func NewOperation() *operation {
	return &operation{Responses: make(map[string]*response)}
}

// SetSummary
func (s *operation) SetSummary(summary string) *operation {
	s.Summary = summary
	return s
}

// SetDescription
func (s *operation) SetDescription(description string) *operation {
	s.Description = description
	return s
}

// SetOperationId
func (s *operation) SetOperationId(operationId string) *operation {
	s.OperationId = operationId
	return s
}

// SetTags
func (s *operation) SetTags(tags ...string) *operation {
	s.Tags = tags
	return s
}

// SetConsumes
func (s *operation) SetConsumes(consumes ...string) *operation {
	s.Consumes = consumes
	return s
}

// SetProduces
func (s *operation) SetProduces(produces ...string) *operation {
	s.Produces = produces
	return s
}

// Parameter
func (s *operation) Parameter(name string, typ string, description string, in string, required bool) *operation {
	p := newParameter(name, typ, description, in, required)
	s.Parameters = append(s.Parameters, p)
	return s
}

// Success
func (s *operation) Success(code int, description string, typ string, ref string) *operation {
	s.Responses[strconv.Itoa(code)] = newResponse(description, typ, ref)
	return s
}

// Failure
func (s *operation) Failure(code int, description string, typ string, ref string) *operation {
	s.Responses[strconv.Itoa(code)] = newResponse(description, typ, ref)
	return s
}
