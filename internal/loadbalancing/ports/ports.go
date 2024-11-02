package ports

import "net"

type Loadbalancer interface {
	SelectDownstreamAddr(conn net.Conn, addresses []string) string
}
