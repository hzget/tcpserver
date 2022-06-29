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
		conn.AddProperty("name", "smart")
		log.Printf("conn [%d] %s is online now", conn.ConnId(), "smart")
	})
	s.SetOnConnStop(func(conn tcpserver.Conn) {
		log.Printf("conn [%d] OnConnStop hookfunc", conn.ConnId())
		name, err := conn.GetProperty("name")
		if err != nil {
			log.Println(err)
			name = "Someone"
		}
		log.Printf("conn [%d] %s is offline now", conn.ConnId(), name)
	})
	s.Serve()
}
