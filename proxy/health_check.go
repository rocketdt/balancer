// Copyright 2022 <mzh.scnu@qq.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"fmt"
	"os"
	"time"
)

// ReadAlive reads the alive status of the site
func (h *HTTPProxy) ReadAlive(url string) bool {
	h.RLock()
	defer h.RUnlock()
	return h.alive[url]
}

// SetAlive sets the alive status to the site
func (h *HTTPProxy) SetAlive(url string, alive bool) {
	h.Lock()
	defer h.Unlock()
	h.alive[url] = alive
}

// HealthCheck enable a health check goroutine for each agent
func (h *HTTPProxy) HealthCheck(interval uint) {
	for host := range h.hostMap {
		go h.healthCheck(host, interval)
	}
}

func (h *HTTPProxy) healthCheck(host string, interval uint) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for range ticker.C {
		alive := false

		if h.fptr != nil {
			alive = h.fptr(host)
		} else {
			alive = IsBackendAlive(host)
		}

		readAlive := h.ReadAlive(host)

		// log.Printf("Alive = %+v, Read alive = %+v\n", alive, readAlive)

		if !alive && readAlive {
			// log.Printf("Site unreachable, remove %s from load balancer.", host)
			fmt.Fprintf(os.Stderr, "[%s] Site unreachable, remove %s from load balancer\n",
				time.Now().Format("2006-01-02 15:04:05"), host)

			h.SetAlive(host, false)
			h.lb.Remove(host)
		} else if alive && !readAlive {
			// log.Printf("Site reachable, add %s to load balancer.", host)
			fmt.Fprintf(os.Stderr, "[%s] Site reachable, add %s from load balancer\n",
				time.Now().Format("2006-01-02 15:04:05"), host)

			h.SetAlive(host, true)
			h.lb.Add(host)
		}
	}

}
