package matcher

import (
	"errors"

	"github.com/blackstorm/ingress-go/pkg/controller/backend"
	"github.com/blackstorm/ingress-go/pkg/tree"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/klog/v2"
)

type PathMatcher struct {
	tree *tree.PathTree
}

func newPathMatcher() *PathMatcher {
	return &PathMatcher{
		tree: tree.NewPathTree(),
	}
}

func (m *PathMatcher) match(path string) (*backend.Backend, error) {
	if b := m.tree.PrefixMatch(path); b != nil {
		return b.(*backend.Backend), nil
	}
	return nil, errors.New("path no found")
}

func (m *PathMatcher) add(ingress *netv1.Ingress, rule netv1.IngressRule) {
	for _, path := range rule.HTTP.Paths {
		p := path.Path
		backend := backend.NewBackend(ingress.Namespace, path.Backend)

		old := m.tree.Put(p, backend)
		if old != nil {
			klog.Warningf("Host %s path %s exist and replaced to %s", rule.Host, p, klog.KObj(ingress))
			continue
		}
		klog.Infof("add host=%s path=%s %s", rule.Host, p, klog.KObj(ingress))
	}
}

func (m *PathMatcher) update(ingress *netv1.Ingress, newRule, oldRule netv1.IngressRule) {
	m.delete(oldRule)
	m.add(ingress, newRule)
}

func (m *PathMatcher) delete(rule netv1.IngressRule) {
	for _, path := range rule.HTTP.Paths {
		m.tree.Delete(path.Path)
	}
}
