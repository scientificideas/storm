/*
Copyright Scientific Ideas 2022. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"github.com/scientificideas/storm/config"
	"github.com/scientificideas/storm/runtime"
	"github.com/scientificideas/storm/runtime/docker"
	"github.com/scientificideas/storm/runtime/k8s"
	"github.com/sirupsen/logrus"
)

func main() {
	configuration := config.GetConfig()

	var r runtime.Runtime
	switch configuration.RuntimeType {
	case "docker":
		docker, err := docker.NewDockerClient(configuration.Chaos, configuration.Filter)
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
		k8sClient, err := k8s.NewK8sClient(configuration.Chaos, configuration.Filter, configuration.Namespace, configuration.Kubeconfig, configuration.K8sContext)
		if err != nil {
			logrus.Fatal(err)
		}
		r = k8sClient
	}

	stopped := make(map[int]struct{})

	// stop containers loop
	ctx := context.Background()
	if configuration.Startfast || configuration.RuntimeType == "k8s" {
		go Loop(ctx, r, stopAndStartImmediately, stopped, configuration.Targets)
	} else {
		go Loop(ctx, r, stop, stopped, configuration.Targets)
		// start containers loop
		go Loop(ctx, r, start, stopped, configuration.Targets)
	}
	ch := make(chan struct{})
	<-ch
}
