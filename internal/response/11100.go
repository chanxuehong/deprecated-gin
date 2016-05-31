package response

import (
	"io"
	"net/http"
)

func newResponseWriter11100() ResponseWriter2 { return new(responseWriter11100) }

var _ io.ReaderFrom = (*responseWriter11100)(nil)
var _ stringWriter = (*responseWriter11100)(nil)
var _ http.Flusher = (*responseWriter11100)(nil)

type responseWriter11100 struct {
	responseWriter01100
}

func (w *responseWriter11100) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
