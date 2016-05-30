package response

import (
	"log"
	"net/http"
)

type responseWriter00000 struct {
	responseWriter http.ResponseWriter
	wroteHeader    bool  // reply header has been written
	status         int   // status code passed to WriteHeader
	written        int64 // number of bytes written in body
}

func (w *responseWriter00000) reset(writer http.ResponseWriter) {
	w.responseWriter = writer
	w.wroteHeader = false
	w.status = http.StatusOK
	w.written = 0
}

func (w *responseWriter00000) WroteHeader() bool {
	return w.wroteHeader
}

func (w *responseWriter00000) Status() int {
	return w.status
}

func (w *responseWriter00000) Written() int64 {
	return w.written
}

func (w *responseWriter00000) Header() http.Header {
	return w.responseWriter.Header()
}

func (w *responseWriter00000) WriteHeader(code int) {
	if w.wroteHeader {
		log.Println("gin: multiple response.WriteHeader calls")
		return
	}
	w.wroteHeader = true
	w.status = code
	w.responseWriter.WriteHeader(code)
}

func (w *responseWriter00000) Write(data []byte) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.Write(data)
	w.written += int64(n)
	return
}
