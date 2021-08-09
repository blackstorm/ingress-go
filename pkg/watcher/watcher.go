package watcher

import (
	"context"

	"k8s.io/client-go/tools/cache"
)

type Watcher interface {
	Watch(context.Context, OnEvent) error
}

type baseWatcher struct {
	informer cache.SharedIndexInformer
}

func (w *baseWatcher) Watch(ctx context.Context, onEvent OnEvent) {
	w.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			onEvent(Add, obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			onEvent(Update, oldObj, newObj)
		},
		DeleteFunc: func(obj interface{}) {
			onEvent(Delete, obj)
		},
	})
	w.informer.Run((ctx.Done()))
}
