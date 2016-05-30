package response

import (
	"io"
	"net/http"
)

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
