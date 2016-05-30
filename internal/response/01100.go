package response

import (
	"net/http"
)

type responseWriter01100 struct {
	responseWriter00100
}

func (w *responseWriter01100) WriteString(data string) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(stringWriter).WriteString(data)
	w.written += int64(n)
	return
}
