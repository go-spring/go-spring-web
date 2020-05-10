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
	"github.com/go-openapi/spec"
	"github.com/go-spring/go-spring-web/spring-echo"
	"github.com/go-spring/go-spring-web/spring-gin"
	"github.com/go-spring/go-spring-web/spring-web"
	"github.com/go-spring/go-spring-web/testcases"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/magiconair/properties/assert"
)

func TestWebContainer(t *testing.T) {
	cfg := SpringWeb.ContainerConfig{Port: 8080}

	SpringWeb.Swagger().
		WithDescription("web container test").
		AddDefinition("Set", new(spec.Schema).
			Typed("object", "").
			AddRequired("name", "age").
			SetProperty("name", *spec.StringProperty()).
			SetProperty("age", *spec.Int32Property()))

	testRun := func(c SpringWeb.WebContainer) {

		l := list.New()
		f2 := testcases.NewNumberFilter(2, l)
		f5 := testcases.NewNumberFilter(5, l)
		f7 := testcases.NewNumberFilter(7, l)

		c.AddFilter(&testcases.LogFilter{}, &testcases.GlobalInterruptFilter{})

		s := testcases.NewService()

		c.GetMapping("/get", s.Get, f5).Swagger("").
			WithDescription("get").
			AddParam(spec.QueryParam("key")).
			WithConsumes(SpringWeb.MIMEApplicationForm).
			WithProduces(SpringWeb.MIMEApplicationJSON).
			RespondsWith(http.StatusOK, spec.NewResponse().
				WithSchema(spec.StringProperty()).
				AddExample(SpringWeb.MIMEApplicationJSON, 2))

		// 等价于 c.GET("/global_interrupt", s.Get)
		c.HandleGet("/global_interrupt", SpringWeb.METHOD(s, "Get"))

		c.GetMapping("/interrupt", s.Get, f5, &testcases.InterruptFilter{})

		// 障眼法
		r := c.Route("", f2, f7)
		{
			r.PostMapping("/set", s.Set).Swagger("").
				WithDescription("set").
				//WithConsumes(SpringWeb.MIMEApplicationForm).
				WithConsumes(SpringWeb.MIMEApplicationJSON).
				//AddParam(spec.QueryParam("name")).
				//AddParam(spec.QueryParam("age")).
				AddParam(&spec.Parameter{
					ParamProps: spec.ParamProps{
						Name:   "body",
						In:     "body",
						Schema: spec.RefSchema("#/definitions/Set"),
					},
				}).
				RespondsWith(http.StatusOK, nil)

			r.Request(SpringWeb.MethodGetPost, "/panic", s.Panic)
		}

		c.GetMapping("/wild_1/*", func(webCtx SpringWeb.WebContext) {
			assert.Equal(t, "anything", webCtx.PathParam("*"))
			assert.Equal(t, []string{"*"}, webCtx.PathParamNames())
			assert.Equal(t, []string{"anything"}, webCtx.PathParamValues())
			webCtx.JSON(http.StatusOK, webCtx.PathParam("*"))
		})

		c.GetMapping("/wild_2/*none", func(webCtx SpringWeb.WebContext) {
			assert.Equal(t, "anything", webCtx.PathParam("*"))
			assert.Equal(t, "anything", webCtx.PathParam("none"))
			assert.Equal(t, []string{"*"}, webCtx.PathParamNames())
			assert.Equal(t, []string{"anything"}, webCtx.PathParamValues())
			webCtx.JSON(http.StatusOK, webCtx.PathParam("*"))
		})

		c.GetMapping("/wild_3/{*}", func(webCtx SpringWeb.WebContext) {
			assert.Equal(t, "anything", webCtx.PathParam("*"))
			assert.Equal(t, []string{"*"}, webCtx.PathParamNames())
			assert.Equal(t, []string{"anything"}, webCtx.PathParamValues())
			webCtx.JSON(http.StatusOK, webCtx.PathParam("*"))
		})

		// 启动 web 服务器
		c.Start()

		time.Sleep(time.Millisecond * 100)
		fmt.Println()

		resp, _ := http.Get("http://127.0.0.1:8080/get?key=a")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.PostForm("http://127.0.0.1:8080/set", url.Values{
			"a": []string{"1"},
		})
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/get?key=a")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/panic")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.PostForm("http://127.0.0.1:8080/panic", nil)
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/interrupt")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/global_interrupt")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/native")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/swagger/doc.json")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/wild_1/anything")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/wild_2/anything")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/wild_3/anything")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		c.Stop(context.TODO())

		time.Sleep(time.Millisecond * 50)
	}

	t.Run("SpringGin", func(t *testing.T) {
		c := SpringGin.NewContainer(cfg)

		fLogger := SpringGin.Filter(gin.Logger())
		c.SetLoggerFilter(fLogger)

		fRecover := SpringGin.Filter(gin.Recovery())
		c.SetRecoveryFilter(fRecover)

		c.HandleGet("/native", SpringGin.Gin(func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "gin")
		}))

		testRun(c)
	})

	t.Run("SpringEcho", func(t *testing.T) {
		c := SpringEcho.NewContainer(cfg)

		fLogger := SpringEcho.Filter(middleware.Logger())
		c.SetLoggerFilter(fLogger)

		fRecover := SpringEcho.Filter(middleware.Recover())
		c.SetRecoveryFilter(fRecover)

		c.HandleGet("/native", SpringEcho.Echo(func(ctx echo.Context) error {
			return ctx.String(http.StatusOK, "echo")
		}))

		testRun(c)
	})
}

func TestEchoServer(t *testing.T) {
	e := echo.New()
	e.HideBanner = true

	// 配合 echo 框架使用
	e.GET("/*", SpringEcho.HandlerWrapper(
		SpringWeb.FUNC(func(webCtx SpringWeb.WebContext) {
			assert.Equal(t, "echo", webCtx.PathParam("*"))
			assert.Equal(t, []string{"*"}, webCtx.PathParamNames())
			assert.Equal(t, []string{"echo"}, webCtx.PathParamValues())
			webCtx.JSON(http.StatusOK, map[string]string{
				"a": "1",
			})
		}), "", nil))

	go func() {
		address := ":8080"
		err := e.Start(address)
		fmt.Println("exit http server on", address, "return", err)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, _ := http.Get("http://127.0.0.1:8080/echo")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	_ = e.Shutdown(context.Background())
	time.Sleep(100 * time.Millisecond)
}

func TestGinServer(t *testing.T) {
	g := gin.New()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: g,
	}

	// 配合 gin 框架使用
	g.GET("/*"+SpringWeb.DefaultWildCardName, SpringGin.HandlerWrapper("/",
		SpringWeb.FUNC(func(webCtx SpringWeb.WebContext) {
			assert.Equal(t, "gin", webCtx.PathParam("*"))
			assert.Equal(t, []string{"*"}, webCtx.PathParamNames())
			assert.Equal(t, []string{"gin"}, webCtx.PathParamValues())
			webCtx.JSON(http.StatusOK, map[string]string{
				"a": "1",
			})
		}), SpringWeb.DefaultWildCardName, nil)...)

	go func() {
		err := httpServer.ListenAndServe()
		fmt.Println("exit http server on", httpServer.Addr, "return", err)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, _ := http.Get("http://127.0.0.1:8080/gin")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

	_ = httpServer.Shutdown(context.Background())
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

	_ = server.Shutdown(context.Background())
	time.Sleep(100 * time.Millisecond)
}
