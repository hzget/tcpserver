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
	data := []byte("ping... pong...")
	cnt, err := conn.TCPConn().Write(data)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("PingRouter - conn [%d] write %d bytes %v: %q",
		conn.ConnId(), cnt, data[:cnt], string(data[:cnt]))
	time.Sleep(time.Second)

	return nil
}

func main() {
	s := tcpserver.NewServer()
	s.AddRouter(&pingRouter{})
	s.Serve()
}
