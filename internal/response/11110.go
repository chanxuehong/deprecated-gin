package response

import (
	"io"
	"net/http"
)

func newResponseWriter11110() ResponseWriter2 { return new(responseWriter11110) }

var _ io.ReaderFrom = (*responseWriter11110)(nil)
var _ stringWriter = (*responseWriter11110)(nil)
var _ http.Flusher = (*responseWriter11110)(nil)
var _ http.Hijacker = (*responseWriter11110)(nil)

type responseWriter11110 struct {
	responseWriter01110
}

func (w *responseWriter11110) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
