package handler

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/phoon884/rev-proxy/internal/httpHdl/models"
	"github.com/phoon884/rev-proxy/internal/httpHdl/utils"
	logger "github.com/phoon884/rev-proxy/pkg/logger/ports"
)

type HTTPHdl struct {
	logger               logger.LoggerApplication
	endpoints            []*Endpoint
	DownStreamConnection *sync.Map
}

func (h *HTTPHdl) HandleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := bufio.NewReader(conn)

	message, err := utils.ParseReq(buffer)
	if err != nil {
		h.handleError(conn, 400, "Invalid HTTP")
		return
	}
	if message == nil {
		return
	}
	address := ""
	var found bool
	for _, endpoint := range h.endpoints {
		result, err := endpoint.ComparePath(message)
		if err != nil {
			h.handleError(conn, 500, err.Error())
			return
		}
		if result {
			endpoint.setHeader(message)
			address = endpoint.GetDownStreamAddr(conn)
			found = true
			if endpoint.Ratelimter != nil {
				clientIp, _, _ := net.SplitHostPort(conn.RemoteAddr().String())
				if allow, err := endpoint.Ratelimter.Allow(clientIp); !allow {
					if err != nil {
						h.handleError(conn, 500, err.Error())
					} else {
						h.handleError(conn, 429, fmt.Sprint("rate limit reached"))
					}
					return
				}
			}
			break
		}
	}
	if !found {
		h.handleError(conn, 404, "Not found")
		return
	}
	if address == "" {
		h.handleError(conn, 503, "Service unavalible") // All health check fails
		return
	}

	downStream, err := net.Dial("tcp", address)
	if err != nil {
		h.handleError(conn, 503, "Service unavalible")
		return
	}
	downStream.SetDeadline(time.Now().Add(time.Duration(3) * time.Second))
	val, _ := h.DownStreamConnection.LoadOrStore(address, new(int64))
	ptr := val.(*int64)
	atomic.AddInt64(ptr, 1)

	defer downStream.Close()
	defer atomic.AddInt64(ptr, -1)

	downStream.Write([]byte(message.StrRepr()))

	ctxLen, err := strconv.ParseInt(message.GetHeader("Content-Length"), 10, 64)
	if err != nil {
		h.handleError(conn, 400, "Context-Length not numeric")
		return
	}
	io.CopyN(downStream, buffer, ctxLen)

	downStreamReader := bufio.NewReader(downStream)
	responseInfo, err := utils.ParseRes(downStreamReader)
	if err != nil {
		h.handleError(conn, 502, "Reverse Proxy can't read response")
		return
	}
	h.logger.Info("", string(responseInfo.ToBytes()))
	conn.Write(responseInfo.ToBytes())
	ctxLen, err = strconv.ParseInt(responseInfo.GetHeader("Content-Length"), 10, 64)
	if err != nil {
		h.handleError(conn, 502, "Reverse Proxy can't read content lenght")
		return
	}
	io.CopyN(conn, downStreamReader, ctxLen)
	h.logger.Info("", message.Method, message.Path, responseInfo.ResponseCode)
}

func (h *HTTPHdl) handleError(conn net.Conn, code int, msg string) {
	response := models.NewErrorFound(code, msg)
	conn.Write([]byte(response.ToBytes()))
}

func NewHTTPHdl(
	endpoints []*Endpoint,
	logger logger.LoggerApplication,
	downStreamConnection *sync.Map,
) *HTTPHdl {
	return &HTTPHdl{
		endpoints:            endpoints,
		logger:               logger,
		DownStreamConnection: downStreamConnection,
	}
}
