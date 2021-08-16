package tree

import "testing"

func TestPathTree(t *testing.T) {
	root := NewPathTree()

	path := "/test/demo"
	value := make(chan int)

	v := root.Put(path, value)
	if v != nil {
		t.Fatalf("error put")
	}

}

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
