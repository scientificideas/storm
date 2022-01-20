/*
Copyright Scientific Ideas 2022. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package docker

import "github.com/docker/docker/api/types"

type Container struct {
	c types.Container
}

func NewContainer(container types.Container) *Container {
	return &Container{c: container}
}

func (c *Container) Name() string {
	return c.Name()
}
