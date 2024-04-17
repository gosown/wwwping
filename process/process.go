package process

import probing "github.com/gosown/wwwping/ping"

type Process interface {
	onRecv(pkt *probing.Packet)
	onDuplicateRecv(pkt *probing.Packet)
	Run() error
}
