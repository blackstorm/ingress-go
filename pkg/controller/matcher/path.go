package matcher

import (
	"errors"

	"github.com/blackstorm/ingress-go/pkg/controller/backend"
	"github.com/blackstorm/ingress-go/pkg/controller/label"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/klog/v2"
)

type PathMatcher struct {
	paths map[string]*backend.Backend
}

func newPathMatcher() *PathMatcher {
	return &PathMatcher{
		paths: make(map[string]*backend.Backend),
	}
}

func (m *PathMatcher) match(path string) (*backend.Backend, error) {
	if backend, ok := m.paths[path]; ok {
		return backend, nil
	}
	return nil, errors.New("path no found")
}

func (m *PathMatcher) add(ingress *netv1.Ingress, rule netv1.IngressRule, labeler label.Labeler) {
	for _, path := range rule.HTTP.Paths {
		p := path.Path
		_, exist := m.paths[p]
		if exist {
			klog.Warningf("host %s path %s exist skip add. ", p, klog.KObj(ingress))
			continue
		}
		backend := backend.NewBackend(ingress.Namespace, path.Backend)
		m.paths[p] = backend

		labeler.Label(backend)
		klog.Infof("add host=%s path=%s %s", rule.Host, p, klog.KObj(ingress))
	}
}
