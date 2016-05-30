package response

import (
	"net/http"
)

type responseWriter00111 struct {
	responseWriter00011
}

func (w *responseWriter00111) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.responseWriter.(http.Flusher).Flush()
}
