package main

import (
	"context"
	"flag"
	"github.com/scientificideas/storm/runtime"
	"github.com/scientificideas/storm/runtime/docker"
	"github.com/sirupsen/logrus"
)

func main() {
	filter := flag.String("filter", "", "filter containers by name")
	chaos := flag.String("chaos", "medium", "easy, medium or hard level of chaos")
	targets := flag.String("targets", "", `if you only want to expose certain containers, list them here ("container1,container2,container3")`)
	startfast := flag.Bool("startfast", false, `start stopped containers immediately ("true" or "false")`)
	runtimeType := flag.String("runtime", "k8s", "which orchestrator to interact with")
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
	}

	stopped := make(map[int]struct{})

	// stop containers loop
	ctx := context.Background()
	if *startfast {
		go Loop(ctx, r, stopAndStartImmediately, stopped, *targets)
	} else {
		go Loop(ctx, r, stop, stopped, *targets)
		// start containers loop
		go Loop(ctx, r, start, stopped, *targets)
	}
	ch := make(chan struct{})
	<-ch
}
