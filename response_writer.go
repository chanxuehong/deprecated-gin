// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
	WroteHeader() bool // WroteHeader returns true if header has been written, otherwise false.
	Status() int       // Status returns the response status code of the current request.
	Written() int64    // Written returns number of bytes written in body.
}

var _ ResponseWriter = (*responseWriter)(nil)

type responseWriter struct {
	responseWriter http.ResponseWriter
	wroteHeader    bool  // reply header has been written
	status         int   // status code passed to WriteHeader
	written        int64 // number of bytes written in body
}

func (w *responseWriter) reset(writer http.ResponseWriter) {
	w.responseWriter = writer
	w.wroteHeader = false
	w.status = http.StatusOK
	w.written = 0
}

func (w *responseWriter) WroteHeader() bool {
	return w.wroteHeader
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Written() int64 {
	return w.written
}

func (w *responseWriter) Header() http.Header {
	return w.responseWriter.Header()
}

func (w *responseWriter) WriteHeader(code int) {
	if w.wroteHeader {
		debugPrintf("[WARNING] multiple ResponseWriter.WriteHeader calls\r\n")
		return
	}
	w.wroteHeader = true
	w.status = code
	w.responseWriter.WriteHeader(code)
}

func (w *responseWriter) Write(data []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.Write(data)
	w.written += int64(n)
	return
}
