package tcpserver

import (
	"log"
)

/*
                       ------------------------------
   <--- writer G <---  |        WorkerPool           |
                       | worker  worker ... worker   |
                       ------------------------------
                       |   ^       ^    ...   ^      |
                       |   |       |    ...   |      |
                       |   |       |    ...   |      |
   ---> reader G --->  |  taskq  taskq  ...  taskq   |
                       ------------------------------

			----------       ----------       ----------
	taskq:	|conn&msg| --->  |conn&msg| --->  |conn&msg| --->  worker
			----------       ----------       ----------
			    req              req              req

*/
type WorkerPool interface {
	Start()
	Stop()
	EnqueueTask(Request)
}

type workerpool struct {
	taskqueues []chan Request
	done       []chan struct{}
}

func NewWorkerPool() WorkerPool {
	w := &workerpool{
		taskqueues: make([]chan Request, config.app.workerpoolsize),
		done:       make([]chan struct{}, config.app.workerpoolsize),
	}
	w.Start()
	return w
}

func (w *workerpool) EnqueueTask(req Request) {
	cid := req.Conn().ConnId()
	wid := cid % config.app.workerpoolsize
	log.Printf("conn [%d] enqueue req to worker[%d]", cid, wid)
	// shall check if the workerpool is closed
	w.taskqueues[wid] <- req
}

func (w *workerpool) Start() {
	log.Printf("start worker pool")
	for i := uint32(0); i < config.app.workerpoolsize; i++ {
		w.taskqueues[i] = make(chan Request, config.app.taskqueuesize)
		w.done[i] = make(chan struct{})
		go w.startworker(i)
	}
}

func (w *workerpool) startworker(wid uint32) {
	log.Printf("worker[%d] start", wid)
	for {
		select {
		case <-w.done[wid]:
			log.Printf("worker[%d] is closed when done chan is closed", wid)
			return
		case req, ok := <-w.taskqueues[wid]:
			if !ok {
				log.Printf("worker[%d] is closed when taskqueue is closed", wid)
				return
			}

			conn := req.Conn()
			log.Printf("conn [%d] in worker [%d] handling ...", conn.ConnId(), wid)
			if err := conn.Msghandler().Handle(req); err != nil {
				log.Printf("conn [%d] in worker [%d] handle failed %v", conn.ConnId(), wid, err)
			}
		}
	}

}

func (w *workerpool) Stop() {

	log.Printf("stop the workerpool")
	// close workerpool goroutines
	for _, v := range w.done {
		close(v)
	}

	// empty taskqueues
	for _, v := range w.taskqueues {
		close(v)
		for range v {
		}
	}
}
