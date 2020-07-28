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

package testcases_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/go-spring/go-spring-web/spring-echo"
	"github.com/go-spring/go-spring-web/spring-web"
	"github.com/go-spring/go-spring-web/testcases"
)

func TestRpc(t *testing.T) {

	cfg := SpringWeb.ContainerConfig{Port: 9090}
	c2 := SpringEcho.NewContainer(cfg)

	server := SpringWeb.NewWebServer()
	server.AddContainer(c2)

	rc := new(testcases.RpcService)
	c2.GetBinding("/echo", rc.Echo)

	// 启动 web 服务器
	server.Start()

	time.Sleep(time.Millisecond * 100)
	fmt.Println()

	resp, _ := http.Get("http://127.0.0.1:9090/echo?str=echo")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	server.Stop(context.TODO())

	time.Sleep(50 * time.Millisecond)
}
