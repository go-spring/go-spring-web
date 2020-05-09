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
	"container/list"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/go-spring/go-spring-web/spring-echo"
	"github.com/go-spring/go-spring-web/spring-gin"
	"github.com/go-spring/go-spring-web/spring-web"
	"github.com/go-spring/go-spring-web/testcases"
)

func TestRpc(t *testing.T) {

	rc := new(testcases.RpcService)

	l := list.New()
	f2 := testcases.NewNumberFilter(2, l)
	f5 := testcases.NewNumberFilter(5, l)
	f7 := testcases.NewNumberFilter(7, l)

	server := SpringWeb.NewWebServer()

	// 添加第一个 web 容器
	{
		c1 := SpringGin.NewContainer()
		server.AddContainer(c1)
		c1.SetPort(8080)

		c1.HandleGet("/ok", SpringWeb.RPC(rc.OK), f2, f5)
	}

	// 添加第二个 web 容器
	{
		c2 := SpringEcho.NewContainer()
		server.AddContainer(c2)
		c2.SetPort(9090)

		r := c2.Route("", f2, f7)
		{
			r.HandleGet("/err", SpringWeb.RPC(rc.Err))
			r.HandleGet("/panic", SpringWeb.RPC(rc.Panic))
		}

		c2.GetBinding("/echo", rc.Echo, f5)
		c2.GetBinding("/ptr_echo", rc.PtrEcho, f5)
		c2.GetBinding("/ctx_echo", rc.CtxEcho, f5)
		c2.GetBinding("/echo_ctx", rc.EchoCtx, f5)
	}

	// 启动 web 服务器
	server.Start()

	time.Sleep(time.Millisecond * 100)
	fmt.Println()

	resp, _ := http.Get("http://127.0.0.1:8080/ok")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body)+"\n")

	resp, _ = http.Get("http://127.0.0.1:9090/err")
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	resp, _ = http.Get("http://127.0.0.1:9090/panic")
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	resp, _ = http.Get("http://127.0.0.1:9090/panic?panic=1")
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	resp, _ = http.Get("http://127.0.0.1:9090/echo?str=echo")
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	resp, _ = http.Get("http://127.0.0.1:9090/ptr_echo?str=echo")
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	resp, _ = http.Get("http://127.0.0.1:9090/ctx_echo?str=echo")
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	resp, _ = http.Get("http://127.0.0.1:9090/echo_ctx?str=echo")
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	server.Stop(context.TODO())

	time.Sleep(50 * time.Millisecond)
}
