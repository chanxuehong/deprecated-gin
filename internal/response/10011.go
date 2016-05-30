package response

import (
	"io"
	"net/http"
)

type responseWriter10011 struct {
	responseWriter00011
}

func (w *responseWriter10011) ReadFrom(r io.Reader) (n int64, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(io.ReaderFrom).ReadFrom(r)
	w.written += n
	return
}
