package watcher

import (
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

type SecretWatcher struct {
	*baseWatcher
}

func NewSecretWatcher(client kubernetes.Interface) *IngressWatcher {
	factory := informers.NewSharedInformerFactory(client, time.Minute)
	informer := factory.Core().V1().Secrets().Informer()
	return &IngressWatcher{
		baseWatcher: &baseWatcher{
			informer: informer,
		},
	}
}
