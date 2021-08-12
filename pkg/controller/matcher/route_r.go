package matcher

/*
type routes map[string]*stream

// 为 route 添加 mapping 和 backend
func (rs routes) add(namesapce string, rule netv1.IngressRule) {
	ruleVlaue := rule.IngressRuleValue
	for _, path := range ruleVlaue.HTTP.Paths {
		backend := path.Backend
		host := fmt.Sprintf("%s.%s", backend.Service.Name, namesapce)
		port := backend.Service.Port.Number
		serviceUrl, _ := url.Parse(fmt.Sprintf("http://%s:%d", host, port))

		// TODO support loadbalance
		rs[path.Path] = &stream{
			backend:    backend,
			serviceUrl: serviceUrl,
		}

		log.InfoWithFields("add route", logrus.Fields{
			"path":       path.Path,
			"serviceUrl": serviceUrl.String(),
		})
	}
}


type router map[string]routes

func (r router) get(req *http.Request) routes {
	host := strings.Split(req.Host, ":")[0]
	if routes, ok := r[host]; ok {
		return routes
	}
	return nil
}

func (r router) add(ingress *netv1.Ingress) {
	ns := ingress.Namespace

	for _, rule := range ingress.Spec.Rules {
		host := rule.Host

		if rule.HTTP == nil {
			log.WarnWithFields("host rule is nil, skip add to router.", logrus.Fields{
				"host":      host,
				"ingress":   ingress.Name,
				"namespace": ns,
			})
			continue
		}

		var rts routes
		var ok bool
		if rts, ok = r[host]; !ok {
			rts = make(routes)
			r[host] = rts
		}

		rts.add(ns, rule)
	}
}
*/
