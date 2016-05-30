package response

import (
	"net/http"
)

type responseWriter01011 struct {
	responseWriter00011
}

func (w *responseWriter01011) WriteString(data string) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(stringWriter).WriteString(data)
	w.written += int64(n)
	return
}
