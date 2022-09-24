package ortc

import (
	"testing"

	"github.com/lainio/err2/try"
	"github.com/pion/webrtc/v3"
)

func TestConnect(t *testing.T) {
	api := webrtc.NewAPI()
	config := webrtc.Configuration{}
	testData := []byte("hello")

	pc1 := try.To1(
		api.NewPeerConnection(config))
	defer pc1.Close()
	pc2 := try.To1(
		api.NewPeerConnection(config))
	defer pc2.Close()

	offer := try.To1(
		CreateOffer(pc1))
	roffer := try.To1(
		HandleConnect(pc2, offer))
	try.To(
		Handshake(pc1, roffer))

	dc := try.To1(
		pc2.CreateDataChannel("www", nil))
	defer dc.Close()

	try.To(
		Wait(dc))

	try.To(
		dc.Send(testData))

}
