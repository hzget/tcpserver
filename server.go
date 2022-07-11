package tcpserver

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
)

type Server struct {
	handler MsgHandler
	ln      net.Listener
	done    chan struct{}
}

func NewServer() *Server {
	s := &Server{
		handler: NewMsgHandler(),
		done:    make(chan struct{}, 2),
	}
	s.AddBasicRouters()
	return s
}

func (s *Server) AddRouter(msgId uint32, router Router) {
	s.handler.AddRouter(msgId, router)
}

func (s *Server) Start() error {
	log.Println("tcpserver is starting...")
	defer func() {
		log.Println("Server.Start() return and send done signal")
		s.done <- struct{}{}
	}()
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
			return err
		}
		/*
					    v1.0: handle tcpconn with a handler
			            	go handleConnectionInteractive(conn)
						v1.1: bind a tcpconn with a handler
			            	c := NewConnection(conn.(*net.TCPConn), cid, handler)
						v1.2: a request combines a connection and its data
						      and a router is registered to handle the request
		*/
		c := NewConnection(conn.(*net.TCPConn), cid, s.handler)
		if err := connmgr.Add(c); err != nil {
			log.Println(err)
			continue
		}
		cid++
		go c.Start()
	}
}

func (s *Server) Serve() error {
	go s.Start()
	go s.GracefullyShutdown()
	_, _ = <-s.done, <-s.done
	log.Println("Serve recv two done signals, close the program")
	return nil
}

func (s *Server) Stop() error {
	log.Printf("tcpserver shutdown...")
	s.ln.Close()
	connmgr.Clear()
	workers.Stop()
	return nil
}

func (s *Server) GracefullyShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)
	s.Stop()
	s.done <- struct{}{}
	log.Println("GracefullyShutdown send done signal")
}

func (s *Server) AddBasicRouters() {
	s.AddRouter(1, NewBaseRouter())
}

func (s *Server) SetOnConnStart(fn func(Conn)) {
	hooks.onconnstart = fn
}

func (s *Server) SetOnConnStop(fn func(Conn)) {
	hooks.onconnstop = fn
}
