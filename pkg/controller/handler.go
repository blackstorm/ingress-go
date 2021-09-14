package controller

import (
	"net/http"

	"github.com/blackstorm/ingress-go/pkg/controller/matcher"
)

type handler func(resp http.ResponseWriter, req *http.Request)

type serverHandler struct {
	matcher        *matcher.RequestMatcher
	beforeHandlers []handler
	afterHandlers  []handler
}

func newServerHandler(matcher *matcher.RequestMatcher) *serverHandler {
	return &serverHandler{
		matcher:        matcher,
		beforeHandlers: make([]handler, 0),
		afterHandlers:  make([]handler, 0),
	}
}

func (s *serverHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	s.doBefores(resp, req)
	resp.Write([]byte("todo"))
	s.doAfters(resp, req)
}

func (s *serverHandler) doBefores(resp http.ResponseWriter, req *http.Request) {
	for _, hr := range s.beforeHandlers {
		hr(resp, req)
	}
}

func (s *serverHandler) doAfters(resp http.ResponseWriter, req *http.Request) {
	for _, hr := range s.afterHandlers {
		hr(resp, req)
	}
}

func (s *serverHandler) addBeforeHandleRequest(hr handler) {
	s.beforeHandlers = append(s.beforeHandlers, hr)
}

func (s *serverHandler) addAfterHandleRequest(hr handler) {
	s.afterHandlers = append(s.afterHandlers, hr)
}
