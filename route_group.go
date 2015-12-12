// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
	"path"
	"regexp"
	"strings"
)

// RouteGroup is used to configure the internal router,
// a RouteGroup is associated with a path prefix and an array of middlewares.
//
//  NOTE:
//  Please serially sets RouteGroup(normally in 'init' or 'main' goroutine).
type RouteGroup struct {
	basePath    string
	middlewares HandlerChain
	engine      *Engine
}

// Group creates a new RouteGroup.
// You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func (group *RouteGroup) Group(relativePath string, middlewares ...HandlerFunc) *RouteGroup {
	if strings.ContainsRune(relativePath, ':') || strings.ContainsRune(relativePath, '*') {
		panic("path parameters can not be used when grouping")
	}
	return &RouteGroup{
		basePath:    pathJoin(group.basePath, relativePath),
		middlewares: combineHandlerChain(group.middlewares, middlewares),
		engine:      group.engine,
	}
}

// BasePath returns the path prefix of this RouteGroup.
func (group *RouteGroup) BasePath() string { return group.basePath }

// Use adds middlewares to the group.
func (group *RouteGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = combineHandlerChain(group.middlewares, middlewares)
}

var httpMethodRegexp = regexp.MustCompile("^[A-Z]+$") // support custom methods

// Handle registers a new request handler and middlewares with the given path and method.
// The last handler should be the real handler, the other ones should be middleware
// that can and should be shared among different routes.
// See the example code in github.
//
// For GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS, CONNECT and TRACE requests
// the respective shortcut functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (group *RouteGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) {
	if !httpMethodRegexp.MatchString(httpMethod) {
		panic(`http method "` + httpMethod + `" is not valid`)
	}
	group.handle(httpMethod, relativePath, handlers)
}

func (group *RouteGroup) handle(httpMethod, relativePath string, handlers HandlerChain) {
	if len(handlers) == 0 {
		panic("there must be at least one handler")
	}
	absolutePath := pathJoin(group.basePath, relativePath)
	handlers = combineHandlerChain(group.middlewares, handlers)
	group.engine.addRoute(httpMethod, absolutePath, handlers)
}

// Any registers a route that matches all the HTTP methods:
// GET, POST, PUT, PATCH, HEAD, OPTIONS, DELETE, CONNECT, TRACE.
func (group *RouteGroup) Any(relativePath string, handlers ...HandlerFunc) {
	for _, method := range httpMethods {
		group.handle(method, relativePath, handlers)
	}
}

// Get is a shortcut for group.Handle("GET", relativePath, handlers)
func (group *RouteGroup) Get(relativePath string, handlers ...HandlerFunc) {
	group.handle(HttpMethodGet, relativePath, handlers)
}

// Head is a shortcut for group.Handle("HEAD", relativePath, handlers)
func (group *RouteGroup) Head(relativePath string, handlers ...HandlerFunc) {
	group.handle(HttpMethodHead, relativePath, handlers)
}

// Post is a shortcut for group.Handle("POST", relativePath, handlers)
func (group *RouteGroup) Post(relativePath string, handlers ...HandlerFunc) {
	group.handle(HttpMethodPost, relativePath, handlers)
}

// Put is a shortcut for group.Handle("PUT", relativePath, handlers)
func (group *RouteGroup) Put(relativePath string, handlers ...HandlerFunc) {
	group.handle(HttpMethodPut, relativePath, handlers)
}

// Patch is a shortcut for group.Handle("PATCH", relativePath, handlers)
func (group *RouteGroup) Patch(relativePath string, handlers ...HandlerFunc) {
	group.handle(HttpMethodPatch, relativePath, handlers)
}

// Delete is a shortcut for group.Handle("DELETE", relativePath, handlers)
func (group *RouteGroup) Delete(relativePath string, handlers ...HandlerFunc) {
	group.handle(HttpMethodDelete, relativePath, handlers)
}

// Connect is a shortcut for group.Handle("CONNECT", relativePath, handlers)
func (group *RouteGroup) Connect(relativePath string, handlers ...HandlerFunc) {
	group.handle(HttpMethodConnect, relativePath, handlers)
}

// Options is a shortcut for group.Handle("OPTIONS", relativePath, handlers)
func (group *RouteGroup) Options(relativePath string, handlers ...HandlerFunc) {
	group.handle(HttpMethodOptions, relativePath, handlers)
}

// Trace is a shortcut for group.Handle("TRACE", relativePath, handlers)
func (group *RouteGroup) Trace(relativePath string, handlers ...HandlerFunc) {
	group.handle(HttpMethodTrace, relativePath, handlers)
}

// StaticFile registers a single route in order to serve a single file of the local filesystem.
//   group.StaticFile("favicon.ico", "./resources/favicon.ico")
func (group *RouteGroup) StaticFile(relativePath, filepath string) {
	handler := func(ctx *Context) {
		ctx.File(filepath)
	}
	group.Get(relativePath, handler)
	group.Head(relativePath, handler)
}

// StaticRoot like nginx root location but relativePath must be ended with '/' if not empty.
//   group.StaticRoot("/abc/", http.Dir("/home/root"))
//   group.StaticRoot("/abc/", gin.Dir("/home/root"))
//
//   is like this in nginx:
//
//   location ^~ /abc/ {
//       root   /home/root;
//       index  index.html;
//   }
func (group *RouteGroup) StaticRoot(relativePath string, root http.FileSystem) {
	if length := len(relativePath); length > 0 && relativePath[length-1] != '/' {
		panic("path must be ended with '/' when serving a static root")
	}
	if strings.ContainsRune(relativePath, '*') {
		panic("path parameter '*' can not be used when serving a static root")
	}
	handler := func(ctx *Context) {
		http.FileServer(root).ServeHTTP(ctx.ResponseWriter, ctx.Request)
	}
	relativePath = path.Join(relativePath, "*filepath")
	group.Get(relativePath, handler)
	group.Head(relativePath, handler)
}

// StaticAlias like nginx alias location(also relativePath must be ended with '/' if not empty).
//   group.StaticAlias("/abc/", http.Dir("/home/root/abc/"))
//   group.StaticAlias("/abc/", gin.Dir("/home/root/abc/"))
//
//   is like this in nginx:
//
//   location ^~ /abc/ {
//       alias  /home/root/abc/;
//       index  index.html;
//   }
func (group *RouteGroup) StaticAlias(relativePath string, dir http.FileSystem) {
	if length := len(relativePath); length > 0 {
		if relativePath[length-1] != '/' {
			panic("path must be ended with '/' when serving a static alias")
		}
	} else {
		relativePath = "/"
	}
	if strings.ContainsRune(relativePath, ':') || strings.ContainsRune(relativePath, '*') {
		panic("path parameters can not be used when serving a static alias")
	}
	prefix := pathJoin(group.basePath, relativePath)
	handler := func(ctx *Context) {
		http.StripPrefix(prefix, http.FileServer(dir)).ServeHTTP(ctx.ResponseWriter, ctx.Request)
	}
	relativePath = path.Join(relativePath, "*filepath")
	group.Get(relativePath, handler)
	group.Head(relativePath, handler)
}

func pathJoin(basePath, relativePath string) string {
	if len(relativePath) == 0 {
		return basePath // basePath is canonical and beginning with '/'.
	}
	return pathClean(basePath + "/" + relativePath)
}

func combineHandlerChain(middlewares, handlers HandlerChain) HandlerChain {
	if len(handlers) == 0 {
		return middlewares
	}
	// middlewares is canonical, since it was returned by combineHandlerChain.
	for _, h := range handlers {
		if h == nil {
			panic("handler can not be nil")
		}
	}
	size := len(middlewares) + len(handlers)
	if size > maxHandlerChainSize {
		panic("too many handlers")
	}
	combinedHandlerChain := make(HandlerChain, size)
	copy(combinedHandlerChain, middlewares)
	copy(combinedHandlerChain[len(middlewares):], handlers)
	return combinedHandlerChain
}
