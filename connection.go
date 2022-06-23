package tcpserver

import (
	"io"
	"log"
	"net"
	"time"
)

const (
	ReadBuffSize = 127
	ReadTimeout  = 3 * time.Second
)

type Conn interface {
	Start()
	Stop()
	TCPConn() *net.TCPConn
	ConnId() uint32
	RemoteAddr() net.Addr
	Router() Router
}

type Connection struct {
	conn     *net.TCPConn
	id       uint32
	isClosed bool
	router   Router
}

func NewConnection(conn *net.TCPConn, id uint32, router Router) Conn {
	return &Connection{
		conn:     conn,
		id:       id,
		isClosed: false,
		router:   router,
	}
}

func (c *Connection) startReader() {
	in := make([]byte, config.tcpserver.readbuffsize)

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

		req := NewRequest(c, in[:cnt])

		if err := c.router.PreHandle(req); err != nil {
			log.Printf("conn [%d] PreHandle failed %v", c.id, err)
		}
		if err := c.router.Handle(req); err != nil {
			log.Printf("conn [%d] Handle failed %v", c.id, err)
			break
		}
		if err := c.router.PostHandle(req); err != nil {
			log.Printf("conn [%d] PostHandle failed %v", c.id, err)
		}

	}
}

func (c *Connection) Start() {
	log.Printf("conn [%d] start %s", c.ConnId(), c.RemoteAddr().String())
	defer c.Stop()

	c.startReader()
}

func (c *Connection) Stop() {

	// add a mutex lock???
	if c.isClosed {
		return
	}

	log.Printf("conn [%d] stop %s", c.ConnId(), c.RemoteAddr().String())

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

func (c *Connection) Router() Router {
	return c.router
}
