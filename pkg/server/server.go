package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	hdl "github.com/phoon884/rev-proxy/internal/ports"
	logger "github.com/phoon884/rev-proxy/pkg/logger/ports"
)

type Server struct {
	name      string
	handler   hdl.Handler
	logger    logger.LoggerApplication
	timeout   int
	waitGroup *sync.WaitGroup
}

func (s *Server) StartServer(port int) {
	defer s.waitGroup.Done()
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.logger.Error("Error while creating server with error message:", err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("Fail to setup server with error:", err.Error())
			continue
		}
		s.logger.Info("Recieve a connection on server:", s.name)
		conn.SetDeadline(time.Now().Add(time.Duration(s.timeout) * time.Second))

		go s.handler.HandleConnection(conn)
		s.logger.Info("handled connection")
	}
}

func NewServer(
	name string,
	handler hdl.Handler,
	logger logger.LoggerApplication,
	timeout int,
	waitGroup *sync.WaitGroup,
) *Server {
	return &Server{
		name:      name,
		handler:   handler,
		logger:    logger,
		timeout:   timeout,
		waitGroup: waitGroup,
	}
}
