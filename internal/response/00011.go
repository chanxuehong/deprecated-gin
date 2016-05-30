package response

import (
	"bufio"
	"net"
	"net/http"
)

type responseWriter00011 struct {
	responseWriter00001
}

func (w *responseWriter00011) Hijack() (rwc net.Conn, buf *bufio.ReadWriter, err error) {
	return w.responseWriter.(http.Hijacker).Hijack()
}
