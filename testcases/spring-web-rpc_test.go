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
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/go-spring/go-spring-parent/spring-error"
	"github.com/go-spring/go-spring-web/spring-echo"
	"github.com/go-spring/go-spring-web/spring-gin"
	"github.com/go-spring/go-spring-web/spring-web"
)

type RpcService struct{}

func (s *RpcService) OK(ctx SpringWeb.WebContext) interface{} {
	return "123"
}

func (s *RpcService) Err(ctx SpringWeb.WebContext) interface{} {
	panic("err")
}

func (s *RpcService) Panic(ctx SpringWeb.WebContext) interface{} {

	err := errors.New("panic")
	isPanic := ctx.QueryParam("panic") == "1"
	SpringError.ERROR.Panic(err).When(isPanic)

	return "ok"
}

func TestRpc(t *testing.T) {

	testRun := func(c SpringWeb.WebContainer) {

		rc := new(RpcService)

		l := list.New()
		f2 := NewNumberFilter(2, l)
		f5 := NewNumberFilter(5, l)
		f7 := NewNumberFilter(7, l)

		c.GET("/ok", SpringWeb.RPC(rc.OK), f2, f5)

		r := c.Route("", f2, f7)
		{
			r.GET("/err", SpringWeb.RPC(rc.Err))
			r.GET("/panic", SpringWeb.RPC(rc.Panic))
		}

		go c.Start(":8080")

		time.Sleep(time.Millisecond * 100)
		fmt.Println()

		resp, _ := http.Get("http://127.0.0.1:8080/ok")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/err")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		resp, _ = http.Get("http://127.0.0.1:8080/panic")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))

		resp, _ = http.Get("http://127.0.0.1:8080/panic?panic=1")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
	}

	t.Run("SpringGin", func(t *testing.T) {
		testRun(SpringGin.NewContainer())
	})

	t.Run("SpringEcho", func(t *testing.T) {
		testRun(SpringEcho.NewContainer())
	})
}
