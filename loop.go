package main

import (
	"context"
	"github.com/scientificideas/storm/runtime"
	"github.com/scientificideas/storm/runtime/docker"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

func Loop(ctx context.Context, runtime runtime.Runtime, loop docker.loopType, stopped map[int]struct{}, targets string) {
	var (
		msg              string
		actionFn         func(ctx context.Context, name string) error
		targetContainers []string
	)

	if targets != "" {
		targetContainers = strings.Split(targets, ",")
	}

	switch {
	case loop == docker.stop || loop == docker.stopAndStartImmediately:
		msg = "stop"
		actionFn = runtime.StopContainer
	case loop == docker.start:
		msg = "start"
		actionFn = runtime.StartContainer
	}

	var targetIndex int
	for {
		rand.Seed(time.Now().UnixNano())
		time.Sleep(runtime.Chaos().Timeout())

		containers, err := runtime.GetContainers(ctx, false)
		if err != nil {
			logrus.Fatal(err)
		}

		switch {
		case loop == docker.stop || loop == docker.stopAndStartImmediately:
			if len(containers) == 0 {
				continue
			}
			if len(targetContainers) != 0 {
				targetIndex = rand.Intn(len(targetContainers))
			} else {
				targetIndex = rand.Intn(len(containers))
			}
		case loop == docker.start:
			containers, err = runtime.GetContainers(ctx, true)
			if err != nil {
				logrus.Fatal(err)
			}
			if len(stopped) == 0 {
				if len(targetContainers) != 0 {
					targetIndex = rand.Intn(len(targetContainers))
				} else {
					targetIndex = rand.Intn(len(containers))
				}
			} else {
				targetIndex = rand.Intn(len(stopped))
				if _, ok := stopped[targetIndex]; !ok {
					continue
				}
			}
		}

		var selectedContainer string
		if len(targetContainers) != 0 {
			selectedContainer = targetContainers[targetIndex]
		} else {
			selectedContainer = containers[targetIndex].Name()
		}

		logrus.Infof("%s %s", msg, selectedContainer)
		if err = actionFn(ctx, selectedContainer); err != nil {
			logrus.Error(err)
		}
		switch loop {
		case docker.stop:
			stopped[targetIndex] = struct{}{}
		case docker.start:
			delete(stopped, targetIndex)
		case docker.stopAndStartImmediately:
			logrus.Infof("start %s", selectedContainer)
			if err = runtime.StartContainer(ctx, selectedContainer); err != nil {
				logrus.Error(err)
			}
		}
	}
}
