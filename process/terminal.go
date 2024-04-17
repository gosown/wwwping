package process

import (
	"fmt"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/gosown/wwwping"
	"github.com/gosown/wwwping/config"
	probing "github.com/gosown/wwwping/ping"
	"time"
)

type ui struct {
	plot *widgets.Plot
	grid *termui.Grid
}

type Terminal struct {
	wwwPing *wwwping.WwwPing
	ui      *ui
}

func NewTerminal(cfg config.Terminal) (terminal *Terminal, err error) {
	if err = termui.Init(); err != nil {
		return
	}

	wwwPing, err := wwwping.NewWwwPing(cfg.WwwPing)
	if err != nil {
		return
	}

	lc := widgets.NewPlot()
	lc.Title = fmt.Sprintf("PING %s (%s):\n", wwwPing.Pinger.Addr(), wwwPing.Pinger.IPAddr())
	lc.Data = [][]float64{{0, 0}}
	lc.AxesColor = cfg.AxesColor
	lc.LineColors[0] = cfg.LineColor

	grid := termui.NewGrid()
	termWidth, termHeight := termui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		termui.NewRow(1.0,
			termui.NewCol(1.0, lc),
		),
	)

	terminal = &Terminal{wwwPing: wwwPing, ui: &ui{plot: lc, grid: grid}}
	return
}

func (terminal *Terminal) onRecv(pkt *probing.Packet) {
	terminal.ui.plot.Data[0] = append(terminal.ui.plot.Data[0], float64(pkt.Rtt)/float64(time.Millisecond))
	termui.Render(terminal.ui.grid)
}

func (terminal *Terminal) onDuplicateRecv(pkt *probing.Packet) {}

func (terminal *Terminal) Statistics() {
	statistics := terminal.wwwPing.Statistics
	fmt.Printf("--- %s ping statistics ---\n", statistics.Addr)
	fmt.Printf("%d packets transmitted, %d packets received, %d duplicates, %v%% packet loss\n",
		statistics.PacketsSent, statistics.PacketsRecv, statistics.PacketsRecvDuplicates, statistics.PacketLoss)
	fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
		statistics.MinRtt, statistics.AvgRtt, statistics.MaxRtt, statistics.StdDevRtt)
}

func (terminal *Terminal) Run() (err error) {
	defer termui.Close()
	termui.Render(terminal.ui.grid)
	uiEvents := termui.PollEvents()
	go func() {
		for {
			select {
			case e := <-uiEvents:
				switch e.ID {
				case "q", "<C-c>":
					terminal.wwwPing.Stop()
				case "<Resize>":
					payload := e.Payload.(termui.Resize)
					terminal.ui.grid.SetRect(0, 0, payload.Width, payload.Height)
					termui.Clear()
					termui.Render(terminal.ui.grid)
				}
			case pkt := <-terminal.wwwPing.OnRecvChan:
				terminal.onRecv(pkt)
			case pkt := <-terminal.wwwPing.OnDuplicateRecvChan:
				terminal.onDuplicateRecv(pkt)
			}
		}
	}()

	err = terminal.wwwPing.Run()
	if err != nil {
		return
	}

	return
}
