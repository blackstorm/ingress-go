package controller

import "k8s.io/client-go/kubernetes"

type certificateStore struct {
	certs  map[string]string
	client kubernetes.Interface
}

func newCertificateStore() *certificateStore {
	// TODO
	return nil
}
