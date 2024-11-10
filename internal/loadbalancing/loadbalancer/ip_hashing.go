package loadbalancer

import (
	"crypto/md5"
	"encoding/binary"
	"net"

	"github.com/phoon884/rev-proxy/internal/loadbalancing/ports"
)

type IPHashingLoadbalancer struct {
	DownStreamAddr []string
}

func (r *IPHashingLoadbalancer) SelectDownstreamAddr(conn net.Conn, hcValidAddr []string) string {
	ipAddr, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
	hash := md5.Sum([]byte(ipAddr))
	idx := binary.BigEndian.Uint16(hash[:]) % uint16(len(r.DownStreamAddr))
	addr := r.DownStreamAddr[idx]
	found := false
	for i := idx; i >= 0; i-- {
		// check if Ip is still in the list of healthcheck passed addr
		if addr == hcValidAddr[i] {
			found = true
			break
		}
	}
	if !found {
		// remodulo with list of hcValidAddr instead if the original list is invalid
		// assumes that healthcheck doesn't fail often and when it fails other downstream server will still...
		// be receiving the same set of users.
		idx = binary.BigEndian.Uint16(hash[:]) % uint16(len(hcValidAddr))
		return hcValidAddr[idx]
	}
	return addr
}

var _ ports.Loadbalancer = (*IPHashingLoadbalancer)(nil)

func NewIPHashingLoadbalancer(downStreamAddr []string) *IPHashingLoadbalancer {
	return &IPHashingLoadbalancer{DownStreamAddr: downStreamAddr}
}
