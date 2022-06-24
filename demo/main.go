package main

import (
	"log"
	"tcpserver"
	"time"
)

type pingRouter struct {
	tcpserver.BaseRouter
}

func (p *pingRouter) Handle(req tcpserver.Request) error {
	conn := req.Conn()
	msg := tcpserver.NewMessage(102, ([]byte("ping...pong...ping...pong...")))
	cnt, err := conn.SendMsg(msg)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("conn [%d] PingRouter - write %d bytes Msg %v, (data=%q)",
		conn.ConnId(), cnt, msg, string(msg.Data()))
	time.Sleep(time.Second)

	return nil
}

func main() {
	s := tcpserver.NewServer()
	s.AddRouter(2, &pingRouter{})
	s.Serve()
}
