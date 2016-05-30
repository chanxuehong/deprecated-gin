package response

import (
	"log"
	"net/http"
)

type responseWriter01110 struct {
	responseWriter00110
}

func (w *responseWriter01110) WriteString(data string) (n int, err error) {
	if w.responseWriter00110.hijacked {
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
