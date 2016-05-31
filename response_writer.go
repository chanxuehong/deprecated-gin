// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"

	"github.com/chanxuehong/gin/internal/response"
)

// type ResponseWriter interface {
//     http.ResponseWriter
//     WroteHeader() bool // WroteHeader returns true if header has been written, otherwise false.
//     Status() int       // Status returns the response status code of the current request.
//     Written() int64    // Written returns number of bytes written in body.
// }
type ResponseWriter response.ResponseWriter

type responseWriterArray [32]response.ResponseWriter2

func (cache *responseWriterArray) ResponseWriter2(w http.ResponseWriter) response.ResponseWriter2 {
	index := response.Bitmap(w)
	resp := cache[index]
	if resp == nil {
		resp = response.NewResponseWriter2(index)
		cache[index] = resp
	}
	resp.Reset(w)
	return resp
}
