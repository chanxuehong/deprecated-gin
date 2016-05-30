package response

import (
	"net/http"
)

type responseWriter01101 struct {
	responseWriter00101
}

func (w *responseWriter01101) WriteString(data string) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(stringWriter).WriteString(data)
	w.written += int64(n)
	return
}
