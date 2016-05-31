package response

import (
	"net/http"
)

func newResponseWriter01000() ResponseWriter2 { return new(responseWriter01000) }

var _ stringWriter = (*responseWriter01000)(nil)

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
