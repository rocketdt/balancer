// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"sort"
	"sync"
)

// -----------------------------------

type HostPriority struct {
	host     string
	priority int
}

type ByHostPriority []HostPriority

func (a ByHostPriority) Len() int           { return len(a) }
func (a ByHostPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHostPriority) Less(i, j int) bool { return a[i].priority < a[j].priority }

func sortHosts(hostL []string, p map[string]int) []string {
	hpL := make([]HostPriority, len(hostL))
	ret := make([]string, len(hostL))

	for k, h := range hostL {
		hpL[k] = HostPriority{
			host:     h,
			priority: p[h],
		}
	}

	sort.Sort(ByHostPriority(hpL))

	for k, hp := range hpL {
		ret[k] = hp.host
	}

	return ret
}

// -----------------------------------

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

	if len(b.hosts) < 2 {
		b.hosts = append(b.hosts, host)
		return
	}

	b.hosts = append(b.hosts, host)
	hL := sortHosts(b.hosts[1:], b.priority)
	b.hosts = append([]string{b.hosts[0]}, hL...)
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
