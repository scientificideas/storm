package main

import (
	"context"
	"flag"
	"github.com/scientificideas/storm/runtime"
	"github.com/scientificideas/storm/runtime/docker"
	"github.com/scientificideas/storm/runtime/k8s"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	filter := flag.String("filter", "", "filter containers by name")
	chaos := flag.String("chaos", "medium", "easy, medium or hard level of chaos")
	targets := flag.String("targets", "", `if you only want to expose certain containers, list them here ("container1,container2,container3")`)
	startfast := flag.Bool("startfast", false, `start stopped containers immediately ("true" or "false")`)
	runtimeType := flag.String("runtime", "k8s", "which orchestrator to interact with")
	namespace := flag.String("kube-namespace", "", "k8s namespace")
	k8sContext := flag.String("kube-context", "", "k8s context")

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	flag.Parse()

	var r runtime.Runtime
	switch *runtimeType {
	case "docker":
		docker, err := docker.NewDockerClient(*chaos, *filter)
		if err != nil {
			logrus.Fatal(err)
		}
		ctx := context.Background()
		_, err = docker.Cli.Ping(ctx)
		if err != nil {
			logrus.Fatal(err)
		}
		r = docker
	case "k8s":
		k8sClient, err := k8s.NewK8sClient(*chaos, *filter, *namespace, *kubeconfig, *k8sContext)
		if err != nil {
			logrus.Fatal(err)
		}
		r = k8sClient
	}

	stopped := make(map[int]struct{})

	// stop containers loop
	ctx := context.Background()
	if *startfast || *runtimeType == "k8s" {
		go Loop(ctx, r, stopAndStartImmediately, stopped, *targets)
	} else {
		go Loop(ctx, r, stop, stopped, *targets)
		// start containers loop
		go Loop(ctx, r, start, stopped, *targets)
	}
	ch := make(chan struct{})
	<-ch
}
