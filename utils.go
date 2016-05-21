// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"encoding/xml"
	"net/http"
	"reflect"
	"runtime"
)

type H map[string]interface{}

var _ xml.Marshaler = H(nil)

// MarshalXML implements xml.Marshaler
func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	if err = e.EncodeToken(start); err != nil {
		return
	}
	for key, value := range h {
		if key == "" {
			continue
		}
		elem := xml.StartElement{
			Name: xml.Name{Local: key},
		}
		if err = e.EncodeElement(value, elem); err != nil {
			return
		}
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

func nameOfFunction(fn interface{}) string {
	if fn == nil {
		return "unknown"
	}
	pc := reflect.ValueOf(fn).Pointer()
	if pc == 0 {
		return "unknown"
	}
	return runtime.FuncForPC(pc).Name()
}

func WrapHandlerFunc(fn http.HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		fn(ctx.ResponseWriter, ctx.Request)
	}
}

func WrapHandler(h http.Handler) HandlerFunc {
	return func(ctx *Context) {
		h.ServeHTTP(ctx.ResponseWriter, ctx.Request)
	}
}
