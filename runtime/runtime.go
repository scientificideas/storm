/*
Copyright Scientific Ideas 2022. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package runtime

import (
	"context"
	"github.com/scientificideas/storm/chaos"
	"github.com/scientificideas/storm/container"
)

// Runtime is an interface for interaction with a container runtime
type Runtime interface {
	GetContainers(ctx context.Context, all bool) ([]container.Container, error)
	StopContainer(ctx context.Context, name string) error
	RmContainer(ctx context.Context, name string) error
	StartContainer(ctx context.Context, name string) error
	Chaos() chaos.Chaos
	Type() string
}
