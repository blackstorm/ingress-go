package controller

import (
	"net/http"
)

type serverHandler struct {
	matcher *matcher
}

func newServerHandler(matcher *matcher) *serverHandler {
	handler := &serverHandler{
		matcher: matcher,
	}
	return handler
}

func (s *serverHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
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
