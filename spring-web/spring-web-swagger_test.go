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

package SpringWeb_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-spring/go-spring-web/spring-web"
)

func TestSwagger(t *testing.T) {

	sw := SpringWeb.Swagger().
		SetDescription("This is a sample server Petstore server.").
		SetTitle("Swagger Example API").
		SetTermsOfService("http://swagger.io/terms/").
		SetContact("API Support", "http://www.swagger.io/support", "support@swagger.io").
		SetLicense("Apache 2.0", "http://www.apache.org/licenses/LICENSE-2.0.html").
		SetVersion("1.0").
		SetHost("petstore.swagger.io").
		SetBasePath("/v2").
		Path("/file/upload", SpringWeb.MethodPost, SpringWeb.NewOperation().
			SetDescription("Upload file").
			SetConsumes("multipart/form-data").
			SetProduces("application/json").
			SetSummary("Upload file").
			SetOperationId("file.upload").
			Param("file", "file", "this is a test file", "formData", true).
			Success(200, "ok", "string", "").
			Failure(400, "We need ID!!", "object", "#/definitions/web.APIError").
			Failure(404, "Can not find ID", "object", "#/definitions/web.APIError")).
		Definition("web.APIError", SpringWeb.NewDefinition("object").
			SetProperty("CreatedAt", "string", "date-time").
			SetProperty("ErrorCode", "integer").
			SetProperty("ErrorMessage", "string"))

	bytes, _ := json.MarshalIndent(sw, "", "  ")
	fmt.Println(string(bytes))
}
