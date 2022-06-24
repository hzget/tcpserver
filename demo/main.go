package main

import (
	"log"
	"tcpserver"
	"time"
)

type pingRouter struct {
	tcpserver.BaseRouter
}

func (p *pingRouter) Handle(request *tcpserver.Request) error {
	conn := request.Conn()
	msg := request.Msg()
	msg.SetData([]byte("ping...pong...ping...pong..."))
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
	//	s.AddRouter(&pingRouter{})
	s.Serve()
}
