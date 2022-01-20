/*
Copyright Scientific Ideas 2022. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

package k8s

import "k8s.io/api/core/v1"

type Pod struct {
	pod v1.Pod
}

func NewPod(pod v1.Pod) *Pod {
	return &Pod{pod: pod}
}

func (p *Pod) Name() string {
	return p.pod.Name
}
