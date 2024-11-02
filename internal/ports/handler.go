package ports

import "net"

type Handler interface {
	HandleConnection(conn net.Conn)
}
