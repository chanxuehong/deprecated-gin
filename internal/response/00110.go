package response

import (
	"net/http"
)

var _ ResponseWriter = (*responseWriter00110)(nil)
var _ http.Flusher = (*responseWriter00110)(nil)
var _ http.Hijacker = (*responseWriter00110)(nil)

type responseWriter00110 struct {
	responseWriter00010
}

func (w *responseWriter00110) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.responseWriter.(http.Flusher).Flush()
}
