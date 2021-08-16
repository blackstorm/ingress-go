package tree

import (
	"strings"
)

// TODO use https://pkg.go.dev/golang.org/x/exp/utf8string for unicode support
type PathTree struct {
	value    interface{}
	children map[string]*PathTree
}

func NewPathTree() *PathTree {
	return &PathTree{
		value:    nil,
		children: make(map[string]*PathTree),
	}
}

func newPathTree(value interface{}) *PathTree {
	return &PathTree{
		value:    value,
		children: make(map[string]*PathTree),
	}
}

func (t *PathTree) Put(path string, value interface{}) interface{} {
	if path[0] == '/' {
		path = path[1:]
	}
	paths := strings.Split(path, "/")
	return t.put(paths, value, t.children)
}

func (t *PathTree) put(paths []string, value interface{}, tree map[string]*PathTree) interface{} {
	isLast := len(paths) == 1
	root := paths[0]

	var node *PathTree
	var ok bool
	node, ok = tree[root]

	if !ok {
		node = newPathTree(nil)
		tree[root] = node
	}

	if isLast {
		nodeOldValue := node.value
		node.value = value
		return nodeOldValue
	}

	return t.put(paths[1:], value, node.children)
}

func (m *PathTree) Match(path string) interface{} {
	if path[0] == '/' {
		path = path[1:]
	}
	paths := strings.Split(path, "/")
	return m.match(paths, m.children, nil)
}

func (m *PathTree) match(paths []string, tree map[string]*PathTree, best interface{}) interface{} {
	if tree == nil {
		return nil
	}

	root := paths[0]
	isLast := len(paths) == 1

	var node *PathTree
	var ok bool
	node, ok = tree[root]

	if isLast {
		if ok {
			return node.value
		} else {
			return nil
		}
	}

	res := m.match(paths[1:], node.children, node.value)
	if res == nil {
		return best
	}
	return res
}
