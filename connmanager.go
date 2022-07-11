package tcpserver

import (
	"fmt"
	"log"
	"sync"
)

type ConnManager interface {
	Add(conn Conn) error
	Remove(conn Conn)
	Clear()
}

type connmanager struct {
	conns map[uint32]Conn
	sync.Mutex
}

func NewConnManager() ConnManager {
	return &connmanager{
		conns: make(map[uint32]Conn),
	}
}

func (mgr *connmanager) Add(conn Conn) error {
	mgr.Lock()
	defer mgr.Unlock()
	// conversion is okay?
	if uint32(len(mgr.conns)) == config.tcpserver.maxconn {
		return fmt.Errorf("Maxconn %d is reached", config.tcpserver.maxconn)
	}
	mgr.conns[conn.ConnId()] = conn
	log.Printf("connmgr Add conn [%d]", conn.ConnId())
	return nil
}

func (mgr *connmanager) Remove(conn Conn) {
	mgr.Lock()
	defer mgr.Unlock()
	delete(mgr.conns, conn.ConnId())
	log.Printf("connmgr Remove conn [%d]", conn.ConnId())
	// stop ???
}

func (mgr *connmanager) Clear() {
	mgr.Lock()
	defer mgr.Unlock()
	for k, v := range mgr.conns {
		log.Printf("connmgr clear and stop conn [%d]", k)
		delete(mgr.conns, k)
		v.Stop()
	}
}
