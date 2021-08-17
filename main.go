package main

import (
	"flag"

	banner "github.com/blackstorm/ingress-go/pkg/banner"
	ctl "github.com/blackstorm/ingress-go/pkg/controller"
	"github.com/blackstorm/ingress-go/pkg/k8s"
	"k8s.io/klog/v2"
)

func main() {
	defaultCertFile := flag.String("defaultCertFile", "/etc/ingress-go/default.crt", "the server default certificate")
	defaultKeyFile := flag.String("defaultKeyFile", "/etc/ingress-go/default.key", "the server default certificate key")
	isHiddenBanner := flag.Bool("hiddenBanner", false, "dont print banner")
	kubeConfPath := flag.String("kubeconfigPath", "", "kube config path")
	flag.Parse()

	banner.Print(isHiddenBanner)

	// get kubernetes client
	client, err := k8s.GetClient(kubeConfPath)
	if err != nil {
		panic(err)
	}

	klog.Info("ingress-go")

	// run server
	signal := make(chan bool)
	err = ctl.Server(client, ctl.NewDefaultCertificate(defaultCertFile, defaultKeyFile))
	if err != nil {
		panic(err)
	}
	<-signal
}
