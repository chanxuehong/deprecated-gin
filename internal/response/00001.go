package response

import (
	"net/http"
)

func newResponseWriter00001() ResponseWriter2 { return new(responseWriter00001) }

var _ http.CloseNotifier = (*responseWriter00001)(nil)

type responseWriter00001 struct {
	responseWriter00000
}

func (w *responseWriter00001) CloseNotify() <-chan bool {
	return w.responseWriter.(http.CloseNotifier).CloseNotify()
}
