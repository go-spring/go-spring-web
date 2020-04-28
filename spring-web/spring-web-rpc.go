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

// rpcHandler Web RPC 处理函数
type rpcHandler func(WebContext) interface{}

func (r rpcHandler) Invoke(ctx WebContext) {
	if r != nil {
		rpcInvoke(ctx, func() interface{} {
			return r(ctx)
		})
	}
}

func (r rpcHandler) FileLine() (file string, line int, fnName string) {
	return SpringUtils.FileLine(r)
}

// RPC Web RPC 处理函数的辅助函数
func RPC(fn func(WebContext) interface{}) Handler {
	return rpcHandler(fn)
}

// bindHandler Web RPC 处理函数
type bindHandler struct {
	fn    interface{}
	inTyp reflect.Type
	fnVal reflect.Value
}

func (b *bindHandler) Invoke(ctx WebContext) {
	rpcInvoke(ctx, func() interface{} {
		inVal := reflect.New(b.inTyp)
		err := ctx.Bind(inVal.Interface())
		SpringError.ERROR.Panic(err).When(err != nil)
		outVal := b.fnVal.Call([]reflect.Value{inVal.Elem()})
		return outVal[0].Interface()
	})
}

func (b *bindHandler) FileLine() (file string, line int, fnName string) {
	return SpringUtils.FileLine(b.fn)
}

// BIND 封装 Bind 操作的 Web RPC 适配函数
func BIND(fn interface{}) Handler {

	fnTyp := reflect.TypeOf(fn)

	// 检查 fn 的类型，必须是 func(req:struct)resp:anything 这样的格式
	if fnTyp.Kind() != reflect.Func || fnTyp.NumIn() != 1 || fnTyp.NumOut() != 1 {
		panic("fn must be func(req:struct)resp:anything")
	}

	return &bindHandler{
		fn:    fn,
		inTyp: fnTyp.In(0),
		fnVal: reflect.ValueOf(fn),
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
