package tcpserver

import (
	"net"
	"log"
	"io"
)

type Conn interface {
	Start()
	Stop()
	TCPConn() *net.TCPConn
	ConnId() uint32
	RemoteAddr() net.Addr
}

type HandlerFunc func(*net.TCPConn, []byte, int) error

type Connection struct {
	conn *net.TCPConn
	id uint32
	isClosed bool
	handler HandlerFunc
}

func NewConnection(conn *net.TCPConn, id uint32, handler HandlerFunc) *Connection {
	return &Connection{
		conn: conn,
		id: id,
		isClosed: false,
		handler: handler,
	}
}

func (c *Connection) startReader() {
	in := make([]byte, ReadBuffSize)

	for {
		cnt, err := c.conn.Read(in)
		if err == io.EOF {
			log.Printf("conn [%d] get EOF", c.id)
			break
		}
		if err != nil {
			log.Printf("conn [%d] read failed %v", c.id, err)
			break
		}
		log.Printf("conn [%d] read %d bytes %v %q", c.id, cnt, in[:cnt], string(in[:cnt]))

		if err := c.handler(c.conn, in, cnt); err != nil {
	//		log.Printf("conn [%d] handler failed %v", c.id, err)
			break
		}
	}
}

func (c *Connection) Start() {
	log.Printf("conn start [%d] %s", c.ConnId(), c.RemoteAddr().String())
	defer c.Stop()

	c.startReader()
}

func (c *Connection) Stop() {

	// add a mutex lock???
	if c.isClosed {
		return
	}

	log.Printf("conn stop [%d] %s", c.ConnId(), c.RemoteAddr().String())

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
