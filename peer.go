package wl

import (
	"net"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/pion/webrtc/v3"
)

type Peer struct {
	PC *webrtc.PeerConnection
}

func (p *Peer) ForwardConns(conns chan net.Conn) (err error) {
	pc := p.PC

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		dc.OnOpen(func() {
			conn, err := dc.Detach()
			if err != nil {
				return
			}
			conns <- NewConn(dc, conn)
		})
	})

	return
}

func (p *Peer) Close() (err error) {
	defer err2.Return(&err)

	try.To(
		p.PC.Close())

	return
}
