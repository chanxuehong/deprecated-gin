package response

import (
	"log"
	"net/http"
)

var _ ResponseWriter = (*responseWriter01111)(nil)
var _ stringWriter = (*responseWriter01111)(nil)
var _ http.Flusher = (*responseWriter01111)(nil)
var _ http.Hijacker = (*responseWriter01111)(nil)
var _ http.CloseNotifier = (*responseWriter01111)(nil)

type responseWriter01111 struct {
	responseWriter00111
}

func (w *responseWriter01111) WriteString(data string) (n int, err error) {
	if w.responseWriter00111.hijacked {
		log.Println("gin: response.WriteString on hijacked connection")
		return 0, http.ErrHijacked
	}
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(stringWriter).WriteString(data)
	w.written += int64(n)
	return
}
