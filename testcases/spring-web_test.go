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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/go-spring/go-spring-parent/spring-utils"
	"github.com/go-spring/go-spring-web/spring-echo"
	"github.com/go-spring/go-spring-web/spring-gin"
	"github.com/go-spring/go-spring-web/spring-swagger"
	"github.com/go-spring/go-spring-web/spring-web"
	"github.com/stretchr/testify/assert"
)

type NumberFilter struct {
	l *list.List
	n int
}

func NewNumberFilter(n int, l *list.List) *NumberFilter {
	return &NumberFilter{
		l: l,
		n: n,
	}
}

func (f *NumberFilter) Invoke(ctx SpringWeb.WebContext, chain *SpringWeb.FilterChain) {

	defer func() {
		fmt.Println("::after", f.n)
		f.l.PushBack(f.n)
	}()

	fmt.Println("::before", f.n)
	f.l.PushBack(f.n)

	chain.Next(ctx)
}

func TestFilterChain(t *testing.T) {
	l := list.New()

	filters := []SpringWeb.Filter{
		NewNumberFilter(2, l),
		NewNumberFilter(5, l),
	}

	chain := SpringWeb.NewFilterChain(filters)
	chain.Next(nil)

	assert.Equal(t, SpringUtils.NewList(2, 5, 5, 2), l)
	assert.NotEqual(t, SpringUtils.NewList(1, 2, 3), l)
}

type Service struct {
	store map[string]string
}

func NewService() *Service {
	return &Service{
		store: make(map[string]string),
	}
}

func (s *Service) Get(ctx SpringWeb.WebContext) {

	key := ctx.QueryParam("key")
	ctx.LogInfo("/get", "key=", key)

	val := s.store[key]
	ctx.LogInfo("/get", "val=", val)

	ctx.String(http.StatusOK, val)
}

func (s *Service) Set(ctx SpringWeb.WebContext) {

	var param struct {
		A string `form:"a" json:"a"`
	}

	ctx.Bind(&param)

	ctx.LogInfo("/set", "param="+SpringUtils.ToJson(param))

	s.store["a"] = param.A
}

func (s *Service) Panic(ctx SpringWeb.WebContext) {
	panic("this is a panic")
}

func TestContainer(t *testing.T) {

	testRun := func(c SpringWeb.WebContainer) {

		s := NewService()

		l := list.New()
		f2 := NewNumberFilter(2, l)
		f5 := NewNumberFilter(5, l)
		f7 := NewNumberFilter(7, l)

		c.GET("/get", s.Get, f2, f5)

		if false {
			// 流式风格
			c.Route("", f2, f7).
				POST("/set", s.Set).
				GET("/panic", s.Panic)
		}

		// 障眼法
		r := c.Route("", f2, f7)
		{
			r.POST("/set", s.Set)
			r.GET("/panic", s.Panic)
		}

		if false {
			// 回调风格
			c.Group("", func(r *SpringWeb.Route) {
				r.GET("/panic", s.Panic)
				r.POST("/set", s.Set)
			}, f2, f7)
		}

		go c.Start(":8080")

		time.Sleep(time.Millisecond * 100)
		fmt.Println()

		resp, _ := http.Get("http://127.0.0.1:8080/get?key=a")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()

		http.PostForm("http://127.0.0.1:8080/set", url.Values{
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
	}

	t.Run("SpringGin", func(t *testing.T) {
		testRun(SpringGin.NewContainer())
	})

	t.Run("SpringEcho", func(t *testing.T) {
		testRun(SpringEcho.NewContainer())
	})
}

func TestSwagger(t *testing.T) {

	testRun := func(c SpringWeb.WebContainer) {

		l := list.New()
		f2 := NewNumberFilter(2, l)
		f5 := NewNumberFilter(5, l)
		f7 := NewNumberFilter(7, l)

		get := func(ctx SpringWeb.WebContext) {
			fmt.Println("invoke get()")
			ctx.String(http.StatusOK, "1")
		}

		c.GET(SpringSwagger.GET("/get", get, f2, f5, f7).Doc("get doc").Build())

		go c.Start(":8080")

		time.Sleep(time.Millisecond * 100)
		fmt.Println()

		resp, _ := http.Get("http://127.0.0.1:8080/get?key=a")
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("code:", resp.StatusCode, "||", "resp:", string(body))
		fmt.Println()
	}

	t.Run("SpringGin", func(t *testing.T) {
		testRun(SpringGin.NewContainer())
	})

	t.Run("SpringEcho", func(t *testing.T) {
		testRun(SpringEcho.NewContainer())
	})
}
