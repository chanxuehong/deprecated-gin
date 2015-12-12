// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	_filepath "path/filepath"

	"github.com/chanxuehong/gin/binder"
)

const (
	// MIME
	MIMEMultipartForm             = "multipart/form-data"
	MIMEApplicationURLEncodedForm = "application/x-www-form-urlencoded"
	MIMEApplicationOctetStream    = "application/octet-stream"

	MIMEApplicationJSON       = "application/json"
	MIMEApplicationXML        = "application/xml"
	MIMEApplicationJavaScript = "application/javascript"
	MIMETextPlain             = "text/plain"
	MIMETextHTML              = "text/html"
	MIMETextXML               = "text/xml"

	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; charset=utf-8"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; charset=utf-8"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; charset=utf-8"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; charset=utf-8"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; charset=utf-8"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; charset=utf-8"
)

const (
	// Headers
	HeaderAcceptEncoding     = "Accept-Encoding"
	HeaderAuthorization      = "Authorization"
	HeaderContentDisposition = "Content-Disposition"
	HeaderContentEncoding    = "Content-Encoding"
	HeaderContentLength      = "Content-Length"
	HeaderContentType        = "Content-Type"
	HeaderLocation           = "Location"
	HeaderWWWAuthenticate    = "WWW-Authenticate"
	HeaderXForwardedFor      = "X-Forwarded-For"
	HeaderXRealIP            = "X-Real-IP"
)

const abortHandlerIndex = maxHandlerChainSize

// Context is the most important part of gin. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
//
//  NOTE:
//  Context and the ResponseWriter, PathParams fields of Context
//  can NOT be safely used outside the request's scope, see Context.Copy().
type Context struct {
	responseWriter ResponseWriter // Context.ResponseWriter points to this
	ResponseWriter *ResponseWriter
	Request        *http.Request

	pathParamsBuffer [256]Param // Context.PathParams points to this
	PathParams       Params
	QueryParams      url.Values

	// Validator will validate object when binding if not nil.
	// By default the Validator is nil, you can set default Validator using
	//     Engine.DefaultValidator()
	// also you can change this in handlers.
	Validator StructValidator

	handlers     HandlerChain
	handlerIndex int

	kvs map[string]interface{}
}

func (ctx *Context) reset(validator StructValidator) {
	ctx.ResponseWriter = &ctx.responseWriter
	ctx.Request = nil
	ctx.PathParams = ctx.pathParamsBuffer[:0]
	ctx.QueryParams = nil
	ctx.Validator = validator
	ctx.handlers = nil
	ctx.handlerIndex = -1
	ctx.kvs = nil
}

// Copy returns a copy of the current context that can be safely used outside the request's scope.
// This have to be used when the context has to be passed to a goroutine.
func (ctx *Context) Copy() *Context {
	var pathParams Params
	if len(ctx.PathParams) > 0 {
		pathParams = append(pathParams, ctx.PathParams...)
	}
	return &Context{
		ResponseWriter: nil,
		Request:        ctx.Request,
		PathParams:     pathParams,
		QueryParams:    ctx.QueryParams,
		Validator:      ctx.Validator,
		handlers:       nil,
		handlerIndex:   abortHandlerIndex,
		kvs:            ctx.kvs,
	}
}

// HandlerName returns the main handle's name.
// For example if the handler is "handleGetUsers()", this function will return "main.handleGetUsers"
func (ctx *Context) HandlerName() string {
	return nameOfFunction(ctx.handlers.last())
}

// IsAborted returns true if the currect context was aborted.
func (ctx *Context) IsAborted() bool {
	return ctx.handlerIndex >= abortHandlerIndex
}

// Abort stops the system to continue calling the pending handlers in the chain.
// Let's say you have an authorization middleware that validates if the request is authorized
// if the authorization fails (the password does not match). This method (Abort()) should be called
// in order to stop the execution of the actual handler.
func (ctx *Context) Abort() {
	ctx.handlerIndex = abortHandlerIndex
}

// AbortWithStatus writes the headers with the specified status code and calls `Abort()`.
// For example, a failed attempt to authentificate a request could use: context.AbortWithStatus(401).
func (ctx *Context) AbortWithStatus(code int) {
	ctx.ResponseWriter.WriteHeader(code)
	ctx.Abort()
}

// AbortWithError writes the status code and error message to response, then stops the chain.
// The error message should be plain text.
// This method calls `String()` internally, see Context.String() for more details.
func (ctx *Context) AbortWithError(code int, format string, values ...interface{}) {
	ctx.String(code, format, values...)
	ctx.Abort()
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
func (ctx *Context) Next() {
	for ctx.handlerIndex++; ctx.handlerIndex < len(ctx.handlers); ctx.handlerIndex++ {
		handler := ctx.handlers[ctx.handlerIndex]
		if handler != nil {
			handler(ctx)
		}
	}
}

// ================================== kvs ======================================

// Set is used to store a new key/value pair exclusivelly for this context.
func (ctx *Context) Set(key string, value interface{}) {
	if ctx.kvs == nil {
		ctx.kvs = make(map[string]interface{})
	}
	ctx.kvs[key] = value
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (ctx *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = ctx.kvs[key]
	return
}

// MustGet Returns the value for the given key if it exists, otherwise it panics.
func (ctx *Context) MustGet(key string) interface{} {
	if value, exists := ctx.Get(key); exists {
		return value
	}
	panic(`[kvs] key "` + key + `" does not exist`)
}

// ================================ request ====================================

// Cookie is a shortcut for ctx.Request.Cookie(name).
// It returns the named cookie provided in the request or http.ErrNoCookie if not found.
func (ctx *Context) Cookie(name string) (*http.Cookie, error) {
	return ctx.Request.Cookie(name)
}

// Param is a shortcut for ctx.PathParams.ByName(name).
func (ctx *Context) Param(name string) string {
	return ctx.PathParams.ByName(name)
}

// ParamByIndex is a shortcut for ctx.PathParams.ByIndex(index).
func (ctx *Context) ParamByIndex(index int) string {
	return ctx.PathParams.ByIndex(index)
}

// Query returns the url query values by name.
func (ctx *Context) Query(name string) string {
	if ctx.QueryParams == nil {
		ctx.QueryParams = ctx.Request.URL.Query()
	}
	return ctx.QueryParams.Get(name)
}

// DefaultQuery like Query if name matched, otherwise it returns defaultValue.
func (ctx *Context) DefaultQuery(name, defaultValue string) string {
	if ctx.QueryParams == nil {
		ctx.QueryParams = ctx.Request.URL.Query()
	}
	if vs := ctx.QueryParams[name]; len(vs) > 0 {
		return vs[0]
	}
	return defaultValue
}

// FormValue is a shortcut for ctx.Request.FormValue(name)
func (ctx *Context) FormValue(name string) string {
	return ctx.Request.FormValue(name)
}

// DefaultFormValue like FormValue if name matched, otherwise it returns defaultValue.
func (ctx *Context) DefaultFormValue(name, defaultValue string) string {
	req := ctx.Request
	if req.Form == nil {
		req.ParseMultipartForm(32 << 20) // 32 MB
	}
	if vs := req.Form[name]; len(vs) > 0 {
		return vs[0]
	}
	return defaultValue
}

// PostFormValue like ctx.Request.PostFormValue(name) but it also gets value from ctx.Request.MultipartForm.Value.
func (ctx *Context) PostFormValue(name string) (value string) {
	value, _ = ctx.postFormValue(name)
	return
}

// DefaultPostFormValue like PostFormValue if name matched, otherwise it returns defaultValue.
func (ctx *Context) DefaultPostFormValue(name, defaultValue string) (value string) {
	var exists bool
	if value, exists = ctx.postFormValue(name); exists {
		return
	}
	return defaultValue
}

func (ctx *Context) postFormValue(name string) (value string, exists bool) {
	req := ctx.Request
	if req.PostForm == nil {
		req.ParseMultipartForm(32 << 20) // 32 MB
	}
	if vs := req.PostForm[name]; len(vs) > 0 {
		return vs[0], true
	}
	if req.MultipartForm != nil && req.MultipartForm.Value != nil {
		if vs := req.MultipartForm.Value[name]; len(vs) > 0 {
			return vs[0], true
		}
	}
	return
}

// BindJSON is a shortcut for ctx.BindWith(obj, binder.JSON).
func (ctx *Context) BindJSON(obj interface{}) (err error) {
	return ctx.BindWith(obj, binder.JSON)
}

// BindXML is a shortcut for ctx.BindWith(obj, binder.XML).
func (ctx *Context) BindXML(obj interface{}) (err error) {
	return ctx.BindWith(obj, binder.XML)
}

// BindWith binds the passed struct pointer using the specified Binder.
func (ctx *Context) BindWith(obj interface{}, b binder.Binder) (err error) {
	if err = b.Bind(ctx.Request, obj); err != nil {
		return
	}
	if validator := ctx.Validator; validator != nil {
		return validator.ValidateStruct(obj)
	}
	return
}

// ================================ response ===================================

// SetCookie is a shortcut for http.SetCookie(ctx.ResponseWriter, cookie).
func (ctx *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.ResponseWriter, cookie)
}

// Redirect redirects the request using http.Redirect with status code.
func (ctx *Context) Redirect(code int, location string) {
	http.Redirect(ctx.ResponseWriter, ctx.Request, location, code)
}

// String writes the given string into the response body.
// It sets the Content-Type as "text/plain; charset=utf-8".
func (ctx *Context) String(code int, format string, values ...interface{}) (err error) {
	w := ctx.ResponseWriter
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	if len(values) > 0 {
		_, err = fmt.Fprintf(w, format, values...)
	} else {
		_, err = w.WriteString(format)
	}
	return
}

// JSON sends a JSON response with status code.
// It sets the Content-Type as "application/json; charset=utf-8".
func (ctx *Context) JSON(code int, obj interface{}) (err error) {
	w := ctx.ResponseWriter
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(obj)
}

// JSONIndent sends a JSON response with status code, but it applies prefix and indent to format the output.
// It sets the Content-Type as "application/json; charset=utf-8".
func (ctx *Context) JSONIndent(code int, obj interface{}, prefix string, indent string) (err error) {
	w := ctx.ResponseWriter
	jsonBytes, err := json.MarshalIndent(obj, prefix, indent)
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_, err = w.Write(jsonBytes)
	return
}

// XML sends an XML response with status code.
// It sets the Content-Type as "application/xml; charset=utf-8".
func (ctx *Context) XML(code int, obj interface{}) (err error) {
	w := ctx.ResponseWriter
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(code)
	if _, err = w.WriteString(xml.Header); err != nil {
		return
	}
	return xml.NewEncoder(w).Encode(obj)
}

// XMLIndent sends an XML response with status code, but it applies prefix and indent to format the output.
// It sets the Content-Type as "application/xml; charset=utf-8".
func (ctx *Context) XMLIndent(code int, obj interface{}, prefix string, indent string) (err error) {
	w := ctx.ResponseWriter
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(code)
	if _, err = w.WriteString(xml.Header); err != nil {
		return
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent(prefix, indent)
	return encoder.Encode(obj)
}

// File writes the specified file into the body stream.
func (ctx *Context) File(filepath string) (err error) {
	http.ServeFile(ctx.ResponseWriter, ctx.Request, filepath)
	return
}

// Attachment adds a `Content-Disposition:attachment;filename="XXX"` header and
// writes the specified file into the body stream.
// filename can be empty, in that case name of the file is used.
func (ctx *Context) Attachment(filepath, filename string) (err error) {
	if filename == "" {
		_, filename = _filepath.Split(filepath)
		if filename == "" {
			return errors.New("invalid filepath")
		}
	}
	ContentDisposition := `attachment;filename="` + url.QueryEscape(filename) + `"`
	w := ctx.ResponseWriter
	w.Header().Set("Content-Disposition", ContentDisposition)
	http.ServeFile(w, ctx.Request, filepath)
	return
}
