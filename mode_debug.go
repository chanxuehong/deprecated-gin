// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build GIN_DEBUG

package gin

import (
	"log"
)

func debugPrintEngineNew() {}

func debugPrint(format string, arg ...interface{}) {
	log.Printf("[GIN_DEBUG] "+format, arg...)
}

func debugPrintError(err error) {
	if err != nil {
		debugPrint("[ERROR] %v\r\n", err)
	}
}

func debugPrintRoute(httpMethod, path string, handlers HandlerChain) {
	handlerName := nameOfFunction(handlers.last())
	debugPrint("%-5s %-25s --> %s (%d handlers)\r\n", httpMethod, path, handlerName, len(handlers))
}
