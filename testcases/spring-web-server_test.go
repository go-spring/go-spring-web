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
	"net/url"
	"testing"
	"time"

	"github.com/go-spring/go-spring-web/spring-echo"
	"github.com/go-spring/go-spring-web/spring-gin"
	"github.com/go-spring/go-spring-web/spring-web"
	"github.com/go-spring/go-spring-web/testcases"
)

func TestWebServer(t *testing.T) {

	l := list.New()
	f2 := testcases.NewNumberFilter(2, l)
	f5 := testcases.NewNumberFilter(5, l)
	f6 := testcases.NewNumberFilter(6, l)
	f7 := testcases.NewNumberFilter(7, l)

	server := SpringWeb.NewWebServer()
	server.SetFilters(f6)

	s := testcases.NewService()

	// 可用于全局的路由分组
	r := SpringWeb.NewRouter("/v1", nil)
	r.GET("/router", func(ctx SpringWeb.WebContext) {
		ctx.String(http.StatusOK, "router:ok")
	})

	// 添加第一个 Web 容器
	{
		g := SpringGin.NewContainer()
		server.AddWebContainer(g)
		g.SetPort(8080)
		g.AddRouter(r)

		g.GET("/get", s.Get, f5)
	}

	// 添加第二个 Web 容器
	{
		e := SpringEcho.NewContainer()
		server.AddWebContainer(e)
		e.SetPort(9090)
		e.AddRouter(r)

		r0 := e.Route("", f2, f7)
		{
			r0.POST("/set", s.Set)
			r0.GET("/panic", s.Panic)
		}
	}

	// 启动 web 服务器
	server.Start()

	time.Sleep(time.Millisecond * 100)
	fmt.Println()

	resp, _ := http.Get("http://127.0.0.1:8080/get?key=a")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
	fmt.Println()

	_, _ = http.PostForm("http://127.0.0.1:9090/set", url.Values{
		"a": []string{"1"},
	})

	fmt.Println()

	resp, _ = http.Get("http://127.0.0.1:8080/get?key=a")
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
	fmt.Println()

	resp, _ = http.Get("http://127.0.0.1:9090/panic")
	body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
	fmt.Println()

	server.Stop(context.TODO())

	time.Sleep(time.Millisecond * 50)
}
