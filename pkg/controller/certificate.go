package controller

import (
	"crypto/tls"
	"sync"

	"github.com/blackstorm/ingress-go/pkg/common"
	log "github.com/blackstorm/ingress-go/pkg/logger"
	"github.com/blackstorm/ingress-go/pkg/watcher"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

type cert struct {
	namespace   string
	secretName  string
	certificate *tls.Certificate
}

type certificateStoreKey string

// TODO is need default certificate
type certificateStore struct {
	sync.Mutex
	certs map[certificateStoreKey]*cert
}

func newCertificateStore() *certificateStore {
	return &certificateStore{
		certs: make(map[certificateStoreKey]*cert),
	}
}

func (c *certificateStore) get(key certificateStoreKey) *cert {
	return c.certs[key]
}

func (c *certificateStore) Update(event watcher.Event, updates ...interface{}) {
	c.Lock()
	defer c.Unlock()
	c.handleEvent(event, updates...)
}

func (c *certificateStore) handleEvent(event watcher.Event, updates ...interface{}) {
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

func (c *certificateStore) update(ns, name string, secret *v1.Secret) {
	if certificate, err := common.SecretToTLSCertificate(secret); err != nil {
		c.certs[c.key(ns, name)].certificate = certificate
	} else {
		// TODO log
	}
}

func (c *certificateStore) add(ns, name string, secret *v1.Secret) {
	if certificate, err := common.SecretToTLSCertificate(secret); err != nil {
		c.certs[c.key(ns, name)] = &cert{
			namespace:   ns,
			secretName:  name,
			certificate: certificate,
		}
	} else {
		// TODO log
	}
}

// remove certificate just set host cert to nil
func (c *certificateStore) delete(ns, name string) {
	key := c.key(ns, name)
	if cert := c.certs[key]; cert != nil {
		c.certs[key] = nil
	}
}

func (c *certificateStore) key(namespace, secretName string) certificateStoreKey {
	return buildCertificateStoreKey(namespace, secretName)
}

func buildCertificateStoreKey(ns, name string) certificateStoreKey {
	return certificateStoreKey(ns + ":" + name)
}
