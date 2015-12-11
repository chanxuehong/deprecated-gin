// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

const (
	HttpMethodGet     = "GET"
	HttpMethodHead    = "HEAD"
	HttpMethodPost    = "POST"
	HttpMethodPut     = "PUT"
	HttpMethodPatch   = "PATCH"
	HttpMethodDelete  = "DELETE"
	HttpMethodConnect = "CONNECT"
	HttpMethodOptions = "OPTIONS"
	HttpMethodTrace   = "TRACE"
)

var httpMethods = [...]string{
	HttpMethodGet,
	HttpMethodHead,
	HttpMethodPost,
	HttpMethodPut,
	HttpMethodPatch,
	HttpMethodDelete,
	HttpMethodConnect,
	HttpMethodOptions,
	HttpMethodTrace,
}
