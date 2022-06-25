package main

import (
	"tcpserver"
)

type pingRouter struct {
	tcpserver.BaseRouter
}

func (p *pingRouter) Handle(req tcpserver.Request) error {
	conn := req.Conn()
	msg := tcpserver.NewMessage(102, ([]byte("ping...pong...ping...pong...")))
	conn.MsgChan() <- msg

	return nil
}

func main() {
	s := tcpserver.NewServer()
	s.AddRouter(2, &pingRouter{})
	s.Serve()
}
