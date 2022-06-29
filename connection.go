package tcpserver

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	MaxPackSize = 256
	MaxConn     = 1000
	ReadTimeout = 3 * time.Second
)

type Conn interface {
	Start()
	Stop()
	TCPConn() *net.TCPConn
	ConnId() uint32
	RemoteAddr() net.Addr
	Msghandler() MsgHandler
	SendMsg(msg Message) error
	AddProperty(key string, value interface{})
	RemoveProperty(key string)
	GetProperty(key string) (interface{}, error)
}

type Connection struct {
	conn     *net.TCPConn
	id       uint32
	isClosed bool
	handler  MsgHandler
	msgch    chan Message
	wch      chan struct{}
	mutex    sync.Mutex
	property map[string]interface{}
	pmutex   sync.RWMutex
}

func NewConnection(conn *net.TCPConn, id uint32, mhr MsgHandler) Conn {
	return &Connection{
		conn:     conn,
		id:       id,
		isClosed: false,
		handler:  mhr,
		msgch:    make(chan Message),
		wch:      make(chan struct{}),
		property: make(map[string]interface{}),
	}
}

func (c *Connection) startReader() {
	p := &packer{}
	for {
		tdata, err := p.UnPackTcp(c.conn)
		if err != nil {
			log.Printf("conn [%d] unpacktcp failed %v", c.id, err)
			break
		}
		log.Printf("conn [%d] read %d bytes %v", c.id, tdata.Size(), tdata.Data())

		msg, err := p.UnPackMessage(tdata.Data())
		if err != nil {
			log.Printf("conn [%d] unpackmessage failed %v", c.id, err)
			break
		}
		log.Printf("conn [%d] msg %v (data=%q)", c.id, msg, string(msg.Data()))

		req := NewRequest(c, msg)

		if config.app.workerpoolsize > 0 {
			// use workerpool to avoid goroutine switch
			// especially when there're millions of goroutines handling the request
			workers.EnqueueTask(req)
		} else {
			if err := c.Msghandler().Handle(req); err != nil {
				log.Println(err)
				break
			}
		}
	}
}

// startWriter: the only method running in the writer goroutine
func (c *Connection) startWriter() {
	for {
		msg, ok := <-c.msgch
		if !ok {
			log.Printf("conn [%d] msg channel is closed and empty", c.ConnId())
			c.wch <- struct{}{}
			return
		}

		cnt, err := c.writeMsg(msg)
		if err != nil {
			log.Printf("conn [%d] writer send msg failed %v", c.ConnId(), err)
		} else {
			log.Printf("conn [%d] writer send %d bytes msg %v", c.ConnId(), cnt, msg)
		}
	}
}

func (c *Connection) Start() {
	log.Printf("conn [%d] start %s", c.ConnId(), c.RemoteAddr().String())
	hooks.onconnstart(c)
	defer c.Stop()

	go c.startWriter()
	c.startReader()
	log.Printf("conn [%d] remove from connmgr after startReader is done", c.ConnId())
	connmgr.Remove(c)
}

func (c *Connection) Stop() {

	c.mutex.Lock()
	defer c.mutex.Unlock()
	defer hooks.onconnstop(c)
	if c.isClosed {
		return
	}

	//	connmgr.Remove(c)
	log.Printf("conn [%d] close msg channel and consume remainings", c.ConnId())
	close(c.msgch)
	<-c.wch

	log.Printf("conn [%d] stop connection now: %s", c.ConnId(), c.RemoteAddr().String())
	c.conn.Close()
	c.isClosed = true
}

func (c *Connection) TCPConn() *net.TCPConn {
	return c.conn
}

func (c *Connection) ConnId() uint32 {
	return c.id
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) Msghandler() MsgHandler {
	return c.handler
}

// SendMsg: send msg to msgch that will be
// consumed by the writer goroutine.
//
// It's called by the reader goroutine
func (c *Connection) SendMsg(msg Message) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.isClosed {
		return fmt.Errorf("conn was closed and failed to send msg %v", msg)
	}
	c.msgch <- msg
	return nil
}

// only called by the writer goroutine
func (c *Connection) writeMsg(msg Message) (int, error) {

	p := &packer{}
	data, err := p.Pack(msg)
	if err != nil {
		return 0, err
	}

	cnt, err := c.TCPConn().Write(data)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	return cnt, nil
}

func (c *Connection) AddProperty(key string, value interface{}) {
	c.pmutex.Lock()
	defer c.pmutex.Unlock()
	c.property[key] = value
}

func (c *Connection) RemoveProperty(key string) {
	c.pmutex.Lock()
	defer c.pmutex.Unlock()
	delete(c.property, key)
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.pmutex.RLock()
	defer c.pmutex.RUnlock()
	v, ok := c.property[key]
	if !ok {
		return nil, fmt.Errorf("conn [%d] fail to GetProperty - %v", c.id, key)
	}
	return v, nil
}
