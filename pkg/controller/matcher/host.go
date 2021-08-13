package matcher

import (
	"fmt"

	"github.com/blackstorm/ingress-go/pkg/controller/label"
	netv1 "k8s.io/api/networking/v1"
)

type Host struct {
	host string
	label.BaseLabeler
}

func newHost(host string, ingress *netv1.Ingress) Host {
	return Host{
		host: host,
		BaseLabeler: label.BaseLabeler{
			LabelName: fmt.Sprintf("%s:%s:%s", host, ingress.Namespace, ingress.Name),
		},
	}
}

type HostMatcher struct {
	hosts map[string][]Host
}

func newHostMatcher() *HostMatcher {
	return &HostMatcher{
		hosts: make(map[string][]Host),
	}
}

func (m *HostMatcher) match(host string) bool {
	_, ok := m.hosts[host]
	return ok
}

func (m *HostMatcher) add(host Host) {
	var hs []Host
	if hs = m.hosts[host.host]; hs == nil {
		hs = make([]Host, 1)
	}
	hs = append(hs, host)
	m.hosts[host.host] = hs
}
