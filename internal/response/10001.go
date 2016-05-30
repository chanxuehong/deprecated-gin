package response

import (
	"io"
	"net/http"
)

type responseWriter10001 struct {
	responseWriter00001
}

func (w *responseWriter10001) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
