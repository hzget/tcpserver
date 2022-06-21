package tcpserver

import (
	"io"
	"log"
	"net"
	"time"
)

const (
	ReadBuffSize = 7
	ReadTimeout  = 3 * time.Second
)

type Server struct {
	ln net.Listener
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println(err)
		return err
	}
	s.ln = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			// handle error
			continue
		}
		go handleConnection(conn)
	}
}

func (s *Server) Serve() error {
	return s.Start()
}

func (s *Server) Shutdown() error {
	if s.ln != nil {
		return s.ln.Close()
	}
	return nil
}

// echo all msg to the client
func handleConnection(conn net.Conn) {
	var msg []byte
	in := make([]byte, ReadBuffSize)
	defer conn.Close()

	remoteaddr := conn.RemoteAddr()
	log.Println("get conn from:", remoteaddr.String())

	// set deadline to avoid blocking on Read()
	err := conn.SetReadDeadline(time.Now().Add(ReadTimeout))
	if err != nil {
		log.Println("handleConnection - SetReadDeadline failed:", err)
		return
	}

	for {
		cnt, err := conn.Read(in)
		if err == io.EOF {
			log.Println("get EOF")
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		log.Println("read", cnt, "bytes:", in, string(in))
		msg = append(msg, in[:cnt]...)
	}

	log.Println("read", len(msg), "bytes msg:", msg, string(msg))

	if _, err := conn.Write(msg); err != nil {
		log.Println(err)
		return
	}
}
