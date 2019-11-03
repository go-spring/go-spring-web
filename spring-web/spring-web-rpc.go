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
	"fmt"
	"net/http"

	"github.com/go-spring/go-spring-parent/spring-error"
)

//
// 定义 Web RPC 处理函数
//
type RpcHandler func(WebContext) interface{}

//
// 定义 Web RPC 适配函数
//
func RPC(fn RpcHandler) Handler {
	return func(webCtx WebContext) {

		// HTTP RPC 只能返回 json 格式的数据
		webCtx.Header("Content-Type", "application/json")

		defer func() {
			if r := recover(); r != nil {
				result, ok := r.(*SpringError.RpcResult)
				if !ok {
					e := fmt.Sprint(r)
					result = SpringError.ERROR.Error(e)
				}
				webCtx.JSON(http.StatusOK, result)
			}
		}()

		result := SpringError.SUCCESS.Data(fn(webCtx))
		webCtx.JSON(http.StatusOK, result)
	}
}
