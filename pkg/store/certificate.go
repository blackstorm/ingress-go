package store

import (
	"crypto/tls"
	"fmt"

	"github.com/blackstorm/ingress-go/pkg/k8s"
	"github.com/blackstorm/ingress-go/pkg/watcher"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

type CertificateStoreKey string

// TODO is need default certificate
type CertificateStore struct {
	certs map[CertificateStoreKey]*tls.Certificate
}

func NewCertificateStore() *CertificateStore {
	return &CertificateStore{
		certs: make(map[CertificateStoreKey]*tls.Certificate),
	}
}

func (c *CertificateStore) Get(key CertificateStoreKey) *tls.Certificate {
	fmt.Printf("\nget ceritficate %s\n", key)
	return c.certs[key]
}

func (c *CertificateStore) Update(event watcher.Event, updates ...interface{}) {
	c.handleEvent(event, updates...)
}

func (c *CertificateStore) handleEvent(event watcher.Event, updates ...interface{}) {
	secret := updates[0].(*v1.Secret)
	if secret.Type == "kubernetes.io/tls" {
		klog.Info("handle tls certificate ", klog.KObj(secret))
		switch event {
		case watcher.Add:
			c.add(secret)
		case watcher.Delete:
			c.delete(secret)
		case watcher.Update:
			c.update(secret)
		}
	}
}

func (c *CertificateStore) update(secret *v1.Secret) {
	if certificate, err := k8s.SecretToTLSCertificate(secret); err != nil {
		c.certs[c.key(secret.Namespace, secret.Name)] = certificate
	} else {
		// TODO log
	}
}

func (c *CertificateStore) add(secret *v1.Secret) {
	if certificate, err := k8s.SecretToTLSCertificate(secret); err != nil {
		klog.Error("secret to tls certificate error. pelease check the secret: ", klog.KObj(secret))
	} else {
		c.certs[c.key(secret.Namespace, secret.Name)] = certificate
	}
}

// remove certificate just set host cert to nil
func (c *CertificateStore) delete(secret *v1.Secret) {
	key := c.key(secret.Namespace, secret.Name)
	delete(c.certs, key)
}

func (c *CertificateStore) key(namespace, secretName string) CertificateStoreKey {
	return BuildCertificateStoreKey(namespace, secretName)
}

func BuildCertificateStoreKey(namespace, secretName string) CertificateStoreKey {
	return CertificateStoreKey(namespace + ":" + secretName)
}
