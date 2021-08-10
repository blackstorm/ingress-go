package controller

import (
	"crypto/tls"
	"reflect"

	"github.com/blackstorm/ingress-go/pkg/common"
	"github.com/blackstorm/ingress-go/pkg/getter"
	log "github.com/blackstorm/ingress-go/pkg/logger"
	"github.com/blackstorm/ingress-go/pkg/watcher"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
)

type cert struct {
	// TOOD
	secretName  string
	certificate *tls.Certificate
}

// TODO is need default certificate
type certificateStore struct {
	secretGetter *getter.SecretGetter
	certs        map[string]*cert
}

func newCertificateStore(
	ingressWatcher *watcher.IngressWatcher,
	secretWatcher *watcher.SecretWatcher,
	secretGetter *getter.SecretGetter) *certificateStore {
	store := &certificateStore{
		certs:        make(map[string]*cert),
		secretGetter: secretGetter,
	}
	ingressWatcher.AddListener(store)
	secretWatcher.AddListener(store)
	return store
}

func (c *certificateStore) Get(sni string) *cert {
	if cert, ok := c.certs[sni]; ok {
		return cert
	} else {
		return c.certs[common.ToWildcardSni(sni)]
	}
}

func (c *certificateStore) Update(event watcher.Event, updates ...interface{}) {
	update := updates[0]
	t := reflect.TypeOf(update)
	switch t {
	case reflect.TypeOf((*netv1.Ingress)(nil)):
		c.handleIngressEvent(event, updates...)
	case reflect.TypeOf((*v1.Secret)(nil)):
		c.handleSecretEvent(event, updates...)
	}
}

func (c *certificateStore) handleIngressEvent(event watcher.Event, updates ...interface{}) {
	switch event {
	case watcher.Add:
		ingress := updates[0].(*netv1.Ingress)
		namesapce := ingress.Namespace
		for _, tls := range ingress.Spec.TLS {
			secretName := tls.SecretName

			// init cert
			cert := &cert{
				secretName: tls.SecretName,
			}

			// get certificate
			if certificate, err := c.getTLSCertificate(namesapce, secretName); err == nil {
				cert.certificate = certificate
			} else {
				log.WarnWithFields("get tls certificate failed", logrus.Fields{
					"namespace":  namesapce,
					"secretName": secretName,
					"error":      err,
				})
			}

			// add tls for hosts
			for _, h := range tls.Hosts {
				c.certs[h] = cert
			}
		}
	}
}

func (c *certificateStore) handleSecretEvent(event watcher.Event, updates ...interface{}) {
	// TODO
}

func (c *certificateStore) getTLSCertificate(namespace string, secretName string) (*tls.Certificate, error) {
	// var err error
	if secret, err := c.secretGetter.Get(namespace, secretName); err == nil {
		if cert, err := tls.X509KeyPair(secret.Data["tls.crt"], secret.Data["tls.key"]); err == nil {
			return &cert, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
