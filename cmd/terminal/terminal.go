package main

import (
	"flag"
	"fmt"
	"github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/gosown/wwwping/config"
	"github.com/gosown/wwwping/process"
	"log"
	"math"
	"time"
)

var usage = `
Usage:

    ping [-c count] [-i interval] [-t timeout] [--privileged] host

Examples:

    # ping google continuously
    ping www.google.com

    # ping google 5 times
    ping -c 5 www.google.com

    # ping google 5 times at 500ms intervals
    ping -c 5 -i 500ms www.google.com

    # ping google for 10 seconds
    ping -t 10s www.google.com

    # Send a privileged raw ICMP ping
    sudo ping --privileged www.google.com

    # Send ICMP messages with a 100-byte payload
    ping -s 100 1.1.1.1
`

func main() {
	timeout := flag.Duration("t", time.Second*100000, "")
	interval := flag.Duration("i", time.Second, "")
	count := flag.Int("c", -1, "")
	size := flag.Int("s", 24, "")
	ttl := flag.Int("l", 64, "TTL")
	privileged := flag.Bool("privileged", false, "")
	flag.Usage = func() {
		fmt.Print(usage)
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	host := flag.Arg(0)
	wpcfg := config.WwwPing{
		Host:       host,
		Interval:   *interval,
		Timeout:    *timeout,
		Count:      *count,
		Size:       *size,
		TTL:        *ttl,
		Privileged: *privileged,
	}

	terminalCfg := config.Terminal{WwwPing: wpcfg, AxesColor: termui.ColorWhite, LineColor: termui.ColorYellow}
	terminal, err := process.NewTerminal(terminalCfg)
	if err != nil {
		log.Fatalln(err)
	}
	err = terminal.Run()
	if err != nil {
		log.Fatalln(err)
	}
	terminal.Statistics()
}

func x() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	sinFloat64 := (func() []float64 {
		n := 400
		data := make([]float64, n)
		for i := range data {
			data[i] = 1 + math.Sin(float64(i)/5)
		}
		return data
	})()

	lc := widgets.NewPlot()
	lc.Title = "braille-mode Line Chart"
	lc.Data = append(lc.Data, sinFloat64)
	lc.AxesColor = ui.ColorWhite
	lc.LineColors[0] = ui.ColorYellow

	grid := ui.NewGrid()
	termWidth, termHeight := ui.TerminalDimensions()
	grid.SetRect(0, 0, termWidth, termHeight)

	grid.Set(
		ui.NewRow(1.0/2,
			ui.NewCol(1.0/2, lc),
		),
	)

	ui.Render(grid)

	tickerCount := 1
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		case <-ticker:
			if tickerCount == 100 {
				return
			}
			lc.Data[0] = sinFloat64[2*tickerCount:]
			ui.Render(grid)
			tickerCount++
		}
	}
}
