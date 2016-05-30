package response

import (
	"log"
	"net/http"
)

var _ ResponseWriter = (*responseWriter01011)(nil)
var _ stringWriter = (*responseWriter01011)(nil)
var _ http.Hijacker = (*responseWriter01011)(nil)
var _ http.CloseNotifier = (*responseWriter01011)(nil)

type responseWriter01011 struct {
	responseWriter00011
}

func (w *responseWriter01011) WriteString(data string) (n int, err error) {
	if w.responseWriter00011.hijacked {
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
