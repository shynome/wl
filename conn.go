package wl

import (
	"bufio"
	"io"
	"math"
	"net"
	"time"
)

type Conn struct {
	rw   *bufio.ReadWriter
	conn io.ReadWriteCloser
}

var _ net.Conn = &Conn{}

func NewConn(conn io.ReadWriteCloser) *Conn {
	r := bufio.NewReaderSize(conn, math.MaxUint16)
	w := bufio.NewWriterSize(conn, math.MaxUint16)
	rw := bufio.NewReadWriter(r, w)
	return &Conn{
		rw:   rw,
		conn: conn,
	}
}

func (c *Conn) Read(b []byte) (n int, err error) {
	return c.rw.Read(b)
}
func (c *Conn) Write(b []byte) (n int, err error) {
	defer c.rw.Flush()
	return c.rw.Write(b)
}
func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) LocalAddr() net.Addr {
	return &Addr{}
}
func (c *Conn) RemoteAddr() net.Addr {
	return &Addr{}
}

// smux stream dial those
func (c *Conn) SetDeadline(t time.Time) error {
	return nil
}
func (c *Conn) SetReadDeadline(t time.Time) error {
	return nil
}
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return nil
}
