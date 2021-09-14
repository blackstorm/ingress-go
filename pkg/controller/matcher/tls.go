package matcher

import (
	"crypto/tls"

	"github.com/blackstorm/ingress-go/pkg/common"
	"github.com/blackstorm/ingress-go/pkg/k8s/api/convert"
	"github.com/blackstorm/ingress-go/pkg/store"
	"github.com/blackstorm/ingress-go/pkg/watcher"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/klog/v2"
)

type TLSMatcher struct {
	store            *store.CertificateStore
	hostTlsStoreKeys map[string]store.CertificateStoreKey
}

func NewTLSMatcher(s *store.CertificateStore) *TLSMatcher {
	return &TLSMatcher{
		store:            s,
		hostTlsStoreKeys: make(map[string]store.CertificateStoreKey),
	}
}

func (m *TLSMatcher) Match(sni string, wildcard bool) *tls.Certificate {
	if key, ok := m.hostTlsStoreKeys[sni]; ok {
		cert := m.store.Get(key)
		if cert != nil {
			return cert
		}
	}

	if wildcard {
		return nil
	}

	return m.Match(common.ToWildcardSni(sni), true)
}

func (m *TLSMatcher) Update(event watcher.Event, updates ...interface{}) {
	ingresses := convert.Ingresses(updates...)
	switch event {
	case watcher.Add:
		m.add(ingresses[0])
	case watcher.Update:
		// delete old
		m.delete(ingresses[1])
		// add new
		m.add(ingresses[0])
	case watcher.Delete:
		m.delete(ingresses[0])
	}
}

// TODO If two ingress or more has same host but diffent tls??
func (m *TLSMatcher) add(ingress *netv1.Ingress) {
	for _, tls := range ingress.Spec.TLS {

		// check tls secret is defined
		if tls.SecretName == "" {
			klog.Warningf("tls %s SecretName not defined %s", tls.Hosts, klog.KObj(ingress))
			continue
		}

		key := store.BuildCertificateStoreKey(ingress.Namespace, tls.SecretName)
		for _, host := range tls.Hosts {
			// klog.Infof("register host %s tls spec %s", host, spec.storeKey)
			m.hostTlsStoreKeys[host] = key
		}
	}
}

func (m *TLSMatcher) delete(ingress *netv1.Ingress) {
	for _, tls := range ingress.Spec.TLS {
		key := store.BuildCertificateStoreKey(ingress.Namespace, tls.SecretName)
		for _, host := range tls.Hosts {
			if k, ok := m.hostTlsStoreKeys[host]; ok {
				if k == key {
					delete(m.hostTlsStoreKeys, host)
				}
			}
		}
	}
}
