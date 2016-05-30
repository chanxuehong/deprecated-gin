package response

import (
	"io"
	"net/http"
)

type responseWriter10000 struct {
	responseWriter00000
}

func (w *responseWriter10000) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
