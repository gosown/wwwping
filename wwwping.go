package wwwping

import (
	"github.com/gosown/wwwping/config"
	probing "github.com/gosown/wwwping/ping"
)

type WwwPing struct {
	Pinger              *probing.Pinger
	OnRecvChan          chan *probing.Packet
	OnDuplicateRecvChan chan *probing.Packet
	Statistics          *probing.Statistics
}

func NewWwwPing(cfg config.WwwPing) (wwwPing *WwwPing, err error) {
	pinger, err := probing.NewPinger(cfg.Host)
	if err != nil {
		return
	}

	wwwPing = &WwwPing{
		OnRecvChan:          make(chan *probing.Packet, 100), // TODO: 100?
		OnDuplicateRecvChan: make(chan *probing.Packet, 100), // TODO: 100?
		Statistics:          &probing.Statistics{},
	}

	pinger.Count = cfg.Count
	pinger.Size = cfg.Size
	pinger.Interval = cfg.Interval
	pinger.Timeout = cfg.Timeout
	pinger.TTL = cfg.TTL
	pinger.SetPrivileged(cfg.Privileged)

	pinger.OnRecv = func(pkt *probing.Packet) {
		wwwPing.OnRecvChan <- pkt
	}
	pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
		wwwPing.OnDuplicateRecvChan <- pkt
	}
	pinger.OnFinish = func(stats *probing.Statistics) {
		wwwPing.Statistics = stats
	}

	wwwPing.Pinger = pinger

	return
}

func (wwwPing *WwwPing) Run() error {
	return wwwPing.Pinger.Run()
}

func (wwwPing *WwwPing) Stop() {
	wwwPing.Pinger.Stop()
}
