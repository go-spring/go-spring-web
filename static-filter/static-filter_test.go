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

package StaticFilter_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/go-spring/go-spring-web/spring-gin"
	"github.com/go-spring/go-spring-web/spring-web"
	"github.com/go-spring/go-spring-web/static-filter"
)

func TestStaticFilter(t *testing.T) {

	dir, _ := os.Getwd()
	dir = path.Dir(dir)

	c := SpringGin.NewContainer(SpringWeb.ContainerConfig{Port: 8080})
	c.GetMapping("/ok", func(ctx SpringWeb.WebContext) {
		ctx.String(http.StatusOK, "yes")
	})
	c.AddFilter(StaticFilter.New(dir, ""))
	c.Start()

	resp, _ := http.Get("http://127.0.0.1:8080/go.mod")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
	fmt.Println()

	resp, _ = http.Get("http://127.0.0.1:8080/ok")
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
	fmt.Println()
}
