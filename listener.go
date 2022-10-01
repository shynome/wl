package wl

import (
	"errors"
	"net"
	"sync"

	"github.com/lainio/err2"
)

type peersMap struct {
	lock  *sync.Mutex
	value map[string]*Peer
}

type Listener struct {
	peers *peersMap
	conns chan net.Conn
}

func Listen() *Listener {
	return &Listener{
		peers: &peersMap{
			lock:  &sync.Mutex{},
			value: map[string]*Peer{},
		},
		conns: make(chan net.Conn, 10),
	}
}

func (l *Listener) Accept() (conn net.Conn, err error) {
	conn, ok := <-l.conns
	if !ok {
		return nil, errors.New("Closed")
	}
	return conn, nil
}

func (l *Listener) Close() error {
	l.closePeers()
	close(l.conns)
	return nil
}

func (l *Listener) Add(peer *Peer) (err error) {
	defer err2.Return(&err)

	peer.ForwardConns(l.conns)

	l.peers.lock.Lock()
	defer l.peers.lock.Unlock()

	k := peer.PC.RemoteDescription().SDP
	l.peers.value[k] = peer

	return
}

func (l *Listener) closePeers() {
	l.peers.lock.Lock()
	defer l.peers.lock.Unlock()

	var wg sync.WaitGroup
	for _, p := range l.peers.value {
		wg.Add(1)
		go func(p *Peer) {
			defer wg.Done()
			p.Close()
		}(p)
	}
	wg.Wait()
}

func (l *Listener) Addr() net.Addr {
	return &Addr{}
}
