// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package middleware

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/chanxuehong/gin"
)

// Logger instances a log middleware that will write the logs to os.Stdout.
func Logger() gin.HandlerFunc {
	return LoggerWithWriter(os.Stdout)
}

// LoggerWithWriter instances a log middleware with the specified writter buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(out io.Writer, notLogged ...string) gin.HandlerFunc {
	var skip map[string]struct{}

	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}

	return func(ctx *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		ctx.Next()

		// Log only when path is not being skipped
		path := ctx.Request.URL.Path
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			clientIP := ctx.ClientIP()
			method := ctx.Request.Method
			methodColor := colorForMethod(method)
			statusCode := ctx.ResponseWriter.Status()
			statusColor := colorForStatus(statusCode)

			fmt.Fprintf(out,
				"[GIN] %s |%s %3d %s| %16v | %s |%s  %s %-7s %s\r\n",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, __colorReset,
				latency,
				clientIP,
				methodColor, __colorReset, method, path,
			)
		}
	}
}
