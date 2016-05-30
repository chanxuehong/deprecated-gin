package response

import (
	"net/http"
)

var _ ResponseWriter = (*responseWriter00001)(nil)
var _ http.CloseNotifier = (*responseWriter00001)(nil)

type responseWriter00001 struct {
	responseWriter00000
}

func (w *responseWriter00001) CloseNotify() <-chan bool {
	return w.responseWriter.(http.CloseNotifier).CloseNotify()
}
