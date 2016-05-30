package response

import (
	"io"
	"net/http"
)

var _ ResponseWriter = (*responseWriter10011)(nil)
var _ io.ReaderFrom = (*responseWriter10011)(nil)
var _ http.Hijacker = (*responseWriter10011)(nil)
var _ http.CloseNotifier = (*responseWriter10011)(nil)

type responseWriter10011 struct {
	responseWriter00011
}

func (w *responseWriter10011) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
