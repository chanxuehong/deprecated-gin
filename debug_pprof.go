// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"github.com/chanxuehong/gin/pprof"
)

// DebugPProf registers internal HandlerFunc to process path "/debug/pprof/*name", for more information please see net/http/pprof.
//
// Please NOTE that DebugPProf does not use any middleware of Engine, you can specify middlewares optionally when you call this method.
func (engine *Engine) DebugPProf(middleware ...HandlerFunc) {
	for _, h := range middleware {
		if h == nil {
			panic("each middleware can not be nil")
		}
	}
	engine.startedChecker.check() // check if engine has been started.
	handlers := combineHandlerChain(middleware, HandlerChain{debugPProfHandler})
	for _, method := range __httpMethods {
		engine.addRoute(method, "/debug/pprof/*name", handlers)
	}
}

func debugPProfHandler(ctx *Context) {
	path := ctx.Request.URL.Path
	switch path {
	default:
		pprof.Index(ctx.ResponseWriter, ctx.Request)
	case "/debug/pprof/cmdline":
		pprof.Cmdline(ctx.ResponseWriter, ctx.Request)
	case "/debug/pprof/profile":
		pprof.Profile(ctx.ResponseWriter, ctx.Request)
	case "/debug/pprof/symbol":
		pprof.Symbol(ctx.ResponseWriter, ctx.Request)
	case "/debug/pprof/trace":
		pprof.Trace(ctx.ResponseWriter, ctx.Request)
	}
}
