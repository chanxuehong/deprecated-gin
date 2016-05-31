package response

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
	WroteHeader() bool // WroteHeader returns true if header has been written, otherwise false.
	Status() int       // Status returns the response status code of the current request.
	Written() int64    // Written returns number of bytes written in body.
}

type ResponseWriter2 interface {
	ResponseWriter
	Reset(w http.ResponseWriter)
}

func NewResponseWriter(bitmap int) ResponseWriter2 {
	newResponseWriterFunc := __newResponseWriterFuncArray[bitmap]
	return newResponseWriterFunc()
}

var __newResponseWriterFuncArray = [...]func() ResponseWriter2{
	newResponseWriter00000,
	newResponseWriter00001,
	newResponseWriter00010,
	newResponseWriter00011,
	newResponseWriter00100,
	newResponseWriter00101,
	newResponseWriter00110,
	newResponseWriter00111,
	newResponseWriter01000,
	newResponseWriter01001,
	newResponseWriter01010,
	newResponseWriter01011,
	newResponseWriter01100,
	newResponseWriter01101,
	newResponseWriter01110,
	newResponseWriter01111,
	newResponseWriter10000,
	newResponseWriter10001,
	newResponseWriter10010,
	newResponseWriter10011,
	newResponseWriter10100,
	newResponseWriter10101,
	newResponseWriter10110,
	newResponseWriter10111,
	newResponseWriter11000,
	newResponseWriter11001,
	newResponseWriter11010,
	newResponseWriter11011,
	newResponseWriter11100,
	newResponseWriter11101,
	newResponseWriter11110,
	newResponseWriter11111,
}
