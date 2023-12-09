// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

import (
	"sync"
)

//RoundRobin will select the server in turn from the server to proxy
type RoundRobin struct {
	BaseBalancer
	i   uint64
	sem sync.Mutex
}

func init() {
	factories[R2Balancer] = NewRoundRobin
}

// NewRoundRobin create new RoundRobin balancer
func NewRoundRobin(hosts []string) Balancer {
	return &RoundRobin{
		i:   0,
		sem: sync.Mutex{},
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
	}
}

// Balance selects a suitable host according
func (r *RoundRobin) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()
	if len(r.hosts) == 0 {
		return "", NoHostError
	}

	r.sem.Lock()
	r.i++
	host := r.hosts[r.i%uint64(len(r.hosts))]
	r.sem.Unlock()

	return host, nil
}
