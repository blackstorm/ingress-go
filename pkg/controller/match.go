package controller

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	log "github.com/blackstorm/ingress-go/pkg/logger"
	"github.com/blackstorm/ingress-go/pkg/watcher"
	"github.com/sirupsen/logrus"
	netv1 "k8s.io/api/networking/v1"
)

type route struct {
	backend    netv1.IngressBackend
	serviceUrl *url.URL
}

func (r route) handleRequest(resp http.ResponseWriter, req *http.Request) {
	proxy := httputil.NewSingleHostReverseProxy(r.serviceUrl)
	proxy.ServeHTTP(resp, req)
}

type routes map[string]*route

// 为 route 添加 mapping 和 backend
func (rs routes) add(namesapce string, rule netv1.IngressRule) {
	ruleVlaue := rule.IngressRuleValue
	for _, path := range ruleVlaue.HTTP.Paths {
		backend := path.Backend
		host := fmt.Sprintf("%s.%s", backend.Service.Name, namesapce)
		port := backend.Service.Port.Number
		serviceUrl, _ := url.Parse(fmt.Sprintf("http://%s:%d", host, port))

		// TODO support loadbalance
		rs[path.Path] = &route{
			backend:    backend,
			serviceUrl: serviceUrl,
		}

		log.InfoWithFields("add route", logrus.Fields{
			"path":       path.Path,
			"serviceUrl": serviceUrl.String(),
		})
	}
}

type hostsRoutes map[string]routes

// 添加一个 host route
func (h hostsRoutes) add(ingress *netv1.Ingress) {
	for _, rule := range ingress.Spec.Rules {
		// ----------
		// contains host but not has http bloks
		// ----------
		// spec:
		// 	rules:
		// 	- host: "foo.bar.com"
		// 	- host: "*.foo.com"
		// 		http:
		// 			paths:
		// 			- pathType: Prefix
		// 				path: "/foo"
		// 				backend:
		// 					service:
		// 						name: service2
		// 						port:
		// 							number: 80
		if rule.HTTP == nil {
			log.WarnWithFields("host rule is nil, skip add route.", logrus.Fields{
				"host":      rule.Host,
				"ingress":   ingress.Name,
				"namespace": ingress.Namespace,
			})
			continue
		}

		var rts routes
		var exist bool

		if rts, exist = h[rule.Host]; !exist {
			rts = make(routes)
			h[rule.Host] = rts
		}

		rts.add(ingress.Namespace, rule)
	}
}

type matcher struct {
	sniCertStoreKeys map[string]certificateStoreKey
	hostsRoutes      hostsRoutes
}

func newMatcher() *matcher {
	return &matcher{
		sniCertStoreKeys: make(map[string]certificateStoreKey),
	}
}

func (m *matcher) Update(event watcher.Event, updates ...interface{}) {
	ingress := updates[0].(*netv1.Ingress)
	switch event {
	case watcher.Add:
		m.add(ingress)
	}
}

func (m *matcher) add(ingress *netv1.Ingress) {
	m.updateSniKeys(ingress)
	m.hostsRoutes.add(ingress)
}

func (m *matcher) updateSniKeys(ingress *netv1.Ingress) {
	ns := ingress.Namespace
	tlss := ingress.Spec.TLS
	for _, tls := range tlss {
		for _, host := range tls.Hosts {
			m.sniCertStoreKeys[host] = buildCertificateStoreKey(ns, tls.SecretName)
		}
	}
}
