// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"sync/atomic"
)

const (
	startedCheckerInitialValue = uintptr(0)
	startedCheckerStartedValue = ^uintptr(0)
)

type startedChecker uintptr

func (p *startedChecker) start() {
	if uintptr(*p) != startedCheckerStartedValue {
		atomic.CompareAndSwapUintptr((*uintptr)(p), startedCheckerInitialValue, startedCheckerStartedValue)
	}
}

func (v startedChecker) check() {
	if uintptr(v) != startedCheckerInitialValue {
		panic("the service has been started.")
	}
}
