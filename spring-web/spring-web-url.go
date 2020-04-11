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

// PathConvert {} 路由风格转换成 : 路由风格
func PathConvert(path string) string {
	var start int
	var brace bool
	var result []byte
	b := []byte(path)
	for k, v := range b {
		if v == '{' {
			start = k
			brace = true
		} else if v == '}' {
			brace = false
			result = append(result, ':')
			result = append(result, b[start+1:k]...)
		} else {
			if !brace {
				result = append(result, v)
			}
		}
	}
	return (string)(result)
}
