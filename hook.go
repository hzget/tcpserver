package tcpserver

import (
	"log"
)

type hook struct {
	onconnstart func(Conn)
	onconnstop  func(Conn)
}

func NewHook() *hook {
	return &hook{
		onconnstart: func(Conn) { log.Printf("Default OnConnStart HookFunc") },
		onconnstop:  func(Conn) { log.Printf("Default OnConnStop HookFunc") },
	}
}
