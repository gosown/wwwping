package config

import (
	"github.com/gizak/termui/v3"
	"time"
)

type WwwPing struct {
	// TODO: The data type here should ideally be defined consistently with that in the pinger.
	Host       string
	Interval   time.Duration
	Timeout    time.Duration
	Count      int
	Size       int
	TTL        int
	Privileged bool
}

type Terminal struct {
	AxesColor termui.Color
	LineColor termui.Color
	WwwPing
}

type Web struct {
	WwwPing
}
