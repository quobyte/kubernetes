package main

import (
	"flag"
	"operator/pkg/controller"
	"operator/pkg/resourcehandler"
	"operator/pkg/web"
	"os"

	"github.com/golang/glog"
)

func init() {

}

func main() {
	kubeconfig := ""
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "kubeconfig file not found")
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Lookup("v").Value.Set("1")
	go web.StartWebServer()
	glog.Info("Starting Quobyte operator")
	resourcehandler.InitClient(kubeconfig)
	err := resourcehandler.CreateAllQuobyteCrd()
	if err != nil {
		glog.Errorf("Terminating operator due to: \n %v", err)
		os.Exit(1)
	}

	if err != nil {
		glog.Errorf("Terminating operator due to: \n %s", err)
		os.Exit(1)
	}
	controller.Start(resourcehandler.QclientConfig)
}
