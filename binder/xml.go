// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binder

import (
	"encoding/xml"
	"net/http"
)

type xmlBinder int

func (*xmlBinder) Bind(req *http.Request, obj interface{}) error {
	return xml.NewDecoder(req.Body).Decode(obj)
}
