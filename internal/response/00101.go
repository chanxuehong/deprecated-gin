package response

import (
	"net/http"
)

func newResponseWriter00101() ResponseWriter2 { return new(responseWriter00101) }

var _ http.Flusher = (*responseWriter00101)(nil)
var _ http.CloseNotifier = (*responseWriter00101)(nil)

type responseWriter00101 struct {
	responseWriter00001
}

func (w *responseWriter00101) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.responseWriter.(http.Flusher).Flush()
}
