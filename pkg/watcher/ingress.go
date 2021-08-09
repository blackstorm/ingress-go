package watcher

import (
	"time"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
)

type IngressWatcher struct {
	*baseWatcher
}

func NewIngressWatcher(client kubernetes.Interface) *IngressWatcher {
	factory := informers.NewSharedInformerFactory(client, time.Minute)
	informer := factory.Networking().V1().Ingresses().Informer()
	return &IngressWatcher{
		baseWatcher: &baseWatcher{
			informer: informer,
		},
	}
}
