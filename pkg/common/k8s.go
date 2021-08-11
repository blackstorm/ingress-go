package common

import (
	"crypto/tls"

	v1 "k8s.io/api/core/v1"
)

func SecretToTLSCertificate(secret *v1.Secret) (*tls.Certificate, error) {
	cert, err := tls.X509KeyPair(secret.Data["tls.crt"], secret.Data["tls.key"])
	if err != nil {
		return nil, err
	}
	return &cert, nil
}
