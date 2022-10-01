package wl

import (
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"testing"

	"github.com/lainio/err2/try"
	"github.com/pion/webrtc/v3"
	"github.com/shynome/wl/ortc"
	"github.com/xtaci/smux"
)

func TestListener(t *testing.T) {
	pc1, pc2 := try.To2(getConnectedPeerConnectionPair())
	go httpServer(pc2)
	dc := try.To1(pc1.CreateDataChannel(SmuxLabel, nil))
	try.To(ortc.Wait(dc))
	conn := try.To1(dc.Detach())
	session := try.To1(smux.Client(NewConn(conn), nil))
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (conn net.Conn, err error) {
				conn, ok := try.To1(session.Open()).(net.Conn)
				if !ok {
					err = fmt.Errorf("session open failed")
				}
				return
			},
		},
	}
	resp := try.To1(client.Get("http://wl.com/hello"))
	b := try.To1(io.ReadAll(resp.Body))
	t.Log(b)

	resp2 := try.To1(client.Get("http://wl.com/big-file"))
	bigFile := try.To1(io.ReadAll(resp2.Body))
	t.Log(bigFile)
}

var bigFile = make([]byte, math.MaxUint16*16)

func httpServer(pc *webrtc.PeerConnection) {
	l := Listen()
	l.Add(&Peer{PC: pc})
	server := &http.ServeMux{}
	server.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "world")
	})
	server.HandleFunc("/big-file", func(w http.ResponseWriter, r *http.Request) {
		w.Write(bigFile)
	})
	http.Serve(l, server)
}
