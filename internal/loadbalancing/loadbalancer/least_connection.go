package loadbalancer

import (
	"math"
	"net"
	"sync"

	"github.com/phoon884/rev-proxy/internal/loadbalancing/ports"
)

type LeastConnectionBalancer struct {
	connTable *sync.Map
}

func (l *LeastConnectionBalancer) SelectDownstreamAddr(conn net.Conn, addresses []string) string {
	bestAddr := ""
	lowestCount := int64(math.MaxInt64)
	for _, address := range addresses {
		var count int64
		structCount, ok := l.connTable.Load(address)
		if ok {
			count = *structCount.(*int64)
		} else {
			count = 0
		}

		if count < lowestCount {
			lowestCount = count
			bestAddr = address
		}
	}
	return bestAddr
}

var _ ports.Loadbalancer = (*LeastConnectionBalancer)(nil)

func NewLeastConnectionBalancer(connTable *sync.Map) *LeastConnectionBalancer {
	return &LeastConnectionBalancer{connTable: connTable}
}
