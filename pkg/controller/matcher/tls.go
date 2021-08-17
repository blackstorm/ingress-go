package matcher

import (
	"crypto/tls"

	"github.com/blackstorm/ingress-go/pkg/common"
	"github.com/blackstorm/ingress-go/pkg/k8s/api/convert"
	"github.com/blackstorm/ingress-go/pkg/store"
	"github.com/blackstorm/ingress-go/pkg/watcher"
	netv1 "k8s.io/api/networking/v1"
)

type tlsSpec struct {
	storeKey store.CertificateStoreKey
}

func newTlsSpec(namespace, secretName string) *tlsSpec {
	return &tlsSpec{
		storeKey: store.BuildCertificateStoreKey(namespace, secretName),
	}
}

type TLSMatcher struct {
	store        *store.CertificateStore
	hostTlsSpecs map[string]*tlsSpec
}

func NewTLSMatcher(store *store.CertificateStore) *TLSMatcher {
	return &TLSMatcher{
		store:        store,
		hostTlsSpecs: make(map[string]*tlsSpec),
	}
}

func (m *TLSMatcher) Match(sni string, wildcard bool) *tls.Certificate {
	if spec, ok := m.hostTlsSpecs[sni]; ok {
		cert := m.store.Get(spec.storeKey)
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
	}
}

func (m *TLSMatcher) add(ingress *netv1.Ingress) {
	for _, tls := range ingress.Spec.TLS {
		spec := newTlsSpec(ingress.Namespace, tls.SecretName)
		for _, host := range tls.Hosts {
			// klog.Infof("register host %s tls spec %s", host, spec.storeKey)
			m.hostTlsSpecs[host] = spec
		}
	}
}
