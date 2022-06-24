package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	c, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("conn type is %T\n", c)
	conn := c.(*net.TCPConn)
	defer conn.Close()

	//s := "GET / HTTP/1.0\r\n\r\n"
	b1 := []byte{9, 0, 0, 0, 1, 0, 0, 0, 'h', 'e', 'l', 'l', 'o'}
	b2 := []byte{7, 0, 0, 0, 1, 0, 0, 0, 'w', 'h', 'o'}
	s := append(b1, b2...)

	for {
		log.Printf("send msg to the server: %v\n", s)
		fmt.Fprintf(conn, string(s))
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
		time.Sleep(20*time.Second)
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
