package getter

import (
	"crypto/tls"
	"time"

	"github.com/blackstorm/ingress-go/pkg/common"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/core/v1"
)

type SecretGetter struct {
	lister v1.SecretLister
}

func NewSecretGetter(client kubernetes.Interface) *SecretGetter {
	factory := informers.NewSharedInformerFactory(client, time.Minute)
	factory.Core().V1().Secrets().Lister().List(labels.Everything())
	return &SecretGetter{
		lister: factory.Core().V1().Secrets().Lister(),
	}
}

func (g *SecretGetter) Get(namespace, name string) (*apiv1.Secret, error) {
	return g.lister.Secrets(namespace).Get(name)
}

func (g *SecretGetter) GetTLSCertificate(namespace, secretName string) (*tls.Certificate, error) {
	// var err error
	if secret, err := g.Get(namespace, secretName); err == nil {
		if cert, err := common.SecretToTLSCertificate(secret); err == nil {
			return cert, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
