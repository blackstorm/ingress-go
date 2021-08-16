package matcher

import (
	"errors"

	"github.com/blackstorm/ingress-go/pkg/controller/backend"
	"github.com/blackstorm/ingress-go/pkg/controller/label"
	"github.com/blackstorm/ingress-go/pkg/tree"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/klog/v2"
)

type PathMatcher struct {
	paths *tree.PathTree
}

func newPathMatcher() *PathMatcher {
	return &PathMatcher{
		paths: tree.NewPathTree(),
	}
}

func (m *PathMatcher) match(path string) (*backend.Backend, error) {
	if be := m.paths.Match(path); be != nil {
		return be.(*backend.Backend), nil
	}
	return nil, errors.New("path no found")
}

func (m *PathMatcher) add(ingress *netv1.Ingress, rule netv1.IngressRule, labeler label.Labeler) {
	for _, path := range rule.HTTP.Paths {
		p := path.Path
		backend := backend.NewBackend(ingress.Namespace, path.Backend)

		old := m.paths.Put(p, backend)
		if old != nil {
			klog.Warningf("Host %s path %s exist and replaced to %s", rule.Host, p, klog.KObj(ingress))
			continue
		}

		labeler.Label(backend)
		klog.Infof("add host=%s path=%s %s", rule.Host, p, klog.KObj(ingress))
	}
}
