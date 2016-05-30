package response

import (
	"net/http"
)

type responseWriter00101 struct {
	responseWriter00001
}

func (w *responseWriter00101) Flush() {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	w.responseWriter.(http.Flusher).Flush()
}
