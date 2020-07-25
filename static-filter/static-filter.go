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

package StaticFilter

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-spring/go-spring-parent/spring-logger"
	"github.com/go-spring/go-spring-web/spring-web"
)

// StaticFilter 本机文件列表过滤器
type StaticFilter struct {
	root       string
	urlPrefix  string
	fileServer http.Handler
}

// New StaticFilter 的构造函数
func New(root string, urlPrefix string) *StaticFilter {
	fileServer := http.FileServer(http.Dir(root))
	if urlPrefix != "" {
		fileServer = http.StripPrefix(urlPrefix, fileServer)
	} else {
		SpringLogger.Warn("strongly recommended non-empty urlPrefix")
	}
	return &StaticFilter{
		root:       root,
		urlPrefix:  urlPrefix,
		fileServer: fileServer,
	}
}

func (f *StaticFilter) Exists(filePath string) bool {

	p := strings.TrimPrefix(filePath, f.urlPrefix)
	if f.urlPrefix != "" && len(p) == len(filePath) {
		return false
	}

	name := path.Join(f.root, p)
	_, err := os.Stat(name)
	return err == nil
}

func (f *StaticFilter) Invoke(ctx SpringWeb.WebContext, chain SpringWeb.FilterChain) {
	if f.Exists(ctx.Request().URL.Path) {
		f.fileServer.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
	} else {
		chain.Next(ctx)
	}
}
