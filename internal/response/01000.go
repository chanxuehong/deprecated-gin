package response

import (
	"net/http"
)

type responseWriter01000 struct {
	responseWriter00000
}

func (w *responseWriter01000) WriteString(data string) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(stringWriter).WriteString(data)
	w.written += int64(n)
	return
}
