// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"sync/atomic"
)

const (
	startCheckerInitialValue = uintptr(0)
	startCheckerStartedValue = ^uintptr(0)
)

type startChecker uintptr

func (p *startChecker) start() {
	if uintptr(*p) != startCheckerStartedValue {
		atomic.CompareAndSwapUintptr((*uintptr)(p), startCheckerInitialValue, startCheckerStartedValue)
	}
}

func (v startChecker) check() {
	if uintptr(v) != startCheckerInitialValue {
		panic("The service has been started.")
	}
}
