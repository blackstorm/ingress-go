package matcher

import (
	"fmt"
	"strings"

	"github.com/blackstorm/ingress-go/pkg/controller/label"
)

type Host struct {
	host      string
	namespace string
	ingress   string
	label     string
}

func NewHost(host, namespace, ingress string) *Host {
	return &Host{
		host:      host,
		namespace: namespace,
		ingress:   ingress,
		label:     fmt.Sprintf("%s:%s:%s", host, namespace, ingress),
	}
}

func (h *Host) LabelResource(res label.Resource) {
	res.Label(h.label)
}

func (h *Host) IsLabeled(res label.Resource) bool {
	if res != nil {
		return strings.Compare(h.label, res.GetLabel()) == 0
	}
	return false
}

func (h *Host) equals(host *Host) bool {
	return strings.Compare(h.label, host.label) == 0
}

type HostMatcher struct {
	hosts map[string][]*Host
}

func NewHostMatcher() *HostMatcher {
	return &HostMatcher{
		hosts: make(map[string][]*Host),
	}
}

func (m *HostMatcher) match(host string) bool {
	_, ok := m.hosts[host]
	return ok
}

func (m *HostMatcher) add(host *Host) {
	var hosts []*Host
	if hosts = m.hosts[host.host]; hosts == nil {
		hosts = make([]*Host, 1)
	}
	hosts = append(hosts, host)
	m.hosts[host.host] = hosts
}

func (m *HostMatcher) delete(host *Host) {
	if hosts, ok := m.hosts[host.host]; ok {
		for i, it := range hosts {
			if it.equals(host) {
				// remove it
				m.hosts[host.host] = append(hosts[:i], hosts[i+1:]...)
				return
			}
		}
	}
}
