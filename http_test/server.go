package main

import (
	"io"
	"net"
	"net/http"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/pion/webrtc/v3"
	"github.com/shynome/wl"
	"github.com/shynome/wl/ortc"
)

func main() {

	pc1, pc2 := try.To2(getConnectedPeerConnectionPair())
	l := wl.Listen()
	l.Add(&wl.Peer{PC: pc1})

	t := wl.NewTransport()
	s := try.To1(wl.NewClientSession(pc2))
	t.Set("wl", s)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "hello world")
	})
	go http.Serve(l, nil)

	l2 := try.To1(net.Listen("tcp", ":8080"))

	for {
		conn, err := l2.Accept()
		if err != nil {
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()
			wconn, err := t.NewConn("wl")
			if err != nil {
				return
			}
			go io.Copy(wconn, conn)
			io.Copy(conn, wconn)
		}(conn)
	}
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
