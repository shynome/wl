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
	defer err2.Handle(&err)

	peer.ForwardConns(l.conns)

	l.peers.lock.Lock()
	defer l.peers.lock.Unlock()

	k := peer.Key()
	l.peers.value[k] = peer

	return
}

func (l *Listener) Remove(peer *Peer) (err error) {
	defer err2.Handle(&err)

	l.peers.lock.Lock()
	defer l.peers.lock.Unlock()

	k := peer.Key()
	delete(l.peers.value, k)

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
