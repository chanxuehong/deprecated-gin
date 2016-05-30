package response

import (
	"io"
	"net/http"
)

var _ ResponseWriter = (*responseWriter10101)(nil)
var _ io.ReaderFrom = (*responseWriter10101)(nil)
var _ http.Flusher = (*responseWriter10101)(nil)
var _ http.CloseNotifier = (*responseWriter10101)(nil)

type responseWriter10101 struct {
	responseWriter00101
}

func (w *responseWriter10101) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
