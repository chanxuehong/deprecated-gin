package response

import (
	"io"
	"net/http"
)

var _ ResponseWriter = (*responseWriter11111)(nil)
var _ io.ReaderFrom = (*responseWriter11111)(nil)
var _ stringWriter = (*responseWriter11111)(nil)
var _ http.Flusher = (*responseWriter11111)(nil)
var _ http.Hijacker = (*responseWriter11111)(nil)
var _ http.CloseNotifier = (*responseWriter11111)(nil)

type responseWriter11111 struct {
	responseWriter01111
}

func (w *responseWriter11111) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
