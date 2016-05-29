// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
)

// Common HTTP methods.
//
// Unless otherwise noted, these are defined in RFC 7231 section 4.3.
const (
	HTTPMethodGet     = http.MethodGet
	HTTPMethodHead    = http.MethodHead
	HTTPMethodPost    = http.MethodPost
	HTTPMethodPut     = http.MethodPut
	HTTPMethodPatch   = http.MethodPatch
	HTTPMethodDelete  = http.MethodDelete
	HTTPMethodConnect = http.MethodConnect
	HTTPMethodOptions = http.MethodOptions
	HTTPMethodTrace   = http.MethodTrace
)

// MIME types
const (
	MIMEApplicationURLEncodedForm = "application/x-www-form-urlencoded"
	MIMEMultipartForm             = "multipart/form-data"
	MIMEApplicationOctetStream    = "application/octet-stream"
	MIMEApplicationProtobuf       = "application/protobuf"

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

// Headers
const (
	HeaderAcceptEncoding                = "Accept-Encoding"
	HeaderAuthorization                 = "Authorization"
	HeaderContentDisposition            = "Content-Disposition"
	HeaderContentEncoding               = "Content-Encoding"
	HeaderContentLength                 = "Content-Length"
	HeaderContentType                   = "Content-Type"
	HeaderCookie                        = "Cookie"
	HeaderSetCookie                     = "Set-Cookie"
	HeaderIfModifiedSince               = "If-Modified-Since"
	HeaderLastModified                  = "Last-Modified"
	HeaderLocation                      = "Location"
	HeaderUpgrade                       = "Upgrade"
	HeaderVary                          = "Vary"
	HeaderWWWAuthenticate               = "WWW-Authenticate"
	HeaderXForwardedProto               = "X-Forwarded-Proto"
	HeaderXHTTPMethodOverride           = "X-HTTP-Method-Override"
	HeaderXForwardedFor                 = "X-Forwarded-For"
	HeaderXRealIP                       = "X-Real-IP"
	HeaderServer                        = "Server"
	HeaderOrigin                        = "Origin"
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderXCSRFToken              = "X-CSRF-Token"
)
