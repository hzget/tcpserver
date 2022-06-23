package tcpserver

import (
	"log"
	"time"
)

type Request struct {
	conn Conn
	msg  Message
}

func NewRequest(conn *Connection, msg Message) *Request {
	return &Request{conn, msg}
}

func (r *Request) Conn() Conn {
	return r.conn
}

func (r *Request) Msg() Message {
	return r.msg
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
	msg := request.Msg()
	cnt, err := conn.SendMsg(msg)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("BaseRouter - conn [%d] write %d bytes Msg %v",
		conn.ConnId(), cnt, msg)
	time.Sleep(time.Second)

	return nil
}
func (r *BaseRouter) PostHandle(request *Request) error {
	log.Println("BaseRouter - Posthandle")
	return nil
}
