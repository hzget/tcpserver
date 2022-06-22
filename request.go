package tcpserver

import (
	"log"
	"time"
)

type Request struct {
	conn Conn
	data []byte
}

func NewRequest(conn *Connection, data []byte) *Request {
	return &Request{conn, data}
}

func (r *Request) Conn() Conn {
	return r.conn
}

func (r *Request) Data() []byte {
	return r.data
}

type Router interface {
	PreHandle(*Request) error
	Handle(*Request) error
	PostHandle(*Request) error
}

type BaseRouter struct {
}

func (r *BaseRouter) PreHandle(request *Request) error {
	log.Println("BaseRouter - Prehandle")
	return nil
}

func (r *BaseRouter) Handle(request *Request) error {
	conn := request.Conn()
	data := request.Data()
	cnt, err := conn.TCPConn().Write(data)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("BaseRouter - conn [%d] write %d bytes %v: %q",
		conn.ConnId(), cnt, data[:cnt], string(data[:cnt]))
	time.Sleep(time.Second)

	return nil
}
func (r *BaseRouter) PostHandle(request *Request) error {
	log.Println("BaseRouter - Posthandle")
	return nil
}
