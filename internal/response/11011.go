package response

import (
	"io"
	"net/http"
)

func newResponseWriter11011() ResponseWriter2 { return new(responseWriter11011) }

var _ io.ReaderFrom = (*responseWriter11011)(nil)
var _ stringWriter = (*responseWriter11011)(nil)
var _ http.Hijacker = (*responseWriter11011)(nil)
var _ http.CloseNotifier = (*responseWriter11011)(nil)

type responseWriter11011 struct {
	responseWriter01011
}

func (w *responseWriter11011) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
