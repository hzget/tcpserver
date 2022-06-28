package main

import (
	"log"
	"tcpserver"
)

type pingRouter struct {
	tcpserver.BaseRouter
}

func (p *pingRouter) Handle(req tcpserver.Request) error {
	conn := req.Conn()
	msg := tcpserver.NewMessage(102, ([]byte("ping...pong...ping...pong...")))
	conn.SendMsg(msg)
	return nil
}

func main() {
	s := tcpserver.NewServer()
	s.AddRouter(2, &pingRouter{})
	s.SetOnConnStart(func(conn tcpserver.Conn) {
		log.Printf("conn [%d] OnConnStart hookfunc", conn.ConnId())
	})
	s.SetOnConnStop(func(conn tcpserver.Conn) {
		log.Printf("conn [%d] OnConnStop hookfunc", conn.ConnId())
	})
	s.Serve()
}
