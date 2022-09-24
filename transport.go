package wl

import (
	"bufio"
	"io"
	"net/http"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/pion/datachannel"
	"github.com/pion/webrtc/v3"
)

type Transport struct {
	PC *webrtc.PeerConnection
}

var _ http.RoundTripper = &Transport{}

func (t *Transport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	defer closeBody(req.Body)
	defer err2.Return(&err)

	conn := try.To1(
		t.NewConn())
	defer conn.Close()

	try.To(
		req.Write(conn))

	res = try.To1(
		http.ReadResponse(bufio.NewReader(conn), req))

	return
}

func (t *Transport) NewConn() (conn io.ReadWriteCloser, err error) {
	defer err2.Return(&err)

	var (
		resultCh = make(chan datachannel.ReadWriteCloser)
		errCh    = make(chan error)
	)

	dc := try.To1(
		t.PC.CreateDataChannel("", nil))

	dc.OnOpen(func() {
		conn, err := dc.Detach()
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- conn
	})
	dc.OnError(func(err error) {
		errCh <- err
	})

	select {
	case conn = <-resultCh:
	case err = <-errCh:
	}

	return
}

func closeBody(body io.ReadCloser) {
	if body == nil {
		return
	}
	body.Close()
}
