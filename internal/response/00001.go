package response

import (
	"net/http"
)

type responseWriter00001 struct {
	responseWriter00000
}

func (w *responseWriter00001) CloseNotify() <-chan bool {
	return w.responseWriter.(http.CloseNotifier).CloseNotify()
}
