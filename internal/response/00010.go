package response

import (
	"bufio"
	"net"
	"net/http"
)

type responseWriter00010 struct {
	responseWriter00000
}

func (w *responseWriter00010) Hijack() (rwc net.Conn, buf *bufio.ReadWriter, err error) {
	return w.responseWriter.(http.Hijacker).Hijack()
}
