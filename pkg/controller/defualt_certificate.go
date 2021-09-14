package controller

type DefaultCertificate struct {
	certFile string
	keyFile  string
}

func NewDefaultCertificate(certFile *string, keyFile *string) DefaultCertificate {
	return DefaultCertificate{
		certFile: *certFile,
		keyFile:  *keyFile,
	}
}
