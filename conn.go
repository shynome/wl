package wl

import (
	"io"
	"net"
	"os"
	"time"

	"github.com/pion/webrtc/v3"
)

type Conn struct {
	Conn        io.ReadWriteCloser
	DataChannel *webrtc.DataChannel
	timeoutOf   timeoutOf
}

type timeoutOf struct {
	read  time.Time
	write time.Time
}

var _ net.Conn = &Conn{}

func NewConn(dc *webrtc.DataChannel, conn io.ReadWriteCloser) *Conn {
	return &Conn{
		Conn:        conn,
		DataChannel: dc,
	}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	if isDeadline(c.timeoutOf.read) {
		return 0, os.ErrDeadlineExceeded
	}
	return c.Conn.Read(b)
}
func (c *Conn) Write(b []byte) (n int, err error) {
	if isDeadline(c.timeoutOf.write) {
		return 0, os.ErrDeadlineExceeded
	}
	return c.Conn.Write(b)
}
func (c *Conn) Close() error {
	return c.Conn.Close()
}

func (c *Conn) getICETransport() *webrtc.ICETransport {
	return c.DataChannel.Transport().Transport().ICETransport()
}
func (c *Conn) LocalAddr() net.Addr {
	return &Addr{
		// ICETransport: c.getICETransport(),
		// PairKey:      PairLocal,
		// ID:           c.DataChannel.ID(),
	}
}
func (c *Conn) RemoteAddr() net.Addr {
	return &Addr{
		// ICETransport: c.getICETransport(),
		// PairKey:      PairRemote,
		// ID:           c.DataChannel.ID(),
	}
}

func isDeadline(t time.Time) bool {
	return !t.IsZero() && time.Now().After(t)
}

func (c *Conn) SetDeadline(t time.Time) error {
	c.SetReadDeadline(t)
	c.SetWriteDeadline(t)
	return nil
}
func (c *Conn) SetReadDeadline(t time.Time) error {
	c.timeoutOf.read = t
	return nil
}
func (c *Conn) SetWriteDeadline(t time.Time) error {
	c.timeoutOf.write = t
	return nil
}
