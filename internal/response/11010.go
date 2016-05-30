package response

import (
	"io"
	"net/http"
)

type responseWriter11010 struct {
	responseWriter01010
}

func (w *responseWriter11010) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
