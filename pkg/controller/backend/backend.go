package backend

import "net/url"

type Backend struct {
	URL   *url.URL
	label string
}

func NewBackend() *Backend {
	return &Backend{}
}

func (b *Backend) Label(label string) {
	b.label = label
}

func (b *Backend) GetLabel() string {
	return b.label
}
