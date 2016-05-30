package response

import (
	"net/http"
)

type responseWriter00100 struct {
	responseWriter00000
}

func (w *responseWriter00100) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.responseWriter.(http.Flusher).Flush()
}
