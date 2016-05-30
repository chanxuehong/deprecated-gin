package response

import (
	"io"
	"net/http"
)

// stringWriter is the interface that wraps the WriteString method.
type stringWriter interface {
	WriteString(s string) (n int, err error)
}

// bitmap = (io.ReaderFrom, stringWriter, http.Flusher, http.Hijacker, http.CloseNotifier)
const (
	httpCloseNotifierBitmap = 1 << iota // http.CloseNotifier
	httpHijackerBitmap                  // http.Hijacker
	httpFlusherBitmap                   // http.Flusher
	stringWriterBitmap                  // stringWriter
	ioReaderFromBitmap                  // io.ReaderFrom
)

func index(w http.ResponseWriter) int {
	var ok bool
	bitmap := 0x00
	if _, ok = w.(io.ReaderFrom); ok {
		bitmap |= ioReaderFromBitmap
	}
	if _, ok = w.(stringWriter); ok {
		bitmap |= stringWriterBitmap
	}
	if _, ok = w.(http.Flusher); ok {
		bitmap |= httpFlusherBitmap
	}
	if _, ok = w.(http.Hijacker); ok {
		bitmap |= httpHijackerBitmap
	}
	if _, ok = w.(http.CloseNotifier); ok {
		bitmap |= httpCloseNotifierBitmap
	}
	return bitmap
}
