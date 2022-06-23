package tcpserver

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Server struct {
	ln     net.Listener
	router Router
}

func NewServer() *Server {
	return &Server{
		router: nil,
	}
}

func (s *Server) AddRouter(router Router) {
	s.router = router
	log.Println("Add router")
}

func (s *Server) Start() error {
	log.Println("tcpserver is starting...")
	if s.router == nil {
		s.router = &BaseRouter{}
	}
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d",
		config.tcpserver.host, config.tcpserver.port))
	if err != nil {
		log.Println(err)
		return err
	}
	s.ln = ln

	cid := uint32(1)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			// handle error
			continue
		}
		/*
					    v1.0: handle tcpconn with a handler
			            	go handleConnectionInteractive(conn)
						v1.1: bind a tcpconn with a handler
			            	c := NewConnection(conn.(*net.TCPConn), cid, handler)
						v1.2: a request combines a connection and its data
						      and a router is registered to handle the request
		*/
		c := NewConnection(conn.(*net.TCPConn), cid, s.router)
		cid++
		go c.Start()
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

/*
 In the interactive mode (read-write-read-write-...),
 need to read all msg in the kernel read buffer via
 1. for-loop
 2. msg-len flag at the beginning of the msg
 3. stop reading via eof or msg-len

 after that, move on and write the response.

*/

// echo all msg to the client
func handleConnectionInteractive(conn net.Conn) {
	in := make([]byte, ReadBuffSize)
	defer conn.Close()

	remoteaddr := conn.RemoteAddr()
	log.Println("get conn from:", remoteaddr.String())

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
		log.Println("read", cnt, "bytes:", in[:cnt], string(in[:cnt]))

		wcnt, err := conn.Write(in[:cnt])
		if err != nil {
			log.Println(err)
			break
		}
		log.Println("write", wcnt, "bytes:", in[:wcnt], string(in[:wcnt]))

		time.Sleep(1 * time.Second)
	}
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
