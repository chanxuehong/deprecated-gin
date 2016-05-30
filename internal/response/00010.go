package response

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

var _ ResponseWriter = (*responseWriter00010)(nil)
var _ http.Hijacker = (*responseWriter00010)(nil)

type responseWriter00010 struct {
	responseWriter00000
	hijacked bool
}

func (w *responseWriter00010) reset(writer http.ResponseWriter) {
	w.responseWriter00000.reset(writer)
	w.hijacked = false
}

func (w *responseWriter00010) WriteHeader(code int) {
	if w.hijacked {
		log.Println("gin: response.WriteHeader on hijacked connection")
		return
	}
	w.responseWriter00000.WriteHeader(code)
}

func (w *responseWriter00010) Write(data []byte) (n int, err error) {
	if w.hijacked {
		log.Println("gin: response.Write on hijacked connection")
		return 0, http.ErrHijacked
	}
	return w.responseWriter00000.Write(data)
}

func (w *responseWriter00010) Hijack() (rwc net.Conn, buf *bufio.ReadWriter, err error) {
	rwc, buf, err = w.responseWriter.(http.Hijacker).Hijack()
	if err == nil {
		w.hijacked = true
	}
	return
}
