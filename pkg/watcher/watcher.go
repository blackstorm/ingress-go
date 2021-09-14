package watcher

import (
	"context"

	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

type Watcher interface {
	Watch(context.Context)
}

type channelEvent struct {
	event   Event
	updates []interface{}
}

type channelEventListener struct {
	ch       chan channelEvent
	listener EventListener
}

type baseWatcher struct {
	name      string
	informer  cache.SharedIndexInformer
	listeners []channelEventListener
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
	cl := channelEventListener{
		ch:       make(chan channelEvent),
		listener: listener,
	}

	if w.listeners == nil {
		w.listeners = make([]channelEventListener, 1)
		w.listeners[0] = cl
	} else {
		w.listeners = append(w.listeners, cl)
	}

	// listen update channel and send event to the channel
	listen := func(l channelEventListener) {
		for {
			event := <-l.ch
			l.listener.Update(event.event, event.updates...)
		}
	}
	go listen(cl)
}

func (w *baseWatcher) AddListeners(listeners ...EventListener) {
	for _, l := range listeners {
		w.AddListener(l)
	}
}

// notify func support publish event to channel
func (w *baseWatcher) notify(event Event, updates ...interface{}) {
	chEvent := channelEvent{
		event:   event,
		updates: updates,
	}

	publish := func(l channelEventListener, event channelEvent) {
		l.ch <- event
	}

	if w.listeners != nil {
		for _, l := range w.listeners {
			go publish(l, chEvent)
		}
	}
}
