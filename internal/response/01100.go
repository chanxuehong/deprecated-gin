package response

import (
	"net/http"
)

var _ ResponseWriter = (*responseWriter01100)(nil)
var _ stringWriter = (*responseWriter01100)(nil)
var _ http.Flusher = (*responseWriter01100)(nil)

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
