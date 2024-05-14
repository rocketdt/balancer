// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"sync"
)

func init() {
	factories[FirstBalancer] = NewFirst
}

type First struct {
	sync.RWMutex
	hosts    []string
	priority map[string]int
}

// NewRandom create new Random balancer
func NewFirst(hosts []string) Balancer {
	priority := map[string]int{}

	for k, host := range hosts {
		priority[host] = k + 1
	}

	return &First{
		hosts:    hosts,
		priority: priority,
	}
}

// Balance selects a suitable host according
func (r *First) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.hosts) == 0 {
		return "", NoHostError
	}

	return r.hosts[0], nil
}

// Add new host to the balancer
func (b *First) Add(host string) {
	b.Lock()
	defer b.Unlock()

	for _, h := range b.hosts {
		if h == host {
			return
		}
	}

	b.hosts = append(b.hosts, host)
}

// Remove new host from the balancer
func (b *First) Remove(host string) {
	b.Lock()
	defer b.Unlock()

	for i, h := range b.hosts {
		if h == host {
			b.hosts = append(b.hosts[:i], b.hosts[i+1:]...)
			return
		}
	}
}

// Inc .
func (b *First) Inc(_ string) {}

// Done .
func (b *First) Done(_ string) {}
