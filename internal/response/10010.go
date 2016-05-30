package response

import (
	"io"
	"net/http"
)

var _ ResponseWriter = (*responseWriter10010)(nil)
var _ io.ReaderFrom = (*responseWriter10010)(nil)
var _ http.Hijacker = (*responseWriter10010)(nil)

type responseWriter10010 struct {
	responseWriter00010
}

func (w *responseWriter10010) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
