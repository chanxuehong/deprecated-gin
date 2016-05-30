package response

import (
	"io"
	"net/http"
)

var _ ResponseWriter = (*responseWriter11101)(nil)
var _ io.ReaderFrom = (*responseWriter11101)(nil)
var _ stringWriter = (*responseWriter11101)(nil)
var _ http.Flusher = (*responseWriter11101)(nil)
var _ http.CloseNotifier = (*responseWriter11101)(nil)

type responseWriter11101 struct {
	responseWriter01101
}

func (w *responseWriter11101) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
