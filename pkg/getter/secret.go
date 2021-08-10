package getter

import (
	"time"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/core/v1"
)

type SecretGetter struct {
	lister v1.SecretLister
}

func NewSecretGetter(client kubernetes.Interface) *SecretGetter {
	factory := informers.NewSharedInformerFactory(client, time.Minute)
	return &SecretGetter{
		lister: factory.Core().V1().Secrets().Lister(),
	}
}

func (g *SecretGetter) Get(namespace string, name string) (*apiv1.Secret, error) {
	return g.lister.Secrets(namespace).Get(name)
}
