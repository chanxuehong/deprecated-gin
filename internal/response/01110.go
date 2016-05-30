package response

import (
	"net/http"
)

type responseWriter01110 struct {
	responseWriter00110
}

func (w *responseWriter01110) WriteString(data string) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(stringWriter).WriteString(data)
	w.written += int64(n)
	return
}
