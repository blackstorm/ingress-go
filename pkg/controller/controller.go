package controller

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"

	"github.com/blackstorm/ingress-go/pkg/common"
	"github.com/blackstorm/ingress-go/pkg/k8s"
	log "github.com/blackstorm/ingress-go/pkg/logger"
	"github.com/blackstorm/ingress-go/pkg/watcher"
)

func Server(kubeConfPath string) error {
	log.Info("running ingress controller. kube config path = %s", kubeConfPath)
	client, err := k8s.GetClientWithFeedback(kubeConfPath)
	if err != nil {
		return err
	}

	// init watchers
	ingressWatcher := watcher.NewIngressWatcher(client)
	secretWatcher := watcher.NewSecretWatcher(client)

	// init matcher
	matcher := newMatcher()
	certificateStore := newCertificateStore()

	// add listeners
	ingressWatcher.AddListener(matcher)
	secretWatcher.AddListener(certificateStore)

	// end
	serverHandler := newServerHandler(matcher)

	// start watch
	log.Info("start watch sercet")
	go secretWatcher.Watch(context.Background())
	log.Info("start watch ingress")
	ingressWatcher.Watch(context.Background())

	// start servers
	go http.ListenAndServe(":8000", serverHandler)
	go listenAndServeTLS(certificateStore, 8443, serverHandler, matcher)

	// todo
	return nil
}

func listenAndServeTLS(certificateStore *certificateStore, port uint, handler *serverHandler, matcher *matcher) error {
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
