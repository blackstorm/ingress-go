package controller

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/blackstorm/ingress-go/pkg/controller/matcher"
	"k8s.io/klog"
)

func listenAndServeTLSHttp(port uint,
	handler http.Handler,
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
