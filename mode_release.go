// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build !GIN_DEBUG

package gin

import (
	"log"
)

func debugPrintEngineNew() {
	const str = "[NOTE] Running in \"release\" mode. Switch to \"debug\" mode using:\r\n" +
		"        go build -tags GIN_DEBUG [your project]\r\n\r\n"
	log.Print(str)
}

func debugPrintf(format string, arg ...interface{}) {}

func debugPrintError(err error) {}

func debugPrintRoute(httpMethod, path string, handlers HandlerChain) {}
