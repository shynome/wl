package wl

import (
	"fmt"
	"net"

	"github.com/pion/webrtc/v3"
)

type PairKey uint

const (
	PairLocal PairKey = iota
	PairRemote
)

type Addr struct {
	ICETransport *webrtc.ICETransport
	PairKey      PairKey
	ID           *uint16
}

var _ net.Addr = &Addr{}

func (a *Addr) Network() string {
	return "wl"
}

func (a *Addr) String() string {
	var addr string
	if a.ICETransport != nil {
		if pair, err := a.ICETransport.GetSelectedCandidatePair(); err == nil {
			switch a.PairKey {
			case PairLocal:
				addr = pair.Local.String()
			case PairRemote:
				addr = pair.Remote.String()
			}
		}
	}
	if a.ID != nil {
		addr = fmt.Sprintf("(%v) %s", *a.ID, addr)
	}
	return addr
}
