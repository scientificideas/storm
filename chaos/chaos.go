/*
Copyright Scientific Ideas 2022. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package chaos

import (
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

const (
	// levels
	Easy   = "easy"
	Medium = "medium"
	Hard   = "hard"

	// timeout ranges for levels
	minTimeoutMillisecsEasy   = 30000
	maxTimeoutMillisecsEasy   = 60000
	minTimeoutMillisecsMedium = 8000
	maxTimeoutMillisecsMedium = 20000
	minTimeoutMillisecsHard   = 500
	maxTimeoutMillisecsHard   = 3000
)

type loopType int

type Chaos interface {
	Timeout() time.Duration
}

func NewChaos(lvl string) Chaos {
	switch lvl {
	case Easy:
		return &ChaosDefault{
			min: minTimeoutMillisecsEasy,
			max: maxTimeoutMillisecsEasy,
		}
	case Medium:
		return &ChaosDefault{
			min: minTimeoutMillisecsMedium,
			max: maxTimeoutMillisecsMedium,
		}
	case Hard:
		return &ChaosDefault{
			min: minTimeoutMillisecsHard,
			max: maxTimeoutMillisecsHard,
		}
	default:
		logrus.Warnf(`chaos level %s is not supported, default level ("medium") will be used`, lvl)
		return &ChaosDefault{
			min: minTimeoutMillisecsMedium,
			max: maxTimeoutMillisecsMedium,
		}
	}
	return nil
}

type ChaosDefault struct {
	min int
	max int
}

func (c *ChaosDefault) Timeout() time.Duration {
	return time.Duration(rand.Intn(c.max-c.min)+c.min) * time.Millisecond
}
