package tcpserver

type Router interface {
	PreHandle(Request) error
	Handle(Request) error
	PostHandle(Request) error
}

type BaseRouter struct {
}

func NewBaseRouter() Router {
	return &BaseRouter{}
}

func (r *BaseRouter) PreHandle(req Request) error {
	return nil
}

// baserouter handle massge 1 ---> msg{101, "thank you for sending me a message"}
func (r *BaseRouter) Handle(req Request) error {
	conn := req.Conn()
	// handle request
	msg := NewMessage(101, []byte("thank you for sending me a message"))

	// send response msg to writer goroutine
	// is it okay to send to a func-returned channel?
	conn.MsgChan() <- msg

	return nil
}
func (r *BaseRouter) PostHandle(req Request) error {
	return nil
}
