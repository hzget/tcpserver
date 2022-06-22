package main

import (
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	c, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("conn type is %T\n", c)
	conn := c.(*net.TCPConn)
	defer conn.Close()

	s := "GET / HTTP/1.0\r\n\r\n"

	for {
		log.Printf("send msg to the server: %q\n", s)
		fmt.Fprintf(conn, s)
		if err != nil {
			log.Println(err)
			return
		}

		in := make([]byte, 127)
		cnt, err := conn.Read(in)
		if err == io.EOF {
			log.Println("get EOF")
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		log.Println("read", cnt, "bytes msg:", in[:cnt], string(in[:cnt]))
	}
}

/*
func main() {
	c, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("conn type is %T\n", c)
	conn := c.(*net.TCPConn)
	defer conn.Close()

	s := "GET / HTTP/1.0\r\n\r\n"
	fmt.Fprintf(conn, s)
	log.Printf("send msg to the server: %q\n", s)
	if err != nil {
		log.Println(err)
		return
	}
	conn.CloseWrite()

	log.Println("waiting for reading from server...")
	var msg []byte
	in := make([]byte, 127)
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
		log.Println("read", cnt, "bytes msg:", in[:cnt], string(in[:cnt]))
		msg = append(msg, in[:cnt]...)
	}
	log.Println("read", len(msg), "bytes msg:", string(msg))

}
*/
