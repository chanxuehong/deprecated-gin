package response

import (
	"net/http"
)

type responseWriter01010 struct {
	responseWriter00010
}

func (w *responseWriter01010) WriteString(data string) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(stringWriter).WriteString(data)
	w.written += int64(n)
	return
}
