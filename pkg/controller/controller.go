package controller

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/blackstorm/ingress-go/pkg/common"
	"github.com/blackstorm/ingress-go/pkg/controller/matcher"
	"github.com/blackstorm/ingress-go/pkg/store"
	"github.com/blackstorm/ingress-go/pkg/watcher"
	"k8s.io/client-go/kubernetes"
)

func Server(client kubernetes.Interface) error {
	// init watchers
	ingressWatcher, secretWatcher := newWatchers(client)

	// init stores
	certificateStore := store.NewCertificateStore()

	// init matcher
	matcher := matcher.NewRequestMatcher()

	// add listeners
	ingressWatcher.AddListener(matcher)
	secretWatcher.AddListener(certificateStore)

	// end
	serverHandler := newServerHandler(matcher)

	// start watch
	runWatchers(context.Background(), ingressWatcher, secretWatcher)

	// start servers
	go http.ListenAndServe(":8000", serverHandler)
	go listenAndServeTLS(8443, serverHandler)

	// todo
	return nil
}

func listenAndServeTLS(port uint, handler *serverHandler) error {
	tlsServer := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	tlsServer.TLSConfig = &tls.Config{
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			return nil, errors.New("todo")
		},
	}
	return tlsServer.ListenAndServeTLS(common.EMPTY_STRING, common.EMPTY_STRING)
}

func newWatchers(client kubernetes.Interface) (*watcher.IngressWatcher, *watcher.SecretWatcher) {
	return watcher.NewIngressWatcher(client), watcher.NewSecretWatcher(client)
}

func runWatchers(ctx context.Context, watchers ...watcher.Watcher) {
	for _, watcher := range watchers {
		go watcher.Watch(ctx)
	}
}
