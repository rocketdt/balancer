// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package balancer

func init() {
	factories[FirstBalancer] = NewFirst
}

type First struct {
	BaseBalancer
}

// NewRandom create new Random balancer
func NewFirst(hosts []string) Balancer {
	return &First{
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
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
