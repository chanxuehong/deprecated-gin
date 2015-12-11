// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"net/http"
	"os"
)

var _ http.FileSystem = Dir("")

// Dir like http.Dir but the http.File returned by Dir.Open() does not list the directory files,
// so we can pass Dir to http.FileServer() to prevent it lists the directory files.
type Dir string

func (d Dir) Open(name string) (http.File, error) {
	file, err := http.Dir(d).Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{file}, nil
}

type neuteredReaddirFile struct{ http.File }

// Overrides the http.File default implementation
func (neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) { return nil, nil }
