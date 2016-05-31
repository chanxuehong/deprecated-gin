package response

import (
	"io"
	"net/http"
)

func newResponseWriter10000() ResponseWriter2 { return new(responseWriter10000) }

var _ io.ReaderFrom = (*responseWriter10000)(nil)

type responseWriter10000 struct {
	responseWriter00000
}

func (w *responseWriter10000) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
