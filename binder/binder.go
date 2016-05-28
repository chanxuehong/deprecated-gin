// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binder

import (
	"net/http"
)

type Binder interface {
	Bind(*http.Request, interface{}) error
}
