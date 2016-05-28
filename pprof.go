// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"

	"github.com/chanxuehong/gin/pprof"
)

// EnablePProf adds "/debug/pprof/*name" to router, for more information please see net/http/pprof.
//
// Please NOTE that EnablePProf without using any middleware inside Engine.
func (engine *Engine) EnablePProf() {
	engine.addRoute(http.MethodGet, "/debug/pprof/*name", HandlerChain{pprofHandlerFunc})
}

func pprofHandlerFunc(ctx *Context) {
	switch path := ctx.Request.URL.Path; path {
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
