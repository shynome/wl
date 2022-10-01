package wl

import (
	"net"
)

type Addr struct {
	Label string
}

var _ net.Addr = &Addr{}

func (a *Addr) Network() string {
	return "wl"
}

func (a *Addr) String() string {
	return a.Label
}
