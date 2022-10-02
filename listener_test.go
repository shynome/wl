package wl

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"testing"

	"github.com/lainio/err2/try"
	"github.com/pion/webrtc/v3"
)

func TestListener(t *testing.T) {
	pc1, pc2 := try.To2(getConnectedPeerConnectionPair())
	go httpServer(pc2)

	session := try.To1(NewClientSession(pc1))

	client := &http.Client{
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
	if !bytes.Equal(b, helloResp) {
		t.Error(b)
	}

	resp2 := try.To1(client.Get("http://wl.com/big-file"))
	bigFileRecived := try.To1(io.ReadAll(resp2.Body))
	if !bytes.Equal(bigFileRecived, bigFile) {
		t.Error(b)
	}
}

var helloResp = []byte("world")
var bigFile = make([]byte, math.MaxUint16*16)

func httpServer(pc *webrtc.PeerConnection) {
	l := Listen()
	l.Add(&Peer{PC: pc})
	server := &http.ServeMux{}
	server.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write(helloResp)
	})
	server.HandleFunc("/big-file", func(w http.ResponseWriter, r *http.Request) {
		w.Write(bigFile)
	})
	http.Serve(l, server)
}
