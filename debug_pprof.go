// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"github.com/chanxuehong/gin/pprof"
)

// EnableDebugPProf adds internal HandlerFunc and current middlewares of Engine to process path "/debug/pprof/*name",
// for more information please see net/http/pprof.
func (engine *Engine) EnableDebugPProf() {
	engine.Any("/debug/pprof/*name", debugPProfHandlerFunc)
}

func debugPProfHandlerFunc(ctx *Context) {
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
