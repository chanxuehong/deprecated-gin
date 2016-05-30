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

func (arr *ResponseWriterArray) ResponseWriter(w http.ResponseWriter) ResponseWriter {
	index := bitmap(w)
	resp := arr[index]
	if resp == nil {
		resp = newResponseWriter(index)
		arr[index] = resp
	}
	resp.reset(w)
	return resp
}

func newResponseWriter(bitmap int) ResponseWriter {
	switch bitmap {
	default:
		return new(responseWriter00000)
	case 0x1f:
		return new(responseWriter11111)
	case 0x00:
		return new(responseWriter00000)
	case 0x01:
		return new(responseWriter00001)
	case 0x02:
		return new(responseWriter00010)
	case 0x03:
		return new(responseWriter00011)
	case 0x04:
		return new(responseWriter00100)
	case 0x05:
		return new(responseWriter00101)
	case 0x06:
		return new(responseWriter00110)
	case 0x07:
		return new(responseWriter00111)
	case 0x08:
		return new(responseWriter01000)
	case 0x09:
		return new(responseWriter01001)
	case 0x0a:
		return new(responseWriter01010)
	case 0x0b:
		return new(responseWriter01011)
	case 0x0c:
		return new(responseWriter01100)
	case 0x0d:
		return new(responseWriter01101)
	case 0x0e:
		return new(responseWriter01110)
	case 0x0f:
		return new(responseWriter01111)
	case 0x10:
		return new(responseWriter10000)
	case 0x11:
		return new(responseWriter10001)
	case 0x12:
		return new(responseWriter10010)
	case 0x13:
		return new(responseWriter10011)
	case 0x14:
		return new(responseWriter10100)
	case 0x15:
		return new(responseWriter10101)
	case 0x16:
		return new(responseWriter10110)
	case 0x17:
		return new(responseWriter10111)
	case 0x18:
		return new(responseWriter11000)
	case 0x19:
		return new(responseWriter11001)
	case 0x1a:
		return new(responseWriter11010)
	case 0x1b:
		return new(responseWriter11011)
	case 0x1c:
		return new(responseWriter11100)
	case 0x1d:
		return new(responseWriter11101)
	case 0x1e:
		return new(responseWriter11110)
	}
}
