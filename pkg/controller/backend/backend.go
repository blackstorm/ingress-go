package backend

import (
	"fmt"
	"net/url"

	netv1 "k8s.io/api/networking/v1"
)

type Backend struct {
	URL *url.URL
	// backend netv1.IngressBackend
	label string
}

func NewBackend(namespace string, backend netv1.IngressBackend) *Backend {
	host := fmt.Sprintf("%s.%s", backend.Service.Name, namespace)
	port := backend.Service.Port.Number
	serviceUrl, _ := url.Parse(fmt.Sprintf("http://%s:%d", host, port))
	return &Backend{
		URL: serviceUrl,
	}
}

func (b *Backend) AddLabel(label string) {
	b.label = label
}

func (b *Backend) GetLabel() string {
	return b.label
}
