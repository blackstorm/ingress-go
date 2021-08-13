package matcher

import (
	"net/http"
	"strings"

	convert "github.com/blackstorm/ingress-go/pkg/k8s/api/convert"
	"github.com/blackstorm/ingress-go/pkg/watcher"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/klog/v2"
)

type RequestMatcher struct {
	hostPathMatcher map[string]*PathMatcher
	hostMatcher     *HostMatcher
}

func NewRequestMatcher() *RequestMatcher {
	return &RequestMatcher{
		hostPathMatcher: make(map[string]*PathMatcher),
		hostMatcher:     newHostMatcher(),
	}
}

func (m *RequestMatcher) Match(req *http.Request) {
	host := strings.Split(req.Host, ":")[0]
	var isMatch bool
	if isMatch = m.hostMatcher.match(host); !isMatch {
		// TODO handle host is not match
	}
}

func (m *RequestMatcher) Update(event watcher.Event, values ...interface{}) {
	updates := convert.Ingresses(values...)
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
	klog.Info("add ingress ", klog.KObj(ingress))
	for _, rule := range ingress.Spec.Rules {
		h := rule.Host

		host := newHost(h, ingress)

		m.hostMatcher.add(host)

		var ok bool
		var matcher *PathMatcher

		matcher, ok = m.hostPathMatcher[h]
		if !ok {
			matcher = newPathMatcher()
			m.hostPathMatcher[h] = matcher
		}

		// add and label
		matcher.add(ingress, rule, host)
	}
}
