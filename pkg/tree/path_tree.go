package tree

import (
	"strings"
)

// TODO use https://pkg.go.dev/golang.org/x/exp/utf8string for unicode support
type PathTree struct {
	path     string
	value    interface{}
	children map[string]*PathTree
}

func NewPathTree() *PathTree {
	return &PathTree{
		children: make(map[string]*PathTree),
	}
}

func newPathTree(path string) *PathTree {
	return &PathTree{
		path:     path,
		children: make(map[string]*PathTree),
	}
}

func treePaths(path string) []string {
	length := len(path)

	if length == 0 || path == "/" {
		return []string{""}
	}

	if path[0] != '/' {
		path = "/" + path
		length++
	}

	if path[length-1] == '/' {
		path = path[:length-1]
	}

	return strings.Split(path, "/")
}

// The ingress path is validated
func (t *PathTree) Put(path string, value interface{}) interface{} {
	return t.put(treePaths(path), value)
}

func (t *PathTree) put(paths []string, value interface{}) interface{} {
	// return recursion
	if len(paths) == 1 {
		oldValue := t.value
		t.value = value
		return oldValue
	}

	subPath := paths[1]
	if t.children[subPath] == nil {
		t.children[subPath] = newPathTree(subPath)
	}

	return t.children[subPath].put(paths[1:], value)
}

func (t *PathTree) PrefixMatch(path string) interface{} {
	return t.prefixMatch(treePaths(path))
}

func (t *PathTree) prefixMatch(paths []string) interface{} {
	if len(paths) == 1 {
		return t.value
	}

	if tree := t.children[paths[1]]; tree != nil {
		if res := tree.prefixMatch(paths[1:]); res != nil {
			return res
		}
	}

	return t.value
}

func (t *PathTree) Delete(path string) {
	t.delete(treePaths(path))
}

func (t *PathTree) delete(paths []string) {
	if len(paths) == 1 {
		t.value = nil
		return
	}

	if tree := t.children[paths[1]]; tree != nil {
		tree.delete(paths[1:])
	}
}
