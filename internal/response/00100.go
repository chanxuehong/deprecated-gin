package response

import (
	"net/http"
)

func newResponseWriter00100() ResponseWriter2 { return new(responseWriter00100) }

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
