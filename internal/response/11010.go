package response

import (
	"io"
	"net/http"
)

func newResponseWriter11010() ResponseWriter2 { return new(responseWriter11010) }

var _ io.ReaderFrom = (*responseWriter11010)(nil)
var _ stringWriter = (*responseWriter11010)(nil)
var _ http.Hijacker = (*responseWriter11010)(nil)

type responseWriter11010 struct {
	responseWriter01010
}

func (w *responseWriter11010) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
