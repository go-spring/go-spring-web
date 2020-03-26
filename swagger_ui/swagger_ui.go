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

package main

import (
	"time"

	"github.com/go-spring/go-spring-web/spring-echo"
	"github.com/go-spring/go-spring-web/spring-web"
	"github.com/swaggo/http-swagger"
)

func main() {

	SpringWeb.Swagger().
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
			Parameter("file", "file", "this is a test file", "formData", true).
			Success(200, "ok", "string", "").
			Failure(400, "We need ID!!", "object", "#/definitions/web.APIError").
			Failure(404, "Can not find ID", "object", "#/definitions/web.APIError")).
		Definition("web.APIError", SpringWeb.NewDefinition("object").
			SetProperty("CreatedAt", "string", "date-time").
			SetProperty("ErrorCode", "integer").
			SetProperty("ErrorMessage", "string"))

	c := SpringEcho.NewContainer()
	c.SetPort(1323)

	c.GET("/swagger/*", SpringWeb.HTTP(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:1323/swagger/doc.json"),
	)))

	c.Start()
	time.Sleep(time.Hour)
}
