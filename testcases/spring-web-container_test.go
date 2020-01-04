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

	"github.com/gin-gonic/gin"
	"github.com/go-spring/go-spring-web/spring-echo"
	"github.com/go-spring/go-spring-web/spring-gin"
	"github.com/go-spring/go-spring-web/spring-web"
	"github.com/go-spring/go-spring-web/testcases"
	"github.com/labstack/echo"
)

func TestWebContainer(t *testing.T) {

	testRun := func(c SpringWeb.WebContainer) {
		c.SetPort(8080)

		l := list.New()
		f2 := testcases.NewNumberFilter(2, l)
		f5 := testcases.NewNumberFilter(5, l)
		f7 := testcases.NewNumberFilter(7, l)

		s := testcases.NewService()

		c.GET("/get", s.Get, f5)

		// 障眼法
		r := c.Route("", f2, f7)
		{
			r.POST("/set", s.Set)
			r.GET("/panic", s.Panic)
		}

		// 启动 web 服务器
		c.Start()

		time.Sleep(time.Millisecond * 100)
		fmt.Println()

		resp, _ := http.Get("http://127.0.0.1:8080/get?key=a")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		_, _ = http.PostForm("http://127.0.0.1:8080/set", url.Values{
			"a": []string{"1"},
		})

		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/get?key=a")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/panic")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		c.Stop(context.TODO())

		time.Sleep(time.Millisecond * 50)
	}

	t.Run("SpringGin", func(t *testing.T) {
		testRun(SpringGin.NewContainer())
	})

	t.Run("SpringEcho", func(t *testing.T) {
		testRun(SpringEcho.NewContainer())
	})
}

func TestEchoServer(t *testing.T) {
	e := echo.New()
	e.HideBanner = true

	// 配合 echo 框架使用
	e.GET("/", SpringEcho.HandlerWrapper(func(webCtx SpringWeb.WebContext) {
		webCtx.JSON(http.StatusOK, map[string]string{
			"a": "1",
		})
	}, nil))

	go func() {
		address := ":8080"
		err := e.Start(address)
		fmt.Println("exit http server on", address, "return", err)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, _ := http.Get("http://127.0.0.1:8080/")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	e.Shutdown(context.Background())
	time.Sleep(100 * time.Millisecond)
}

func TestGinServer(t *testing.T) {
	g := gin.New()

	// 配合 gin 框架使用
	g.GET("/", SpringGin.HandlerWrapper("/", func(webCtx SpringWeb.WebContext) {
		webCtx.JSON(http.StatusOK, map[string]string{
			"a": "1",
		})
	}, nil))

	go func() {
		address := ":8080"
		err := g.Run(address)
		fmt.Println("exit http server on", address, "return", err)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, _ := http.Get("http://127.0.0.1:8080/")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	//g.Shutdown(context.Background())
	time.Sleep(100 * time.Millisecond)
}

func TestHttpServer(t *testing.T) {
	var server *http.Server

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		// 暂未实现
	})

	go func() {
		address := ":8080"
		server = &http.Server{Addr: address}
		err := server.ListenAndServe()
		fmt.Println("exit http server on", address, "return", err)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, _ := http.Get("http://127.0.0.1:8080/")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	server.Shutdown(context.Background())
	time.Sleep(100 * time.Millisecond)
}
