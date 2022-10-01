package wl

import (
	"fmt"
	"net"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/pion/webrtc/v3"
	"github.com/xtaci/smux"
)

type Peer struct {
	PC      *webrtc.PeerConnection
	Session *smux.Session
}

const SmuxLabel = "smux"

func (p *Peer) ForwardConns(conns chan net.Conn) (err error) {
	pc := p.PC

	pc.OnDataChannel(func(dc *webrtc.DataChannel) {
		if dc.Label() != SmuxLabel {
			return
		}
		dc.OnOpen(func() {
			var err error
			defer func() {
				if err != nil {
					fmt.Println("err:", err)
				}
			}()
			defer err2.Return(&err)
			conn := try.To1(dc.Detach())
			session := try.To1(smux.Server(NewConn(conn), nil))
			defer session.Close()
			p.Session = session
			for {
				stream := try.To1(session.Accept())
				conn, ok := stream.(net.Conn)
				if !ok {
					break
				}
				conns <- conn
			}
		})
	})

	return
}

func (p *Peer) Close() (err error) {
	defer err2.Return(&err)

	try.To(p.PC.Close())
	try.To(p.Session.Close())

	return
}
