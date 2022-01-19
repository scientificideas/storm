package runtime

import (
	"context"
	"github.com/scientificideas/storm/chaos"
	"github.com/scientificideas/storm/container"
)

// Runtime is a container runtime
type Runtime interface {
	GetContainers(ctx context.Context, all bool) ([]container.Container, error)
	StopContainer(ctx context.Context, name string) error
	RmContainer(ctx context.Context, name string) error
	StartContainer(ctx context.Context, name string) error
	Chaos() chaos.Chaos
}
