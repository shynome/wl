package wl

import (
	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/pion/webrtc/v3"
	"github.com/shynome/wl/ortc"
	"github.com/xtaci/smux"
)

func NewClientSession(pc *webrtc.PeerConnection) (session *smux.Session, err error) {
	defer err2.Return(&err)

	dc := try.To1(pc.CreateDataChannel(SmuxLabel, nil))
	ortc.Wait(dc)
	conn := try.To1(dc.Detach())
	session = try.To1(smux.Client(NewConn(conn), nil))

	return
}
