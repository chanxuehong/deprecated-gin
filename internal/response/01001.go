package response

import (
	"net/http"
)

type responseWriter01001 struct {
	responseWriter00001
}

func (w *responseWriter01001) WriteString(data string) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(stringWriter).WriteString(data)
	w.written += int64(n)
	return
}
