package healthcheck

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/phoon884/rev-proxy/internal/httpHdl/models"
	"github.com/phoon884/rev-proxy/internal/httpHdl/utils"
	logger "github.com/phoon884/rev-proxy/pkg/logger/application"
)

type healthStatus struct {
	Failed int
	Pass   int
	Status atomic.Bool
}

type HttpHealthcheck struct {
	interval     int
	fail         int
	pass         int
	healthStatus map[string]*healthStatus
	logger       *logger.Logger
}

func (h *HttpHealthcheck) StartHealthcheck() {
	ticker := time.NewTicker(time.Second * time.Duration(h.interval))
	go h.healthcheckAll()
	for {
		select {
		case <-ticker.C:
			go h.healthcheckAll()
		}
	}
}

func (h *HttpHealthcheck) healthcheckAll() {
	wg := sync.WaitGroup{}
	for addr, status := range h.healthStatus {
		wg.Add(1)
		go h.healthcheck(addr, status, &wg)
	}
	wg.Wait()
}

func (h *HttpHealthcheck) healthcheck(addr string, status *healthStatus, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := net.Dial("tcp", addr)
	var res *models.HTTPRes
	if conn != nil {
		header := make(map[string]string, 1)
		header["Host"] = addr
		payload := models.HTTPReq{Method: "GET", Path: "/", Header: header}

		conn.Write([]byte(payload.StrRepr()))
		reader := bufio.NewReader(conn)
		res, err = utils.ParseRes(reader)
	}
	if err != nil {
		res = &models.HTTPRes{ResponseCode: 503}
	}
	if res.ResponseCode == 200 {
		h.logger.Debug("HEALTHCHECK: PASSED", addr)
		status.Pass += 1
	} else {

		fmt.Println("dialing")
		h.logger.Debug("HEALTHCHECK: FAILED", addr)
		status.Failed += 1
	}
	if status.Pass >= h.pass {
		if !status.Status.Load() {
			h.logger.Info("HEALTHCHECK STATUS BACK UP", addr)
		}
		status.Status.Store(true)
		status.Failed = 0
		status.Pass = 0
	} else if status.Failed >= h.fail {
		if status.Status.Load() {
			h.logger.Warn("HEALTHCHECK STATUS WENT DOWN", addr)
		}
		status.Status.Store(false)
		status.Failed = 0
		status.Pass = 0
	}
}

func (h *HttpHealthcheck) GetStatus(addr string) bool {
	var status bool
	if h.healthStatus[addr] != nil {
		status = h.healthStatus[addr].Status.Load()
	}
	return status
}

func NewHttpHealthcheck(
	interval int,
	fail int,
	pass int,
	addrList []string,
	logger *logger.Logger,
) *HttpHealthcheck {
	ret := HttpHealthcheck{
		interval:     interval,
		fail:         fail,
		pass:         pass,
		logger:       logger,
		healthStatus: make(map[string]*healthStatus, len(addrList)),
	}
	for _, addr := range addrList {
		ret.healthStatus[addr] = &healthStatus{}
	}
	go ret.StartHealthcheck()
	return &ret
}
