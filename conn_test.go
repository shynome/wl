package wl

import (
	"io"
	"net"
	"testing"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/pion/webrtc/v3"
	"github.com/shynome/wl/ortc"
	"github.com/xtaci/smux"
)

func TestSmuxConn(t *testing.T) {
	pc1, pc2 := try.To2(getConnectedPeerConnectionPair())
	pc2.OnDataChannel(func(dc *webrtc.DataChannel) {
		dc.OnOpen(func() {
			var err error
			defer err2.Return(&err)
			conn := try.To1(dc.Detach())
			session := try.To1(smux.Server(NewConn(conn), nil))
			for {
				conn := try.To1(session.Accept()).(net.Conn)
				go func(conn net.Conn) {
					defer conn.Close()
					try.To1(io.WriteString(conn, "world"))
				}(conn)
			}
		})
	})
	dc := try.To1(pc1.CreateDataChannel(SmuxLabel, nil))
	try.To(ortc.Wait(dc))
	conn := try.To1(dc.Detach())
	session := try.To1(smux.Client(NewConn(conn), nil))
	kconn := try.To1(session.Open()).(net.Conn)
	b := try.To1(io.ReadAll(kconn))
	t.Log(b)
}

func getConnectedPeerConnectionPair() (pc1 *webrtc.PeerConnection, pc2 *webrtc.PeerConnection, err error) {
	defer err2.Return(&err)

	settingEngine := webrtc.SettingEngine{}
	settingEngine.DetachDataChannels()

	api := webrtc.NewAPI(webrtc.WithSettingEngine(settingEngine))
	config := webrtc.Configuration{}

	pc1 = try.To1(api.NewPeerConnection(config))
	pc2 = try.To1(api.NewPeerConnection(config))

	offer := try.To1(ortc.CreateOffer(pc1))
	roffer := try.To1(ortc.HandleConnect(pc2, offer))
	try.To(ortc.Handshake(pc1, roffer))

	return
}
