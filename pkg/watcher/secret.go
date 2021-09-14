package watcher

import (
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

type SecretWatcher struct {
	*baseWatcher
}

func NewSecretWatcher(client kubernetes.Interface) *SecretWatcher {
	factory := informers.NewSharedInformerFactory(client, time.Minute)
	informer := factory.Core().V1().Secrets().Informer()
	return &SecretWatcher{
		baseWatcher: &baseWatcher{
			name:     "secrets",
			informer: informer,
		},
	}
}
