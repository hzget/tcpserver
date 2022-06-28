package tcpserver

import (
	"fmt"
)

type MsgHandler interface {
	Handle(req Request) error
	AddRouter(msgID uint32, r Router)
}

type msghandler struct {
	handlers map[uint32]Router
}

func NewMsgHandler() MsgHandler {
	return &msghandler{
		handlers: make(map[uint32]Router),
	}
}

func (mh *msghandler) Handle(req Request) error {
	id := req.Msg().Id()
	v, ok := mh.handlers[id]
	if !ok {
		return fmt.Errorf("conn [%d] not find the router for msgid %d", req.Conn().ConnId(), id)
	}

	if err := v.Handle(req); err != nil {
		return fmt.Errorf("conn [%d] Handle Msg [%d] failed: %v", req.Conn().ConnId(), id, err)
	}

	return nil
}

func (mhr *msghandler) AddRouter(msgID uint32, r Router) {
	mhr.handlers[msgID] = r
}
