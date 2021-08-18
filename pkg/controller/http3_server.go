package controller

import (
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
)

func listenAndServeHttp3(handler http.Handler, defaultCert DefaultCertificate) {
	http3.ListenAndServeQUIC(":443", defaultCert.certFile, defaultCert.keyFile, handler)
}
