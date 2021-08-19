package controller

import (
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
)

func listenAndServeHttp3(shr *serverHandler, defaultCert DefaultCertificate) {
	server := &http3.Server{
		Server: &http.Server{
			Addr:    ":443",
			Handler: shr,
		},
	}

	shr.addAfterHandleRequest(func(resp http.ResponseWriter, req *http.Request) {
		server.SetQuicHeaders(resp.Header())
	})

	server.ListenAndServeTLS(defaultCert.certFile, defaultCert.keyFile)
}
