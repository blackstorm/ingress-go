package controller

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/blackstorm/ingress-go/pkg/controller/matcher"
	"github.com/blackstorm/ingress-go/pkg/store"
	"github.com/blackstorm/ingress-go/pkg/watcher"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

func Server(client kubernetes.Interface, defaultCert DefaultCertificate) error {
	// init watchers
	ingressWatcher, secretWatcher := newWatchers(client)

	// init stores
	certificateStore := store.NewCertificateStore()

	// init matcher
	requestMatcher := matcher.NewRequestMatcher()
	tlsMatcher := matcher.NewTLSMatcher(certificateStore)

	// add listeners
	ingressWatcher.AddListeners(requestMatcher, tlsMatcher)
	secretWatcher.AddListener(certificateStore)

	// end
	serverHandler := newServerHandler(requestMatcher)

	// start watch
	runWatchers(context.Background(), ingressWatcher, secretWatcher)

	// start servers
	go http.ListenAndServe(":80", serverHandler)
	go listenAndServeTLS(443, serverHandler, tlsMatcher, defaultCert)
	go listenAndServeHttp3(serverHandler, defaultCert)

	// todo
	return nil
}

func listenAndServeTLS(port uint,
	handler *serverHandler,
	tlsMatcher *matcher.TLSMatcher, defaultCert DefaultCertificate) error {
	tlsServer := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}
	tlsServer.TLSConfig = &tls.Config{
		// CipherSuites
		// if not set CipherSuites. will use initDefaultCipherSuites for defualt sets
		GetCertificate: func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
			// klog.Infof("ssl handshake hello server %s", hello.ServerName)
			cert := tlsMatcher.Match(hello.ServerName, false)
			if cert == nil {
				klog.Warningf("server %s certificate no found", hello.ServerName)
				return nil, fmt.Errorf("no found %s certificate", hello.ServerName)
			} else {
				return cert, nil
			}
		},
	}
	return tlsServer.ListenAndServeTLS(defaultCert.certFile, defaultCert.keyFile)
}

func newWatchers(client kubernetes.Interface) (*watcher.IngressWatcher, *watcher.SecretWatcher) {
	return watcher.NewIngressWatcher(client), watcher.NewSecretWatcher(client)
}

func runWatchers(ctx context.Context, watchers ...watcher.Watcher) {
	for _, watcher := range watchers {
		go watcher.Watch(ctx)
	}
}
