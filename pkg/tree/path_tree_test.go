package tree

import "testing"

func TestPathTree(t *testing.T) {
	root := NewPathTree()

	path := "test/demo/"
	value := make(chan int)

	v := root.Put(path, value)
	if v != nil {
		t.Fatalf("error put")
	}

}

func TestPathTreePrefixMatch(t *testing.T) {
	tree := NewPathTree()

	test := "/test/demo"
	value2 := "test"
	tree.Put(test, value2)

	test3 := "/test/demo/name"
	value3 := "test3"
	tree.Put(test3, value3)

	if tree.PrefixMatch(test) != value2 {
		t.Fatalf("error prefixMatch")
	}

	if tree.PrefixMatch(test+"/123213/sadasd") != value2 {
		t.Fatalf("error prefixMatch")
	}

	if tree.PrefixMatch(test3) != value3 {
		t.Fatalf("error prefixMatch")
	}

	if tree.PrefixMatch(test3+"/asdasd/asdasd/asdasda") != value3 {
		t.Fatalf("error prefixMatch")
	}

	if tree.PrefixMatch("/") != nil {
		t.Fatalf("error prefixMatch")
	}

	if tree.PrefixMatch("/test") != nil {
		t.Fatalf("error prefixMatch")
	}

	if tree.PrefixMatch("/asd/testss") != nil {
		t.Fatalf("error prefixMatch")
	}

}
