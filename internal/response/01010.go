package response

import (
	"log"
	"net/http"
)

var _ ResponseWriter = (*responseWriter01010)(nil)
var _ stringWriter = (*responseWriter01010)(nil)
var _ http.Hijacker = (*responseWriter01010)(nil)

type responseWriter01010 struct {
	responseWriter00010
}

func (w *responseWriter01010) WriteString(data string) (n int, err error) {
	if w.responseWriter00010.hijacked {
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
