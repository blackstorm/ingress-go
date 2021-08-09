package controller

import (
	"context"
	"net/http"

	"github.com/blackstorm/ingress-go/pkg/k8s"
	log "github.com/blackstorm/ingress-go/pkg/logger"
	"github.com/blackstorm/ingress-go/pkg/watcher"
	netv1 "k8s.io/api/networking/v1"
)

func Server(kubeConfPath string) error {
	log.Info("running ingress controller. kube config path = %s", kubeConfPath)
	client, err := k8s.GetClient(kubeConfPath)
	if err != nil {
		return err
	}

	watcherEventChannel := newWatcherEventChannel()
	serverHandler := newServerHandler()

	log.Info("start watch ingress")
	ingressWatcher := watcher.NewIngressWatcher(client)
	go ingressWatcher.Watch(context.Background(), func(event watcher.Event, updates ...interface{}) {
		watcherEventChannel.send(ChannelEvent{
			event:  event,
			values: updates,
		})
	})

	go http.ListenAndServe(":8000", serverHandler)
	go http.ListenAndServe(":8443", serverHandler)

	go watcherEventChannel.listen(func(ce ChannelEvent) {
		switch ce.event {
		case watcher.Add:
			ingress := ce.GetFirst().(*netv1.Ingress)
			serverHandler.add(ingress)
		case watcher.Delete:
			ingress := ce.GetFirst().(*netv1.Ingress)
			serverHandler.delete(ingress)
		case watcher.Update:
			old := ce.GetFirst().(*netv1.Ingress)
			update := ce.GetLast().(*netv1.Ingress)
			serverHandler.update(old, update)
		}
	})

	return nil
}
