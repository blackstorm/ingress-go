package getter

import "k8s.io/client-go/kubernetes"

type getter struct {
	client kubernetes.Interface
}
