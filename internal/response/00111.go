package response

import (
	"net/http"
)

func newResponseWriter00111() ResponseWriter2 { return new(responseWriter00111) }

var _ http.Flusher = (*responseWriter00111)(nil)
var _ http.Hijacker = (*responseWriter00111)(nil)
var _ http.CloseNotifier = (*responseWriter00111)(nil)

type responseWriter00111 struct {
	responseWriter00011
}

func (w *responseWriter00111) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.responseWriter.(http.Flusher).Flush()
}
