package response

import (
	"net/http"
)

type responseWriter00110 struct {
	responseWriter00010
}

func (w *responseWriter00110) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.responseWriter.(http.Flusher).Flush()
}
