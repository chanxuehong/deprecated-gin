package response

import (
	"io"
	"net/http"
)

func newResponseWriter10110() ResponseWriter2 { return new(responseWriter10110) }

var _ io.ReaderFrom = (*responseWriter10110)(nil)
var _ http.Flusher = (*responseWriter10110)(nil)
var _ http.Hijacker = (*responseWriter10110)(nil)

type responseWriter10110 struct {
	responseWriter00110
}

func (w *responseWriter10110) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
