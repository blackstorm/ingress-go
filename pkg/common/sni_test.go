package common

import "testing"

func TestToWildcardSni(t *testing.T) {

	// case 1
	example := "blog.example.com"
	wildcard := ToWildcardSni(example)
	answer := "*.example.com"
	if wildcard != answer {
		t.Fatalf("answer is %s", wildcard)
	}

	// case 2
	example = "example.com"
	wildcard = ToWildcardSni(example)
	answer = "*.com"
	if wildcard != answer {
		t.Fatalf("answer is %s", wildcard)
	}

	example = "blog.tech.example.com"
	wildcard = ToWildcardSni(example)
	answer = "*.tech.example.com"
	if wildcard != answer {
		t.Fatalf("answer is %s", wildcard)
	}
}
