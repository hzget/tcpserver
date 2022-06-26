package tcpserver

import (
	"log"
	"net"
	"time"
)

const (
	MaxPackSize = 256
	ReadTimeout = 3 * time.Second
)

type Conn interface {
	Start()
	Stop()
	TCPConn() *net.TCPConn
	ConnId() uint32
	RemoteAddr() net.Addr
	Msghandler() MsgHandler
	MsgChan() chan Message
	SendMsg(msg Message) (int, error)
}

type Connection struct {
	conn     *net.TCPConn
	id       uint32
	isClosed bool
	handler  MsgHandler
	msgch    chan Message
}

func NewConnection(conn *net.TCPConn, id uint32, mhr MsgHandler) Conn {
	return &Connection{
		conn:     conn,
		id:       id,
		isClosed: false,
		handler:  mhr,
		msgch:    make(chan Message),
	}
}

func (c *Connection) enqueueRequest(req Request) {
	wid := c.id % config.app.workerpoolsize
	log.Printf("conn[%d] enqueue req to worker[%d]", c.id, wid)
	c.handler.TaskQueue()[wid] <- req
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
			c.enqueueRequest(req)
		} else {
			if err := c.Msghandler().Handle(req); err != nil {
				log.Println(err)
				break
			}
		}
	}
}

func (c *Connection) startWriter() {
	for {
		msg, ok := <-c.msgch
		if !ok {
			log.Printf("conn [%d] is closed", c.ConnId())
			return
		}

		cnt, err := c.SendMsg(msg)
		if err != nil {
			log.Printf("conn [%d] writer send msg failed %v", c.ConnId(), err)
		} else {
			log.Printf("conn [%d] writer send %d bytes msg %v", c.ConnId(), cnt, msg)
		}
	}
}

func (c *Connection) Start() {
	log.Printf("conn [%d] start %s", c.ConnId(), c.RemoteAddr().String())
	defer c.Stop()

	go c.startWriter()
	c.startReader()
}

func (c *Connection) Stop() {

	// add a mutex lock???
	if c.isClosed {
		return
	}

	log.Printf("conn [%d] stop %s", c.ConnId(), c.RemoteAddr().String())

	close(c.msgch)
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

func (c *Connection) MsgChan() chan Message {
	return c.msgch
}

func (c *Connection) SendMsg(msg Message) (int, error) {

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
