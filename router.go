package tcpserver

import (
	"log"
	"time"
)

type Router interface {
	PreHandle(Request) error
	Handle(Request) error
	PostHandle(Request) error
}

type BaseRouter struct {
}

func NewBaseRouter() Router {
	return &BaseRouter{}
}

func (r *BaseRouter) PreHandle(req Request) error {
	return nil
}

// baserouter handle massge 1 ---> msg{101, "thank you for sending me a message"}
func (r *BaseRouter) Handle(req Request) error {
	conn := req.Conn()
	msg := NewMessage(101, []byte("thank you for sending me a message"))
	cnt, err := conn.SendMsg(msg)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("conn [%d] BaseRouter - write %d bytes Msg %v, (data=%q)",
		conn.ConnId(), cnt, msg, string(msg.Data()))
	time.Sleep(time.Second)

	return nil
}
func (r *BaseRouter) PostHandle(req Request) error {
	return nil
}
