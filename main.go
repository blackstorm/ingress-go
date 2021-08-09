package main

import (
	"flag"

	"github.com/blackstorm/ingress-go/pkg/common"
	ctl "github.com/blackstorm/ingress-go/pkg/controller"
	log "github.com/blackstorm/ingress-go/pkg/logger"
)

func main() {
	log.Info("fast ingress controller.")

	kubeConfPath := flag.String("kubeconfigPath", "", "kube config path")
	flag.Parse()

	signal := make(chan bool)

	err := ctl.Server(common.CheckOrDefault(kubeConfPath, common.EMPTY_STRING))
	if err != nil {
		panic(err)
	}

	<-signal
}
