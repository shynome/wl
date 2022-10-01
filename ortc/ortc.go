package ortc

import (
	"fmt"
	"time"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/pion/webrtc/v3"
)

type Signal webrtc.SessionDescription

func CreateOffer(pc *webrtc.PeerConnection) (sdp Signal, err error) {
	defer err2.Return(&err)

	try.To(
		makeOfferWithCandidates(pc))

	offer := try.To1(
		pc.CreateOffer(nil))

	try.To(
		pc.SetLocalDescription(offer))
	sdp = Signal(*pc.LocalDescription())

	return
}

func makeOfferWithCandidates(pc *webrtc.PeerConnection) (err error) {
	defer err2.Return(&err)

	dc := try.To1(
		pc.CreateDataChannel("_for_collect_candidates", nil))
	defer dc.Close()

	wait := make(chan struct{})
	pc.OnNegotiationNeeded(func() {
		close(wait)
	})
	<-wait

	return
}

func HandleConnect(pc *webrtc.PeerConnection, offer Signal) (roffer Signal, err error) {
	defer err2.Return(&err)

	try.To(
		pc.SetRemoteDescription(webrtc.SessionDescription(offer)))

	answer := try.To1(
		pc.CreateAnswer(nil))

	gatherComplete := webrtc.GatheringCompletePromise(pc)

	try.To(
		pc.SetLocalDescription(answer))

	<-gatherComplete

	roffer = Signal(*pc.LocalDescription())

	return
}

func Handshake(pc *webrtc.PeerConnection, offer Signal) (err error) {
	defer err2.Return(&err)

	try.To(
		pc.SetRemoteDescription(webrtc.SessionDescription(offer)))

	return
}

func Wait(dc *webrtc.DataChannel) (err error) {
	defer err2.Return(&err)

	switch dc.ReadyState() {
	case webrtc.DataChannelStateOpen:
		return
	case webrtc.DataChannelStateClosing:
		fallthrough
	case webrtc.DataChannelStateClosed:
		err = fmt.Errorf("dc closed")
		return
	case webrtc.DataChannelStateConnecting:
		break
	}

	var (
		wait = make(chan struct{})
	)

	dc.OnOpen(func() {
		close(wait)
	})
	// dc.OnError(func(err error) {
	// 	errCh <- err
	// })

	select {
	case <-wait:
	case <-time.After(10 * time.Second):
		err = fmt.Errorf("wait dc timeout")
	}

	return
}
