package response

import (
	"io"
	"net/http"
)

func newResponseWriter10111() ResponseWriter2 { return new(responseWriter10111) }

var _ io.ReaderFrom = (*responseWriter10111)(nil)
var _ http.Flusher = (*responseWriter10111)(nil)
var _ http.Hijacker = (*responseWriter10111)(nil)
var _ http.CloseNotifier = (*responseWriter10111)(nil)

type responseWriter10111 struct {
	responseWriter00111
}

func (w *responseWriter10111) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
