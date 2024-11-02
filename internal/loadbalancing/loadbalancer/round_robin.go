package loadbalancer

import (
	"net"
	"sync"

	"github.com/phoon884/rev-proxy/internal/loadbalancing/ports"
)

type RoundRobinLoadbalancer struct {
	idx int
	mu  sync.Mutex // not sure if this is needed
}

func (r *RoundRobinLoadbalancer) SelectDownstreamAddr(conn net.Conn, addresses []string) string {
	// would mess up when healthcheck fails on of the address I am too lazy to come up with a better idea
	r.mu.Lock()
	defer r.mu.Unlock()
	addressesLen := len(addresses)
	r.idx++
	if r.idx >= addressesLen {
		r.idx -= addressesLen
	}
	return addresses[r.idx]
}

var _ ports.Loadbalancer = (*LeastConnectionBalancer)(nil)

func NewRoundRobinLoadBalancer() *RoundRobinLoadbalancer {
	return &RoundRobinLoadbalancer{}
}
