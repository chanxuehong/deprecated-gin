package response

import (
	"io"
	"net/http"
)

var _ ResponseWriter = (*responseWriter10100)(nil)
var _ io.ReaderFrom = (*responseWriter10100)(nil)
var _ http.Flusher = (*responseWriter10100)(nil)

type responseWriter10100 struct {
	responseWriter00100
}

func (w *responseWriter10100) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
