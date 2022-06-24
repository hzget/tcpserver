package tcpserver

type Request interface {
	Conn() Conn
	Msg() Message
}

type request struct {
	conn Conn
	msg  Message
}

func NewRequest(conn *Connection, msg Message) Request {
	return &request{conn, msg}
}

func (r *request) Conn() Conn {
	return r.conn
}

func (r *request) Msg() Message {
	return r.msg
}
