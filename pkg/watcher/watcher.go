package watcher

import (
	"context"

	"k8s.io/client-go/tools/cache"
)

type Watcher interface {
	Watch(context.Context) error
}

type baseWatcher struct {
	informer  cache.SharedIndexInformer
	listeners []EventListener
}

func (w *baseWatcher) Watch(ctx context.Context) {
	w.informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			w.notify(Add, obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			w.notify(Update, newObj, oldObj)
		},
		DeleteFunc: func(obj interface{}) {
			w.notify(Delete, obj)
		},
	})
	w.informer.Run((ctx.Done()))
}

func (w *baseWatcher) AddListener(listener EventListener) {
	if w.listeners == nil {
		w.listeners = make([]EventListener, 0)
	}
	w.listeners = append(w.listeners, listener)
}

func (w *baseWatcher) notify(event Event, updates ...interface{}) {
	if w.listeners != nil {
		for _, lis := range w.listeners {
			go lis.Update(event, updates...)
		}
	}
}
