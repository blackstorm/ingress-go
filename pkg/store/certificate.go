package store

import (
	"crypto/tls"
	"sync"

	"github.com/blackstorm/ingress-go/pkg/k8s"
	log "github.com/blackstorm/ingress-go/pkg/logger"
	"github.com/blackstorm/ingress-go/pkg/watcher"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type cert struct {
	namespace   string
	secretName  string
	Certificate *tls.Certificate
}

type CertificateStoreKey string

// TODO is need default certificate
type CertificateStore struct {
	sync.Mutex
	certs map[CertificateStoreKey]*cert
}

func NewCertificateStore() *CertificateStore {
	return &CertificateStore{
		certs: make(map[CertificateStoreKey]*cert),
	}
}

func (c *CertificateStore) Get(key CertificateStoreKey) *cert {
	return c.certs[key]
}

func (c *CertificateStore) Update(event watcher.Event, updates ...interface{}) {
	c.Lock()
	defer c.Unlock()
	c.handleEvent(event, updates...)
}

func (c *CertificateStore) handleEvent(event watcher.Event, updates ...interface{}) {
	secret := updates[0].(*v1.Secret)
	if secret.Type == "kubernetes.io/tls" {
		namespace := secret.Namespace
		secretName := secret.Name
		log.InfoWithFields("add tls certificate to store", logrus.Fields{
			"namespace":  namespace,
			"secretName": secretName,
		})
		switch event {
		case watcher.Add:
			c.add(namespace, secretName, secret)
		case watcher.Delete:
			c.delete(namespace, secretName)
		case watcher.Update:
			c.update(namespace, secretName, secret)
		}
	}
}

func (c *CertificateStore) update(ns, name string, secret *v1.Secret) {
	if certificate, err := k8s.SecretToTLSCertificate(secret); err != nil {
		c.certs[c.key(ns, name)].Certificate = certificate
	} else {
		// TODO log
	}
}

func (c *CertificateStore) add(ns, name string, secret *v1.Secret) {
	if certificate, err := k8s.SecretToTLSCertificate(secret); err != nil {
		c.certs[c.key(ns, name)] = &cert{
			namespace:   ns,
			secretName:  name,
			Certificate: certificate,
		}
	} else {
		// TODO log
	}
}

// remove certificate just set host cert to nil
func (c *CertificateStore) delete(ns, name string) {
	key := c.key(ns, name)
	delete(c.certs, key)
}

func (c *CertificateStore) key(namespace, secretName string) CertificateStoreKey {
	return BuildCertificateStoreKey(namespace, secretName)
}

func BuildCertificateStoreKey(ns, name string) CertificateStoreKey {
	return CertificateStoreKey(ns + ":" + name)
}
