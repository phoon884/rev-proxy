package loadbalancer

import (
	"crypto/md5"
	"encoding/binary"
	"net"

	"github.com/phoon884/rev-proxy/internal/loadbalancing/ports"
)

type IPHashingLoadbalancer struct{}

func (r *IPHashingLoadbalancer) SelectDownstreamAddr(conn net.Conn, addresses []string) string {
	ipAddr := conn.RemoteAddr().String()
	hash := md5.Sum([]byte(ipAddr))
	idx := binary.BigEndian.Uint16(hash[:]) % uint16(len(addresses))
	return addresses[idx]
}

var _ ports.Loadbalancer = (*IPHashingLoadbalancer)(nil)

func NewIPHashingLoadbalancer() *IPHashingLoadbalancer {
	return &IPHashingLoadbalancer{}
}
