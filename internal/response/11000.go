package response

import (
	"io"
	"net/http"
)

var _ ResponseWriter = (*responseWriter11000)(nil)
var _ io.ReaderFrom = (*responseWriter11000)(nil)
var _ stringWriter = (*responseWriter11000)(nil)

type responseWriter11000 struct {
	responseWriter01000
}

func (w *responseWriter11000) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
