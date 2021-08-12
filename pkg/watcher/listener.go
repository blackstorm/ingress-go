package watcher

import (
	netv1 "k8s.io/api/networking/v1"
)

type EventListener interface {
	Update(event Event, updates ...interface{})
}

type IngressEventListener interface {
	OnUpdate(event Event, updates ...*netv1.Ingress)
}

type BaseIngressEventListener struct {
	IngressEventListener
}

func (l BaseIngressEventListener) Update(event Event, updates ...interface{}) {
	ingresses := make([]*netv1.Ingress, len(updates))
	for _, update := range updates {
		ingresses = append(ingresses, update.(*netv1.Ingress))
	}
	l.OnUpdate(event, ingresses...)
}
