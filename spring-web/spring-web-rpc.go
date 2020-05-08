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

package SpringWeb

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-spring/go-spring-parent/spring-error"
	"github.com/go-spring/go-spring-parent/spring-utils"
)

// rpcHandler RPC 形式的 Web 处理接口
type rpcHandler func(WebContext) interface{}

func (r rpcHandler) Invoke(ctx WebContext) {
	rpcInvoke(ctx, func() interface{} { return r(ctx) })
}

func (r rpcHandler) FileLine() (file string, line int, fnName string) {
	return SpringUtils.FileLine(r)
}

// RPC 转换成 RPC 形式的 Web 处理接口
func RPC(fn func(WebContext) interface{}) Handler {
	return rpcHandler(fn)
}

// bindHandler BIND 形式的 Web 处理接口
type bindHandler struct {
	fn       interface{}   // 原始函数的指针
	fnVal    reflect.Value // 原始函数的值
	bindType reflect.Type  // 待绑定的类型
	ctxIndex int           // ctx 变量的位置
}

func (b *bindHandler) Invoke(ctx WebContext) {
	rpcInvoke(ctx, func() interface{} {

		var (
			err     error
			bindVal reflect.Value
		)

		// 获取待绑定的值
		if b.bindType.Kind() == reflect.Ptr {
			bindVal = reflect.New(b.bindType.Elem())
			err = ctx.Bind(bindVal.Interface())
		} else {
			bindVal = reflect.New(b.bindType)
			err = ctx.Bind(bindVal.Interface())
			bindVal = bindVal.Elem()
		}

		SpringError.ERROR.Panic(err).When(err != nil)

		var in []reflect.Value

		// 组装请求参数
		if b.ctxIndex == 0 {
			// func(WebContext,Request)Response
			in = append(in, reflect.ValueOf(ctx))
			in = append(in, bindVal)
		} else if b.ctxIndex == 1 {
			// func(Request,WebContext)Response
			in = append(in, bindVal)
			in = append(in, reflect.ValueOf(ctx))
		} else {
			// func(Request)Response
			in = append(in, bindVal)
		}

		// 执行处理函数，并返回结果
		outVal := b.fnVal.Call(in)
		return outVal[0].Interface()
	})
}

func (b *bindHandler) FileLine() (file string, line int, fnName string) {
	return SpringUtils.FileLine(b.fn)
}

func validBindFn(fn interface{}) (reflect.Type, int, bool) {
	fnTyp := reflect.TypeOf(fn)

	// 必须是函数
	if fnTyp.Kind() != reflect.Func {
		return nil, -1, false
	}

	// 只能有一个返回值
	if fnTyp.NumOut() != 1 {
		return nil, -1, false
	}

	// 待绑定参数必须是结构体或者结构体的指针
	validBindType := func(t reflect.Type) bool {
		return SpringUtils.Indirect(t).Kind() == reflect.Struct
	}

	// 只有一个入参
	if fnTyp.NumIn() == 1 {
		// func(Request)Response
		bindType := fnTyp.In(0)
		if !validBindType(bindType) {
			return nil, -1, false
		}
		return bindType, -1, true
	}

	// 有两个入参
	if fnTyp.NumIn() == 2 {
		t0 := fnTyp.In(0)
		if t0 == webContextType {
			// func(WebContext,Request)Response
			bindType := fnTyp.In(1)
			if !validBindType(bindType) {
				return nil, -1, false
			}
			return bindType, 0, true
		} else {
			// func(Request,WebContext)Response
			bindType := t0
			if !validBindType(bindType) {
				return nil, -1, false
			}
			if fnTyp.In(1) != webContextType {
				return nil, -1, false
			}
			return bindType, 1, true
		}
	}

	return nil, -1, false
}

// BIND 转换成 BIND 形式的 Web 处理接口
func BIND(fn interface{}) Handler {

	var (
		ok       bool
		ctxIndex int
		bindType reflect.Type
	)

	if bindType, ctxIndex, ok = validBindFn(fn); !ok {
		panic(errors.New("fn should be func(req:struct)resp:anything or " +
			"func(ctx:WebContext,req:struct)resp:anything or " +
			"func(req:struct,ctx:WebContext)resp:anything"))
	}

	return &bindHandler{
		fn:       fn,
		fnVal:    reflect.ValueOf(fn),
		bindType: bindType,
		ctxIndex: ctxIndex,
	}
}

func rpcInvoke(webCtx WebContext, fn func() interface{}) {

	// 目前 HTTP RPC 只能返回 json 格式的数据
	webCtx.Header("Content-Type", "application/json")

	defer func() {
		if r := recover(); r != nil {
			result, ok := r.(*SpringError.RpcResult)
			if !ok {
				var err error
				if err, ok = r.(error); !ok {
					err = errors.New(fmt.Sprint(r))
				}
				result = SpringError.ERROR.Error(err)
			}
			webCtx.JSON(http.StatusOK, result)
		}
	}()

	result := SpringError.SUCCESS.Data(fn())
	webCtx.JSON(http.StatusOK, result)
}
