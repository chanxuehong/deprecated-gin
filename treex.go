// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

// Param is a single URL parameter, consisting of a key and a value.
type Param struct {
	Key   string
	Value string // May be empty
}

// Params is a Param-slice, as returned by the router.
// The slice is ordered, the first URL parameter is also the first slice value.
// It is therefore safe to read values by the index.
type Params []Param

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (params Params) ByName(name string) (value string) {
	for i := 0; i < len(params); i++ {
		if params[i].Key == name {
			return params[i].Value
		}
	}
	return
}

// ByIndex returns the value with given index.
// If the index out of range, an empty string is returned.
func (params Params) ByIndex(index int) (value string) {
	if index >= 0 && index < len(params) {
		return params[index].Value
	}
	return
}

type tree struct {
	method string
	root   *node
}

type trees []tree

func (v trees) getTree(method string) *node {
	for i, l := 0, len(v); i < l; i++ {
		if v[i].method == method {
			return v[i].root
		}
	}
	return nil
}

func (v trees) addTree(method string, root *node) trees {
	return append(v, tree{method: method, root: root})
}

type Route struct {
	Method  string
	Path    string
	Handler string // handler name
}

// routes returns a slice of registered routes, including some useful information, such as:
// the http method, path and the handler name.
func (v trees) routes() (routes []Route) {
	for i, l := 0, len(v); i < l; i++ {
		routes = iterateTree(routes, v[i].root, "", v[i].method)
	}
	return routes
}

func iterateTree(routes []Route, root *node, pathPrefix, treeMethod string) []Route {
	path := pathPrefix + root.path
	if len(root.handlers) > 0 {
		routes = append(routes, Route{
			Method:  treeMethod,
			Path:    path,
			Handler: nameOfFunction(root.handlers.last()),
		})
	}
	for _, childRoot := range root.children {
		routes = iterateTree(routes, childRoot, path, treeMethod)
	}
	return routes
}
