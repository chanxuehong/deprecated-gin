package response

import (
	"net/http"
)

var _ ResponseWriter = (*responseWriter01101)(nil)
var _ stringWriter = (*responseWriter01101)(nil)
var _ http.Flusher = (*responseWriter01101)(nil)
var _ http.CloseNotifier = (*responseWriter01101)(nil)

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
