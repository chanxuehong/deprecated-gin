// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/chanxuehong/gin/binder"
	"github.com/chanxuehong/gin/internal/response"
)

const (
	__initHandlerIndex  = -1
	__abortHandlerIndex = __maxHandlerChainSize
)

type ResponseWriter interface {
	http.ResponseWriter
	WroteHeader() bool // WroteHeader returns true if header has been written, otherwise false.
	Status() int       // Status returns the response status code of the current request.
	Written() int64    // Written returns number of bytes written in body.
}

// Context is the most important part of gin. It allows us to pass variables between middleware,
// manage the flow, validate the JSON of a request and render a JSON response for example.
//
//  NOTE:
//  Context and the ResponseWriter, PathParams fields of Context
//  can NOT be safely used outside the request's scope, see Context.Copy().
type Context struct {
	responseWriterArray response.ResponseWriterArray // cache pre-allocate ResponseWriter
	ResponseWriter      ResponseWriter
	Request             *http.Request

	pathParamsBuffer [8]Param // Context.PathParams points to this
	PathParams       Params
	queryParams      url.Values

	// Validator will validate object when binding if not nil.
	// By default the Validator is nil, you can set default Validator using
	//     Engine.DefaultValidator()
	// also you can change this in handlers.
	Validator StructValidator

	fetchClientIPFromHeader bool

	handlers     HandlerChain
	handlerIndex int

	kvs map[string]interface{}
}

func (ctx *Context) reset(w http.ResponseWriter, r *http.Request, validator StructValidator, fetchClientIPFromHeader bool) {
	ctx.ResponseWriter = ctx.responseWriterArray.ResponseWriter(w)
	ctx.Request = r
	ctx.PathParams = ctx.pathParamsBuffer[:0]
	ctx.queryParams = nil
	ctx.Validator = validator
	ctx.fetchClientIPFromHeader = fetchClientIPFromHeader
	ctx.handlers = nil
	ctx.handlerIndex = __initHandlerIndex
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
		ResponseWriter:          nil,
		Request:                 ctx.Request,
		PathParams:              pathParams,
		queryParams:             ctx.queryParams,
		Validator:               ctx.Validator,
		fetchClientIPFromHeader: ctx.fetchClientIPFromHeader,
		handlers:                nil,
		handlerIndex:            __abortHandlerIndex,
		kvs:                     ctx.kvs,
	}
}

// HandlerName returns the main handle's name.
// For example if the handler is "handleGetUsers()", this function will return "main.handleGetUsers"
func (ctx *Context) HandlerName() string {
	return nameOfFunction(ctx.handlers.last())
}

// IsAborted returns true if the currect context was aborted.
func (ctx *Context) IsAborted() bool {
	return ctx.handlerIndex >= __abortHandlerIndex
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized. If the
// authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (ctx *Context) Abort() {
	ctx.handlerIndex = __abortHandlerIndex
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
	for {
		ctx.handlerIndex++
		if ctx.handlerIndex >= len(ctx.handlers) {
			ctx.handlerIndex--
			break
		}
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

// Delete remove the value with given key.
func (ctx *Context) Delete(key string) {
	delete(ctx.kvs, key)
}

// Get returns the value for the given key, ie: (value, true).
// If the value does not exists it returns (nil, false)
func (ctx *Context) Get(key string) (value interface{}, exists bool) {
	value, exists = ctx.kvs[key]
	return
}

// MustGet returns the value for the given key if it exists, otherwise it panics.
func (ctx *Context) MustGet(key string) interface{} {
	if value, exists := ctx.Get(key); exists {
		return value
	}
	panic(`[kvs] value with key "` + key + `" does not exist`)
}

// ================================ request ====================================

func (ctx *Context) ClientIP() (ip string) {
	if ctx.fetchClientIPFromHeader {
		header := ctx.Request.Header
		if ip = header.Get("X-Real-Ip"); ip != "" {
			if ip = strings.TrimSpace(ip); ip != "" {
				return ip
			}
		}
		if ip = header.Get("X-Forwarded-For"); ip != "" {
			if index := strings.IndexByte(ip, ','); index >= 0 {
				ip = ip[:index]
			}
			if ip = strings.TrimSpace(ip); ip != "" {
				return ip
			}
		}
	}
	ip, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err != nil {
		return "unknown"
	}
	return ip
}

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

// QueryParams returns the url query values.
func (ctx *Context) QueryParams() url.Values {
	if ctx.queryParams == nil {
		ctx.queryParams = ctx.Request.URL.Query()
	}
	return ctx.queryParams
}

// Query returns the url query value by name.
func (ctx *Context) Query(name string) string {
	return ctx.QueryParams().Get(name)
}

// DefaultQuery like Query if name matched, otherwise it returns defaultValue.
func (ctx *Context) DefaultQuery(name, defaultValue string) string {
	queryParams := ctx.QueryParams()
	if vs := queryParams[name]; len(vs) > 0 {
		return vs[0]
	}
	return defaultValue
}

// FormValue is a shortcut for ctx.Request.FormValue(name).
func (ctx *Context) FormValue(name string) string {
	return ctx.Request.FormValue(name)
}

// DefaultFormValue like FormValue if name matched, otherwise it returns defaultValue.
func (ctx *Context) DefaultFormValue(name, defaultValue string) string {
	r := ctx.Request
	if r.Form == nil {
		r.ParseMultipartForm(32 << 20) // 32 MB
	}
	if vs := r.Form[name]; len(vs) > 0 {
		return vs[0]
	}
	return defaultValue
}

// PostFormValue is a shortcut for ctx.Request.PostFormValue(name).
func (ctx *Context) PostFormValue(name string) (value string) {
	return ctx.Request.PostFormValue(name)
}

// DefaultPostFormValue like PostFormValue if name matched, otherwise it returns defaultValue.
func (ctx *Context) DefaultPostFormValue(name, defaultValue string) (value string) {
	r := ctx.Request
	if r.PostForm == nil {
		r.ParseMultipartForm(32 << 20) // 32 MB
	}
	if vs := r.PostForm[name]; len(vs) > 0 {
		return vs[0]
	}
	return defaultValue
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

// NoContent sends a response with no body and a status code.
func (ctx *Context) NoContent(code int) error {
	ctx.ResponseWriter.WriteHeader(code)
	return nil
}

// String writes the given string into the response body.
// It sets the Content-Type as "text/plain; charset=utf-8".
func (ctx *Context) String(code int, format string, values ...interface{}) (err error) {
	w := ctx.ResponseWriter
	w.Header().Set(HeaderContentType, MIMETextPlainCharsetUTF8)
	w.Header().Set(HeaderXContentTypeOptions, "nosniff")
	w.WriteHeader(code)
	if len(values) > 0 {
		_, err = fmt.Fprintf(w, format, values...)
	} else {
		_, err = io.WriteString(w, format)
	}
	return
}

// JSON sends a JSON response with status code.
// It sets the Content-Type as "application/json; charset=utf-8".
func (ctx *Context) JSON(code int, obj interface{}) (err error) {
	w := ctx.ResponseWriter
	w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(obj)
}

// JSONIndent sends a JSON response with status code, but it applies prefix and indent to format the output.
// It sets the Content-Type as "application/json; charset=utf-8".
func (ctx *Context) JSONIndent(code int, obj interface{}, prefix string, indent string) (err error) {
	w := ctx.ResponseWriter
	body, err := json.MarshalIndent(obj, prefix, indent)
	if err != nil {
		return
	}
	w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(code)
	_, err = w.Write(body)
	return
}

// JSONBlob sends a JSON response with status code.
// It sets the Content-Type as "application/json; charset=utf-8".
func (ctx *Context) JSONBlob(code int, blob []byte) (err error) {
	w := ctx.ResponseWriter
	w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(code)
	_, err = w.Write(blob)
	return
}

// XML sends an XML response with status code.
// It sets the Content-Type as "application/xml; charset=utf-8".
func (ctx *Context) XML(code int, obj interface{}) (err error) {
	w := ctx.ResponseWriter
	w.Header().Set(HeaderContentType, MIMEApplicationXMLCharsetUTF8)
	w.WriteHeader(code)
	if _, err = io.WriteString(w, xml.Header); err != nil {
		return
	}
	return xml.NewEncoder(w).Encode(obj)
}

// XMLIndent sends an XML response with status code, but it applies prefix and indent to format the output.
// It sets the Content-Type as "application/xml; charset=utf-8".
func (ctx *Context) XMLIndent(code int, obj interface{}, prefix string, indent string) (err error) {
	w := ctx.ResponseWriter
	w.Header().Set(HeaderContentType, MIMEApplicationXMLCharsetUTF8)
	w.WriteHeader(code)
	if _, err = io.WriteString(w, xml.Header); err != nil {
		return
	}
	encoder := xml.NewEncoder(w)
	encoder.Indent(prefix, indent)
	return encoder.Encode(obj)
}

var __xmlHeaderPrefix = []byte(`<?xml `) // <?xml version="1.0" encoding="UTF-8"?>

// XMLBlob sends an XML response with status code.
// It sets the Content-Type as "application/xml; charset=utf-8".
func (ctx *Context) XMLBlob(code int, blob []byte) (err error) {
	w := ctx.ResponseWriter
	w.Header().Set(HeaderContentType, MIMEApplicationXMLCharsetUTF8)
	w.WriteHeader(code)
	if !bytes.HasPrefix(bytes.TrimLeftFunc(blob, unicode.IsSpace), __xmlHeaderPrefix) {
		if _, err = io.WriteString(w, xml.Header); err != nil {
			return
		}
	}
	_, err = w.Write(blob)
	return
}

// ServeFile is wrapper for http.ServeFile.
func (ctx *Context) ServeFile(filepath string) (err error) {
	http.ServeFile(ctx.ResponseWriter, ctx.Request, filepath)
	return
}

// ServeContent is wrapper for http.ServeContent.
func (ctx *Context) ServeContent(content io.ReadSeeker, name string, modtime time.Time) (err error) {
	http.ServeContent(ctx.ResponseWriter, ctx.Request, name, modtime, content)
	return
}

// AttachmentFile adds a `Content-Disposition:attachment;filename="XXX"` header
// and writes the specified file into the body stream.
//
// filename can be empty, in that case name of the file is used.
func (ctx *Context) AttachmentFile(filepath, filename string) (err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return
	}
	if fileInfo.IsDir() {
		return errors.New("can't set a dir as file attachment")
	}

	if filename == "" {
		filename = fileInfo.Name()
	}
	return ctx.Attachment(file, filename)
}

// Attachment adds a `Content-Disposition:attachment;filename="XXX"` header
// and writes the specified content into the body stream.
func (ctx *Context) Attachment(content io.Reader, filename string) (err error) {
	w := ctx.ResponseWriter
	w.Header().Set(HeaderContentType, contentTypeFromName(filename))
	w.Header().Set(HeaderContentDisposition, `attachment;filename="`+url.QueryEscape(filename)+`"`)
	_, err = io.Copy(w, content)
	return
}

func contentTypeFromName(name string) string {
	t := mime.TypeByExtension(filepath.Ext(name))
	if t != "" {
		return t
	}
	return MIMEApplicationOctetStream
}
