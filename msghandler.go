package tcpserver

import (
	"fmt"
	"log"
)

type MsgHandler interface {
	Handle(req Request) error
	AddRouter(msgID uint32, r Router)
	StartWorkerPool()
	TaskQueue() []chan Request
}

type msghandler struct {
	handlers  map[uint32]Router
	taskqueue []chan Request
}

func NewMsgHandler() MsgHandler {
	return &msghandler{
		handlers:  make(map[uint32]Router),
		taskqueue: make([]chan Request, config.app.workerpoolsize),
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

func (mhr *msghandler) TaskQueue() []chan Request {
	return mhr.taskqueue
}

func (mhr *msghandler) StartWorkerPool() {
	log.Printf("start worker pool")
	for i := uint32(0); i < config.app.workerpoolsize; i++ {
		mhr.taskqueue[i] = make(chan Request, config.app.taskqueuesize)
		go func(wid uint32) {
			log.Printf("worker[%d] start", wid)
			for {
				req, ok := <-mhr.taskqueue[wid]
				if !ok {
					log.Printf("worker[%d] is closed", wid)
					return
				}

				if err := mhr.Handle(req); err != nil {
					log.Println(err)
				}
			}
		}(i)
	}
}
