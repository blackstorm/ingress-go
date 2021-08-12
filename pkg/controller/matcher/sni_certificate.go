package matcher

import (
	"crypto/tls"
	"sync"

	"github.com/blackstorm/ingress-go/pkg/store"
	"github.com/blackstorm/ingress-go/pkg/watcher"
	netv1 "k8s.io/api/networking/v1"
)

type SniCertificateMatcher struct {
	certStore *store.CertificateStore
	keys      map[string]store.CertificateStoreKey
	sync.Mutex
}

func NewSniCertificateMatcher(certStore *store.CertificateStore) *SniCertificateMatcher {
	return &SniCertificateMatcher{
		certStore: certStore,
		keys:      make(map[string]store.CertificateStoreKey),
	}
}

func (m *SniCertificateMatcher) Get(sni string) *tls.Certificate {
	if key, ok := m.keys[sni]; ok {
		if cert := m.certStore.Get(key); cert != nil {
			return cert.Certificate
		}
	}
	return nil
}

func (m *SniCertificateMatcher) Update(event watcher.Event, updates ...interface{}) {
	m.Lock()
	defer m.Unlock()
	ingress := updates[0].(*netv1.Ingress)
	switch event {
	case watcher.Add:
		m.add(ingress)
	}
}

func (m *SniCertificateMatcher) add(ingress *netv1.Ingress) {
	ns := ingress.Namespace
	tlss := ingress.Spec.TLS
	for _, tls := range tlss {
		for _, host := range tls.Hosts {
			m.keys[host] = store.BuildCertificateStoreKey(ns, tls.SecretName)
		}
	}
}
