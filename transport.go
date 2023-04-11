package wl

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"

	"github.com/lainio/err2"
	"github.com/lainio/err2/try"
	"github.com/xtaci/smux"
)

type sessionMap struct {
	mu    *sync.RWMutex
	value map[string]*smux.Session
}

type Transport struct {
	sessions *sessionMap
}

var _ http.RoundTripper = (*Transport)(nil)

func NewTransport() (t *Transport) {
	t = &Transport{
		sessions: &sessionMap{
			mu:    &sync.RWMutex{},
			value: map[string]*smux.Session{},
		},
	}
	return
}

func (t *Transport) RoundTrip(req *http.Request) (res *http.Response, err error) {
	defer closeBody(req.Body)
	defer err2.Handle(&err)

	conn := try.To1(t.NewConn(req.Host))
	try.To(req.Write(conn))
	res = try.To1(http.ReadResponse(bufio.NewReader(conn), req))

	return
}

func (t *Transport) NewConn(addr string) (conn net.Conn, err error) {
	defer err2.Handle(&err)
	session := try.To1(t.Get(addr))
	conn, ok := try.To1(session.Open()).(net.Conn)
	if !ok {
		err = fmt.Errorf("")
		return
	}
	return
}

var ErrSessionNotExists = fmt.Errorf("session is not exists")

func (t *Transport) Get(addr string) (session *smux.Session, err error) {
	t.sessions.mu.RLock()
	defer t.sessions.mu.RUnlock()
	session, ok := t.sessions.value[addr]
	if !ok || session == nil {
		err = fmt.Errorf("%w", ErrSessionNotExists)
	}
	return
}

func (t *Transport) Set(addr string, session *smux.Session) {
	t.sessions.mu.Lock()
	defer t.sessions.mu.Unlock()
	t.sessions.value[addr] = session
}

func closeBody(body io.ReadCloser) {
	if body == nil {
		return
	}
	body.Close()
}
