package tcpserver

var workers WorkerPool
var connmgr ConnManager

func init() {
	connmgr = NewConnManager()
	workers = NewWorkerPool()
}
