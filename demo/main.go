package main

import "tcpserver"

func main(){
	s := tcpserver.NewServer()
	s.Serve()
}
