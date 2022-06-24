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

	if err := v.PreHandle(req); err != nil {
		return fmt.Errorf("conn [%d] PreHandle Msg [%d] failed: %v", req.Conn().ConnId(), id, err)
	}
	if err := v.Handle(req); err != nil {
		return fmt.Errorf("conn [%d] Handle Msg [%d] failed: %v", req.Conn().ConnId(), id, err)
	}
	if err := v.PostHandle(req); err != nil {
		return fmt.Errorf("conn [%d] PostHandle Msg [%d] failed: %v", req.Conn().ConnId(), id, err)
	}

	return nil
}

func (mh *msghandler) AddRouter(msgID uint32, r Router) {
	mh.handlers[msgID] = r
}
