package response

import (
	"net/http"
)

var _ ResponseWriter = (*responseWriter00100)(nil)
var _ http.Flusher = (*responseWriter00100)(nil)

type responseWriter00100 struct {
	responseWriter00000
}

func (w *responseWriter00100) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.responseWriter.(http.Flusher).Flush()
}
