// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"bufio"
	"io"
	"net"
	"net/http"
)

var _ http.ResponseWriter = (*ResponseWriter)(nil)

// ResponseWriter implements http.ResponseWriter, http.Flusher, http.Hijacker, http.CloseNotifier and io.ReaderFrom.
type ResponseWriter struct {
	responseWriter http.ResponseWriter
	hijacked       bool  // connection has been hijacked by handler
	wroteHeader    bool  // reply header has been (logically) written
	status         int   // status code passed to WriteHeader
	written        int64 // number of bytes written in body
}

func (w *ResponseWriter) reset(writer http.ResponseWriter) {
	w.responseWriter = writer
	w.hijacked = false
	w.wroteHeader = false
	w.status = http.StatusOK
	w.written = 0
}

// Header returns the header map that will be sent by
// WriteHeader. Changing the header after a call to
// WriteHeader (or Write) has no effect unless the modified
// headers were declared as trailers by setting the
// "Trailer" header before the call to WriteHeader (see example).
// To suppress implicit response headers, set their value to nil.
func (w *ResponseWriter) Header() http.Header {
	return w.responseWriter.Header()
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (w *ResponseWriter) WriteHeader(code int) {
	if w.hijacked {
		debugPrint("[WARNING] ResponseWriter.WriteHeader on hijacked connection\r\n")
		return
	}
	if w.wroteHeader {
		debugPrint("[WARNING] multiple ResponseWriter.WriteHeader calls\r\n")
		return
	}
	w.wroteHeader = true
	w.status = code
	w.responseWriter.WriteHeader(code)
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (w *ResponseWriter) Write(data []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.Write(data)
	w.written += int64(n)
	return
}

// WriteString writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, WriteString calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, WriteString adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (w *ResponseWriter) WriteString(s string) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = io.WriteString(w.responseWriter, s)
	w.written += int64(n)
	return
}

// ReadFrom implements the io.ReaderFrom interface.
func (w *ResponseWriter) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = io.Copy(w.responseWriter, r)
	w.written += n
	return
}

// Flush implements the http.Flusher interface.
func (w *ResponseWriter) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.responseWriter.(http.Flusher).Flush()
}

// Hijack implements the http.Hijacker interface.
func (w *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if !w.hijacked {
		w.hijacked = true
	}
	return w.responseWriter.(http.Hijacker).Hijack()
}

// CloseNotify implements the http.CloseNotifier interface.
func (w *ResponseWriter) CloseNotify() <-chan bool {
	return w.responseWriter.(http.CloseNotifier).CloseNotify()
}

// WroteHeader replies header has been written.
func (w *ResponseWriter) WroteHeader() bool {
	return w.wroteHeader
}

// Status returns response status code of the current request.
func (w *ResponseWriter) Status() int {
	return w.status
}

// Written returns number of bytes written in body.
func (w *ResponseWriter) Written() int64 {
	return w.written
}
