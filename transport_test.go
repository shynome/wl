package wl

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/lainio/err2/try"
	"github.com/pion/webrtc/v3"
	"github.com/shynome/wl/ortc"
)

var testData = []byte("hello world")

func BenchmarkTransport(b *testing.B) {

	settingEngine := webrtc.SettingEngine{}
	settingEngine.DetachDataChannels()

	api := webrtc.NewAPI(
		webrtc.WithSettingEngine(settingEngine),
	)
	config := webrtc.Configuration{}

	pc1 := try.To1(
		api.NewPeerConnection(config))
	defer pc1.Close()
	pc2 := try.To1(
		api.NewPeerConnection(config))
	defer pc2.Close()

	offer := try.To1(
		ortc.CreateOffer(pc1))
	roffer := try.To1(
		ortc.HandleConnect(pc2, offer))
	try.To(
		ortc.Handshake(pc1, roffer))

	createServer(pc1)

	transport := &Transport{PC: pc2}
	client := http.Client{Transport: transport}

	for i := 0; i < b.N; i++ {
		resp := try.To1(
			client.Get("/hello"))
		respBytes := try.To1(
			io.ReadAll(resp.Body))

		if !bytes.Equal(testData, respBytes) {
			b.Error(respBytes)
		}
	}

}

func createServer(pc *webrtc.PeerConnection) (err error) {
	h := &http.ServeMux{}
	h.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write(testData)
	})
	var l = Listen()
	l.Add(&Peer{PC: pc})
	go http.Serve(l, h)
	return
}
