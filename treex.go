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
func (params Params) ByName(name string) string {
	for i := 0; i < len(params); i++ {
		if params[i].Key == name {
			return params[i].Value
		}
	}
	return ""
}

// ByIndex returns the value with given index.
// If the index out of range, an empty string is returned.
func (params Params) ByIndex(index int) string {
	if index < 0 || index >= len(params) {
		return ""
	}
	return params[index].Value
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

func (p *trees) addTree(method string, root *node) {
	*p = append(*p, tree{method: method, root: root})
}

type Route struct {
	Method  string
	Path    string
	Handler string // handler name
}

// routes returns a slice of registered routes, including some useful information, such as:
// the http method, path and the handler name.
func (v trees) routes() []Route {
	routes := make([]Route, 0, 64)
	for i, l := 0, len(v); i < l; i++ {
		iterateTree(&routes, v[i].root, "", v[i].method)
	}
	return routes
}

func iterateTree(routesPtr *[]Route, root *node, pathPrefix, treeMethod string) {
	path := pathPrefix + root.path
	if len(root.handlers) > 0 {
		*routesPtr = append(*routesPtr, Route{
			Method:  treeMethod,
			Path:    path,
			Handler: nameOfFunction(root.handlers.last()),
		})
	}
	for _, childRoot := range root.children {
		iterateTree(routesPtr, childRoot, path, treeMethod)
	}
}
