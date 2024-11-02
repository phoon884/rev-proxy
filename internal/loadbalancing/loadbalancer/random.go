package loadbalancer

import (
	"math/rand/v2"
	"net"

	"github.com/phoon884/rev-proxy/internal/loadbalancing/ports"
)

type RandomLoadbalancer struct{}

func (r *RandomLoadbalancer) SelectDownstreamAddr(conn net.Conn, addresses []string) string {
	return addresses[rand.IntN(len(addresses))]
}

var _ ports.Loadbalancer = (*RandomLoadbalancer)(nil)

func NewRandomLoadbalancer() *RandomLoadbalancer {
	return &RandomLoadbalancer{}
}
