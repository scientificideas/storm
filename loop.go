/*
Copyright Scientific Ideas 2022. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"github.com/scientificideas/storm/runtime"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

const (
	stop = iota
	start
	stopAndStartImmediately
	undefined
)

type LoopType int

func Loop(ctx context.Context, runtime runtime.Runtime, loop LoopType, stopped map[int]struct{}, targets string) {
	if runtime.Type() == "k8s" {
		loop = undefined
	}

	var (
		msg              string
		actionFn         func(ctx context.Context, name string) error
		targetContainers []string
	)

	if targets != "" {
		targetContainers = strings.Split(targets, ",")
	}

	switch {
	case loop == stop || loop == stopAndStartImmediately:
		msg = "stop"
		actionFn = runtime.StopContainer
	case loop == start:
		msg = "start"
		actionFn = runtime.StartContainer
	default:
		msg = "delete pod"
		actionFn = runtime.RmContainer
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
		case loop == stop || loop == stopAndStartImmediately || loop == undefined:
			if len(containers) == 0 {
				continue
			}
			if len(targetContainers) != 0 {
				targetIndex = rand.Intn(len(targetContainers))
			} else {
				targetIndex = rand.Intn(len(containers))
			}
		case loop == start:
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
		case stop:
			stopped[targetIndex] = struct{}{}
		case start:
			delete(stopped, targetIndex)
		case stopAndStartImmediately:
			logrus.Infof("start %s", selectedContainer)
			if err = runtime.StartContainer(ctx, selectedContainer); err != nil {
				logrus.Error(err)
			}
		default:

		}
	}
}
