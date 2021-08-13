package watcher

import (
	"context"

	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

type Watcher interface {
	Watch(context.Context)
}

type baseWatcher struct {
	name      string
	informer  cache.SharedIndexInformer
	listeners []EventListener
}

func (w *baseWatcher) Watch(ctx context.Context) {
	klog.Infof("watching %s", w.name)
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
		w.listeners = make([]EventListener, 1)
		w.listeners[0] = listener
	} else {
		w.listeners = append(w.listeners, listener)
	}
}

func (w *baseWatcher) AddListeners(listeners ...EventListener) {
	for _, l := range listeners {
		w.AddListener(l)
	}
}

func (w *baseWatcher) notify(event Event, updates ...interface{}) {
	if w.listeners != nil {
		for _, lis := range w.listeners {
			go lis.Update(event, updates...)
		}
	}
}
