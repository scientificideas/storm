package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/scientificideas/storm/chaos"
	"github.com/scientificideas/storm/container"
	"github.com/scientificideas/storm/container/docker"
	"math/rand"
	"strings"
	"time"
)

const (
	stop = iota
	start
	stopAndStartImmediately
)

type Docker struct {
	cli    *client.Client
	chaos  chaos.Chaos
	filter string
}

func NewDockerClient(chaosType, filter string) (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &Docker{cli, chaos.NewChaos(chaosType), filter}, nil
}

func (d *Docker) Chaos() chaos.Chaos {
	return d.chaos
}

func (d *Docker) GetContainers(ctx context.Context, all bool) ([]container.Container, error) {
	var filterArg filters.Args
	if d.filter != "" {
		filtersSlice := strings.Split(d.filter, ",")
		var filterKeyValues []filters.KeyValuePair
		for _, filter := range filtersSlice {
			filterKeyValues = append(filterKeyValues, filters.KeyValuePair{
				Key:   "name",
				Value: filter,
			})
		}
		filterArg = filters.NewArgs(filterKeyValues...)
	}
	containers, err := d.cli.ContainerList(ctx, types.ContainerListOptions{
		Quiet:   false,
		Size:    false,
		All:     all,
		Latest:  false,
		Since:   "",
		Before:  "",
		Limit:   0,
		Filters: filterArg,
	})
	if err != nil {
		return nil, err
	}
	rand.Shuffle(len(containers), func(i, j int) {
		containers[i], containers[j] = containers[j], containers[i]
	})

	var result []container.Container
	for _, c := range containers {
		result = append(result, docker.NewContainer(c))
	}
	return result, nil
}

func (d *Docker) StopContainer(ctx context.Context, name string) error {
	stopTime := 1 * time.Millisecond
	return d.cli.ContainerStop(ctx, name, &stopTime)
}
func (d *Docker) RmContainer(ctx context.Context, name string) error {
	return d.cli.ContainerKill(ctx, name, "SIGKILL")
}
func (d *Docker) StartContainer(ctx context.Context, name string) error {
	return d.cli.ContainerStart(ctx, name, types.ContainerStartOptions{})
}
