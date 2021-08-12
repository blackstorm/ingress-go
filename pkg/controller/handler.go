package controller

import (
	"net/http"

	"github.com/blackstorm/ingress-go/pkg/controller/matcher"
)

type serverHandler struct {
	matcher *matcher.RequestMatcher
}

func newServerHandler(matcher *matcher.RequestMatcher) *serverHandler {
	return &serverHandler{
		matcher: matcher,
	}
}

func (s *serverHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("todo"))

	/*
		host := strings.Split(req.Host, ":")[0]
		if routes, ok := s.hostsRoutes[host]; !ok {
			resp.Write([]byte(fmt.Sprintf("TODO: no found host %s, is feedback no host backend", host)))
		} else {
			if route, ok := routes[req.URL.Path]; ok {
				route.handleRequest(resp, req)
			} else {
				resp.WriteHeader(404)
				resp.Write([]byte(fmt.Sprintf("host = %s path = %s not found", host, req.URL.Path)))
			}
		}
	*/
}
