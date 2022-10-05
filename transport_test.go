package wl

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"testing"

	"github.com/lainio/err2/try"
)

func TestTranposrt(t *testing.T) {
	pc1, pc2 := try.To2(getConnectedPeerConnectionPair())
	go httpServer(pc2)

	session := try.To1(NewClientSession(pc1))

	tt := NewTransport()
	tt.Set("wl.com", session)

	client := &http.Client{
		Transport: tt,
	}
	client2 := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return tt.NewConn(addr)
			},
		},
	}

	testClient(t, client)
	testClient(t, client)

	t.Log(client, client2)
}

func testClient(t *testing.T, client *http.Client) {
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

func BenchmarkTranposrt(b *testing.B) {
	pc1, pc2 := try.To2(getConnectedPeerConnectionPair())
	go httpServer(pc2)

	session := try.To1(NewClientSession(pc1))

	tt := NewTransport()
	tt.Set("wl.com", session)

	client := &http.Client{
		Transport: tt,
	}

	for i := 0; i < b.N; i++ {
		resp := try.To1(client.Get("http://wl.com/hello"))
		respBytes := try.To1(io.ReadAll(resp.Body))
		if !bytes.Equal(respBytes, helloResp) {
			b.Error(respBytes)
		}
	}
}

func BenchmarkHttpTranposrt(b *testing.B) {
	pc1, pc2 := try.To2(getConnectedPeerConnectionPair())
	go httpServer(pc2)

	session := try.To1(NewClientSession(pc1))

	tt := NewTransport()
	tt.Set("wl.com:80", session)

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return tt.NewConn(addr)
			},
		},
	}

	for i := 0; i < b.N; i++ {
		resp := try.To1(client.Get("http://wl.com/hello"))
		respBytes := try.To1(io.ReadAll(resp.Body))
		if !bytes.Equal(respBytes, helloResp) {
			b.Error(respBytes)
		}
	}
}
