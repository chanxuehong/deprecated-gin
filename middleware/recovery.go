// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http/httputil"
	"os"
	"runtime"

	"github.com/chanxuehong/gin"
)

var (
	__dunno     = []byte("???")
	__centerDot = []byte("·")
	__dot       = []byte(".")
	__slash     = []byte("/")
	__newline   = []byte("\n")
)

// Recovery returns a middleware that recovers from any panics
// and writes a 500 to os.Stderr if there was one.
func Recovery() gin.HandlerFunc {
	return RecoveryWithWriter(os.Stderr)
}

// RecoveryWithWriter returns a middleware that recovers from any panics
// and writes a 500 to the specified io.Writer if there was one.
//
// Example: os.Stdout, a file opened in write mode, a socket...
func RecoveryWithWriter(out io.Writer) gin.HandlerFunc {
	var logger *log.Logger
	if out != nil {
		logger = log.New(out, "\n\n\x1b[31m", log.LstdFlags)
	}
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				if logger != nil {
					stack := stack(3)
					httprequest, _ := httputil.DumpRequest(ctx.Request, false)
					logger.Printf("[Recovery] panic recovered:\n%s\n%s\n%s%s", string(httprequest), err, stack, __colorReset)
				}
				ctx.AbortWithStatus(500)
			}
		}()
		ctx.Next()
	}
}

// stack returns a nicely formated stack frame, skipping skip frames
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least. If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, __newline)
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return __dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return __dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included. Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastslash := bytes.LastIndex(name, __slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}
	if period := bytes.Index(name, __dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, __centerDot, __dot, -1)
	return name
}
