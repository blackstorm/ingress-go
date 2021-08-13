package k8s

import (
	netv1 "k8s.io/api/networking/v1"
)

func Ingresses(values ...interface{}) []*netv1.Ingress {
	if values != nil {
		len := len(values)
		if len > 0 {
			ingresses := make([]*netv1.Ingress, len)
			for i, value := range values {
				ingresses[i] = ingress(value)
			}
			return ingresses
		}
		return make([]*netv1.Ingress, 0)
	}
	return nil
}

func Ingress(value interface{}) *netv1.Ingress {
	if value != nil {
		return ingress(value)
	}
	return nil
}

func ingress(value interface{}) *netv1.Ingress {
	return value.(*netv1.Ingress)
}
