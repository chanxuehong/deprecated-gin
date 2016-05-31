package response

import (
	"io"
	"net/http"
)

func newResponseWriter11001() ResponseWriter2 { return new(responseWriter11001) }

var _ io.ReaderFrom = (*responseWriter11001)(nil)
var _ stringWriter = (*responseWriter11001)(nil)
var _ http.CloseNotifier = (*responseWriter11001)(nil)

type responseWriter11001 struct {
	responseWriter01001
}

func (w *responseWriter11001) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
