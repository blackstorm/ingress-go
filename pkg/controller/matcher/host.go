package matcher

import (
	"fmt"

	"github.com/blackstorm/ingress-go/pkg/controller/label"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/klog/v2"
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

type HostMatcher map[string]*linkedHost

func newHostMatcher() HostMatcher {
	return make(HostMatcher)
}

func (m HostMatcher) match(host string) bool {
	linked, ok := m[host]
	return ok && linked != nil
}

func (m HostMatcher) add(host Host) {
	if linked, ok := m[host.host]; ok {
		if linked.find(host) != nil {
			klog.Infof("host %s existed skip add", host.host)
		} else {
			klog.Infof("add host %s to matcher", host.host)
			linked.append(&host)
		}
	} else {
		klog.Infof("add host %s to matcher", host.host)
		m[host.host] = newLinkedHost(&host)
	}
}

func (m HostMatcher) delete(host Host) {
	if linked, ok := m[host.host]; ok {
		if res := linked.find(host); res != nil {
			if res.remove() == -1 {
				m[host.host] = nil
			}
		}
	}
}
