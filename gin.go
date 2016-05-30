// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
	"sync"
)

const __maxHandlerChainSize = 64

type (
	HandlerChain []HandlerFunc
	HandlerFunc  func(*Context)
)

func (chain HandlerChain) last() HandlerFunc {
	if n := len(chain); n > 0 {
		return chain[n-1]
	}
	return nil
}

var __httpMethods = [...]string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

var _ http.Handler = (*Engine)(nil)

// Engine is the framework's instance, it contains the muxer, middleware and configuration settings.
// Create an instance of Engine, by using New() or Default().
//
//  NOTE:
//  Engine is NOT copyable.
//  Please serially sets Engine(normally in 'init' or 'main' goroutine).
type Engine struct {
	startedChecker startedChecker

	RouteGroup

	contextPool      sync.Pool
	defaultValidator StructValidator // Validate object when binding if set.

	noRoute     HandlerChain
	noMethod    HandlerChain
	allNoRoute  HandlerChain // always == combineHandlerChain(middlewares, noRoute) if noRoute is not empty, otherwise is nil.
	allNoMethod HandlerChain // always == combineHandlerChain(middlewares, noMethod) if noMethod is not empty, otherwise is nil.

	trees       trees // point treesBuffer
	treesBuffer [len(__httpMethods)]tree

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 307 for all other request methods.
	redirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 307 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// redirectTrailingSlash is independent of this option.
	redirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	handleMethodNotAllowed bool

	// If enabled, the engine get client IP from http hearder X-Real-IP, X-Forwarded-For.
	fetchClientIPFromHeader bool
}

// New returns a new blank Engine instance without any middleware attached.
// By default the configuration is:
// - RedirectTrailingSlash:  true
// - RedirectFixedPath:      false
// - HandleMethodNotAllowed: false
func New() *Engine {
	debugPrintEngineNew()
	engine := &Engine{
		redirectTrailingSlash:   true,
		redirectFixedPath:       false,
		handleMethodNotAllowed:  false,
		fetchClientIPFromHeader: false,
	}
	engine.RouteGroup.basePath = "/"
	engine.RouteGroup.engine = engine
	engine.contextPool.New = contextPoolNew
	engine.trees = engine.treesBuffer[:0]
	return engine
}

func contextPoolNew() interface{} { return new(Context) }

func (engine *Engine) addRoute(method, path string, handlers HandlerChain) {
	if method == "" {
		panic("http method can not be empty")
	}
	if path == "" || path[0] != '/' {
		panic("path must begin with '/'")
	}
	if len(handlers) == 0 {
		panic("handlers can not be empty")
	}
	for _, h := range handlers {
		if h == nil {
			panic("each handler of handlers can not be nil")
		}
	}

	engine.startedChecker.check() // check if engine has been started.
	debugPrintRoute(method, path, handlers)

	root := engine.trees.getTree(method)
	if root == nil {
		root = new(node)
		engine.trees.addTree(method, root)
	}
	root.addRoute(path, handlers)
}

// Routes returns a slice of registered routes, including some useful information, such as:
// the http method, path and the handler name.
func (engine *Engine) Routes() []Route {
	return engine.trees.routes()
}

// Attachs a global middleware to Engine. ie. the middleware attached though Use() will be
// included in the handler chain for every single request. Even 404, 405, static files...
// For example, this is the right place for a logger or error management middleware.
func (engine *Engine) Use(middleware ...HandlerFunc) {
	for _, h := range middleware {
		if h == nil {
			panic("each middleware can not be nil")
		}
	}
	engine.startedChecker.check() // check if engine has been started.
	engine.RouteGroup.Use(middleware...)
	engine.rebuild404Handlers()
	engine.rebuild405Handlers()
}

// NoRoute set handlers for NoRoute. It return a 404 code by default.
// Engine.NoRoute() removes all no-route handlers.
func (engine *Engine) NoRoute(handlers ...HandlerFunc) {
	for _, h := range handlers {
		if h == nil {
			panic("each handler of handlers can not be nil")
		}
	}
	engine.startedChecker.check() // check if engine has been started.
	engine.noRoute = handlers
	engine.rebuild404Handlers()
}

func (engine *Engine) rebuild404Handlers() {
	if len(engine.noRoute) == 0 {
		engine.allNoRoute = nil
		return
	}
	engine.allNoRoute = combineHandlerChain(engine.middlewares, engine.noRoute)
}

// NoMethod set handlers for NoMethod. It return a 405 code by default.
// Engine.NoMethod() removes all no-method handlers.
func (engine *Engine) NoMethod(handlers ...HandlerFunc) {
	for _, h := range handlers {
		if h == nil {
			panic("each handler of handlers can not be nil")
		}
	}
	engine.startedChecker.check() // check if engine has been started.
	engine.noMethod = handlers
	engine.rebuild405Handlers()
}

func (engine *Engine) rebuild405Handlers() {
	if len(engine.noMethod) == 0 {
		engine.allNoMethod = nil
		return
	}
	engine.allNoMethod = combineHandlerChain(engine.middlewares, engine.noMethod)
}

// Enables automatic redirection if the current route can't be matched but a
// handler for the path with (without) the trailing slash exists.
// For example if /foo/ is requested but a route only exists for /foo, the
// client is redirected to /foo with http status code 301 for GET requests
// and 307 for all other request methods.
//
// Default is true.
func (engine *Engine) RedirectTrailingSlash(b bool) {
	engine.startedChecker.check() // check if engine has been started.
	engine.redirectTrailingSlash = b
}

// If enabled, the router tries to fix the current request path, if no
// handle is registered for it.
// First superfluous path elements like ../ or // are removed.
// Afterwards the router does a case-insensitive lookup of the cleaned path.
// If a handle can be found for this route, the router makes a redirection
// to the corrected path with status code 301 for GET requests and 307 for
// all other request methods.
// For example /FOO and /..//Foo could be redirected to /foo.
// RedirectTrailingSlash is independent of this option.
//
// Default is false.
func (engine *Engine) RedirectFixedPath(b bool) {
	engine.startedChecker.check() // check if engine has been started.
	engine.redirectFixedPath = b
}

// If enabled, the router checks if another method is allowed for the
// current route, if the current request can not be routed.
// If this is the case, the request is answered with 'Method Not Allowed'
// and HTTP status code 405.
// If no other Method is allowed, the request is delegated to the NotFound
// handler.
//
// Default is false.
func (engine *Engine) HandleMethodNotAllowed(b bool) {
	engine.startedChecker.check() // check if engine has been started.
	engine.handleMethodNotAllowed = b
}

// If enabled, the engine get client IP from http hearder X-Real-IP, X-Forwarded-For.
//
// Default is false.
func (engine *Engine) FetchClientIPFromHeader(b bool) {
	engine.startedChecker.check() // check if engine has been started.
	engine.fetchClientIPFromHeader = b
}

// DefaultValidator sets the default Validator to validate object when binding.
func (engine *Engine) DefaultValidator(v StructValidator) {
	engine.startedChecker.check() // check if engine has been started.
	engine.defaultValidator = v
}

// ================================================================================================================

// Run attaches the engine to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, engine)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) Run(addr string) (err error) {
	engine.startedChecker.start()
	defer func() {
		debugPrintError(err)
	}()

	debugPrintf("Listening and serving HTTP on %s\r\n", addr)
	return http.ListenAndServe(addr, engine)
}

// RunTLS attaches the engine to a http.Server and starts listening and serving HTTPS (secure) requests.
// It is a shortcut for http.ListenAndServeTLS(addr, certFile, keyFile, engine)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (engine *Engine) RunTLS(addr string, certFile string, keyFile string) (err error) {
	engine.startedChecker.start()
	defer func() {
		debugPrintError(err)
	}()

	debugPrintf("Listening and serving HTTPS on %s\r\n", addr)
	return http.ListenAndServeTLS(addr, certFile, keyFile, engine)
}

// ServeHTTP implements the http.Handler interface.
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	engine.startedChecker.start()
	ctx := engine.contextPool.Get().(*Context)

	ctx.reset()
	ctx.ResponseWriter = ctx.responseWriterArray.ResponseWriter(w)
	ctx.Request = r
	ctx.PathParams = ctx.pathParamsBuffer[:0]
	ctx.Validator = engine.defaultValidator
	ctx.fetchClientIPFromHeader = engine.fetchClientIPFromHeader
	engine.serveHTTP(ctx)

	engine.contextPool.Put(ctx)
}

func (engine *Engine) serveHTTP(ctx *Context) {
	httpMethod := ctx.Request.Method
	path := ctx.Request.URL.Path

	// find root of the tree for the given HTTP method
	root := engine.trees.getTree(httpMethod)
	if root != nil {
		// find route in tree
		handlers, params, tsr := root.getValue(path, ctx.PathParams[:0])
		if handlers != nil {
			ctx.handlers = handlers
			ctx.PathParams = params
			ctx.Next()
			return
		}
		if httpMethod != http.MethodConnect && path != "/" {
			if tsr && engine.redirectTrailingSlash {
				redirectTrailingSlash(ctx)
				return
			}
			if engine.redirectFixedPath && redirectFixedPath(ctx, root, engine.redirectTrailingSlash) {
				return
			}
		}
	}

	// Handle 405
	if engine.handleMethodNotAllowed {
		trees := engine.trees
		for i := 0; i < len(trees); i++ {
			if trees[i].method == httpMethod {
				continue // Skip the requested method - we already tried this one
			}
			if handlers, params, _ := trees[i].root.getValue(path, ctx.PathParams[:0]); handlers != nil {
				ctx.handlers = engine.allNoMethod
				ctx.PathParams = params
				serveError(ctx, 405, __default405Body)
				return
			}
		}
	}

	// Handle 404
	ctx.handlers = engine.allNoRoute
	serveError(ctx, 404, __default404Body)
}

var (
	__default404Body = []byte("404 page not found")
	__default405Body = []byte("405 method not allowed")
)

func serveError(ctx *Context, defaultCode int, defaultMessage []byte) {
	ctx.Next()
	if w := ctx.ResponseWriter; !w.WroteHeader() {
		w.Header().Set(HeaderContentType, MIMETextPlainCharsetUTF8)
		w.Header().Set(HeaderXContentTypeOptions, "nosniff")
		w.WriteHeader(defaultCode)
		w.Write(defaultMessage)
	}
}

func redirectTrailingSlash(ctx *Context) {
	req := ctx.Request
	path := req.URL.Path // path != "/"

	code := 301
	if req.Method != http.MethodGet {
		code = 307
	}

	if len(path) > 1 && path[len(path)-1] == '/' {
		req.URL.Path = path[:len(path)-1]
	} else {
		req.URL.Path = path + "/"
	}

	debugPrintf("redirecting request %d: %s --> %s\r\n", code, path, req.URL.Path)
	http.Redirect(ctx.ResponseWriter, req, req.URL.String(), code)
}

func redirectFixedPath(ctx *Context, root *node, fixTrailingSlash bool) bool {
	req := ctx.Request
	path := req.URL.Path // path != "/"

	fixedPath, found := root.findCaseInsensitivePath(pathClean(path), fixTrailingSlash)
	if found {
		code := 301
		if req.Method != http.MethodGet {
			code = 307
		}
		req.URL.Path = string(fixedPath)
		debugPrintf("redirecting request %d: %s --> %s\r\n", code, path, req.URL.Path)
		http.Redirect(ctx.ResponseWriter, req, req.URL.String(), code)
		return true
	}
	return false
}
