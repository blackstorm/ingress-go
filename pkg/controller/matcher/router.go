package matcher

import (
	"errors"

	"github.com/blackstorm/ingress-go/pkg/controller/backend"
	"github.com/blackstorm/ingress-go/pkg/controller/label"
	netv1 "k8s.io/api/networking/v1"
)

type router struct {
	paths map[string]*backend.Backend
}

func newRouter() *router {
	return &router{
		paths: make(map[string]*backend.Backend),
	}
}

func (r *router) get(path string) *backend.Backend {
	return r.paths[path]
}

func (r *router) add(path string, backend *backend.Backend) {
	r.paths[path] = backend
}

type RouterMatcher struct {
	routes map[string]*router
}

func NewRouterMatcher() *RouterMatcher {
	return &RouterMatcher{
		routes: make(map[string]*router),
	}
}

func (m *RouterMatcher) match(host string, path string) (*backend.Backend, error) {
	if router, ok := m.routes[host]; ok {
		return router.get(path), nil
	}
	return nil, errors.New("router no found")
}

func (m *RouterMatcher) add(rule netv1.IngressRule, labeler label.Labeler) {
	host := rule.Host

	router, ok := m.routes[host]
	if !ok {
		router = newRouter()
		m.routes[host] = router
	}

	for _, path := range rule.HTTP.Paths {
		backend := backend.NewBackend()
		router.add(path.Path, backend)
		labeler.LabelResource(backend)
	}
}
