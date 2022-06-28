package tcpserver

var workers WorkerPool
var connmgr ConnManager
var hooks *hook

func init() {
	connmgr = NewConnManager()
	workers = NewWorkerPool()
	hooks = NewHook()
}
