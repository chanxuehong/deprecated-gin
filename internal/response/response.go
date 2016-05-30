package response

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
	WroteHeader() bool // WroteHeader returns true if header has been written, otherwise false.
	Status() int       // Status returns the response status code of the current request.
	Written() int64    // Written returns number of bytes written in body.
	reset(w http.ResponseWriter)
}

type ResponseWriterArray [32]ResponseWriter

func NewResponseWriterArray() ResponseWriterArray {
	return ResponseWriterArray{
		&responseWriter00000{},
		&responseWriter00001{},
		&responseWriter00010{},
		&responseWriter00011{},
		&responseWriter00100{},
		&responseWriter00101{},
		&responseWriter00110{},
		&responseWriter00111{},
		&responseWriter01000{},
		&responseWriter01001{},
		&responseWriter01010{},
		&responseWriter01011{},
		&responseWriter01100{},
		&responseWriter01101{},
		&responseWriter01110{},
		&responseWriter01111{},
		&responseWriter10000{},
		&responseWriter10001{},
		&responseWriter10010{},
		&responseWriter10011{},
		&responseWriter10100{},
		&responseWriter10101{},
		&responseWriter10110{},
		&responseWriter10111{},
		&responseWriter11000{},
		&responseWriter11001{},
		&responseWriter11010{},
		&responseWriter11011{},
		&responseWriter11100{},
		&responseWriter11101{},
		&responseWriter11110{},
		&responseWriter11111{},
	}
}

func (arr *ResponseWriterArray) SelectResponseWriter(w http.ResponseWriter) ResponseWriter {
	resp := arr[index(w)]
	resp.reset(w)
	return resp
}
