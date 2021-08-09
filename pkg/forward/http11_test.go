package forward

import (
	"bytes"
	"net/http"
	"testing"
)

func TestFastHttp11(t *testing.T) {
	req, err := http.NewRequest("get", "https://baidu.com", bytes.NewReader([]byte{}))
	if err != nil {
		t.Fatal(err)
	}

	res, err := FastHttp11("https://baidu.com", req)
	if err != nil {
		t.Fatal(err)
	}

	if res == nil {
		t.Fatal("res is nil")
	}
}
