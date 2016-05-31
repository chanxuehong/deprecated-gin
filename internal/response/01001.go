package response

import (
	"net/http"
)

func newResponseWriter01001() ResponseWriter2 { return new(responseWriter01001) }

var _ stringWriter = (*responseWriter01001)(nil)
var _ http.CloseNotifier = (*responseWriter01001)(nil)

type responseWriter01001 struct {
	responseWriter00001
}

func (w *responseWriter01001) WriteString(data string) (n int, err error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.responseWriter.(stringWriter).WriteString(data)
	w.written += int64(n)
	return
}
