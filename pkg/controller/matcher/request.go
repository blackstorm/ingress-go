package matcher

import (
	"net/http"
	"strings"

	"github.com/blackstorm/ingress-go/pkg/watcher"
	netv1 "k8s.io/api/networking/v1"
)

type RequestMatcher struct {
	routerMatcher *RouterMatcher
	hostMatcher   *HostMatcher
	watcher.BaseIngressEventListener
}

func NewRequestMatcher() *RequestMatcher {
	return &RequestMatcher{
		hostMatcher:   NewHostMatcher(),
		routerMatcher: NewRouterMatcher(),
	}
}

func (m *RequestMatcher) Match(req *http.Request) {
	host := strings.Split(req.Host, ":")[0]
	var isMatch bool
	if isMatch = m.hostMatcher.match(host); !isMatch {
		// TODO handle host is not match
	}
}

func (m *RequestMatcher) OnUpdate(event watcher.Event, updates ...*netv1.Ingress) {
	// TODO lock
	switch event {
	case watcher.Add:
		m.add(updates[0])
	case watcher.Update:
		// TODO
	case watcher.Delete:
		// TODO
	}
}

// 使用 label 删除 path 和 route
func (m *RequestMatcher) add(ingress *netv1.Ingress) {
	name := ingress.Name
	namespace := ingress.Namespace
	for _, rule := range ingress.Spec.Rules {
		host := NewHost(rule.Host, namespace, name)

		m.hostMatcher.add(host)
		m.routerMatcher.add(rule, host)
	}
}
