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

	root := "/"
	value := make(chan int)
	tree.Put(root, value)

	test := "/test/demo"
	value2 := make(chan int)
	tree.Put(test, value2)

	test3 := "/test/demo/name"
	value3 := make(chan int)
	tree.Put(test3, value3)

	if tree.PrefixMatch(test) != value2 {
		t.Fatalf("error prefixMatch")
	}
}

/*
func TestPathTreeMatch(t *testing.T) {
	root := NewPathTree()

	path := "/test/demo"
	value := make(chan int)
	root.Put(path, value)

	path2 := "/test/dem1o"
	value2 := make(chan int)
	root.Put(path2, value2)

	v := root.Match(path)
	if v == nil {
		t.Fatalf("no found path")
	} else {
		if v != value {
			t.Fatalf("value not match")
		}
	}

	v = root.Match(path2)
	if v == nil {
		t.Fatalf("no found path2")
	} else {
		if v != value2 {
			t.Fatalf("value2 not match")
		}
	}
}

func TestPathTreeDelete(t *testing.T) {
	root := NewPathTree()

	path := "/test/demo"
	value := make(chan int)
	root.Put(path, value)

	path2 := "/test/demo1"
	value2 := make(chan int)
	root.Put(path2, value2)

	path3 := "/test/demo2"
	value3 := make(chan int)
	root.Put(path3, value3)

	path4 := "/test/demo3"
	value4 := make(chan int)
	root.Put(path4, value4)

	root.Delete("/test/demo2")

	if v := root.Match("/test/demo2"); v != nil {
		t.Fatalf("delete failed")
	}

}
*/
