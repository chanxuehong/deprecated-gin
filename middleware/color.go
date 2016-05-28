// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

var (
	__colorGreen   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	__colorWhite   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	__colorYellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	__colorRed     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	__colorBlue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	__colorMagenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	__colorCyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	__colorReset   = string([]byte{27, 91, 48, 109})
)

func colorForStatus(code int) string {
	switch {
	case code < 200:
		return __colorRed
	case code < 300: // [200, 300)
		return __colorGreen
	case code < 400: // [300, 400)
		return __colorWhite
	case code < 500: // [400, 500)
		return __colorYellow
	default: // >= 500
		return __colorRed
	}
}

func colorForMethod(method string) string {
	switch method {
	case "GET":
		return __colorBlue
	case "POST":
		return __colorCyan
	case "PUT":
		return __colorYellow
	case "DELETE":
		return __colorRed
	case "PATCH":
		return __colorGreen
	case "HEAD":
		return __colorMagenta
	case "OPTIONS":
		return __colorWhite
	default:
		return __colorReset
	}
}
