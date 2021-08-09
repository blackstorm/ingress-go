package forward

import (
	"io/ioutil"
	"net/http"

	"github.com/valyala/fasthttp"
)

// http1.1
func FastHttp11(url string, req *http.Request) ([]byte, error) {
	// fasthttpproxy.FasthttpHTTPDialerTimeout()
	// acquire request and response
	freq := fasthttp.AcquireRequest()
	fresp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(freq)
		fasthttp.ReleaseResponse(fresp)
	}()

	var body []byte
	var err error
	// read request body
	if req.Body != http.NoBody {
		body, err = ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			return nil, err
		}
	}

	// set request body
	freq.Header.SetMethod(req.Method)
	freq.SetRequestURI(url)
	if body != nil {
		freq.SetBody(body)
	}

	// send request
	if err = fasthttp.Do(freq, fresp); err != nil {
		return nil, err
	} else {
		return fresp.Body(), nil
	}
}
