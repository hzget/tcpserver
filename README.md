# tcpserver

A tcpserver framework for many service.

## structure

* reader and writer split
* a workerpool to consume taskqueue
* connection management
* tcpdata pack/unpack
* message pack/unpack
* register router handlers for different msg type
* register hook funcs on conn start/stop
* add/get properties for specific connection

```golang
                       ------------------------------
   <--- writer G <---  |        WorkerPool           |
                       | worker  worker ... worker   |
                       ------------------------------
                       |   ^       ^    ...   ^      |
                       |   |       |    ...   |      |
                       |   |       |    ...   |      |
   ---> reader G --->  |  taskq  taskq  ...  taskq   |
                       ------------------------------
```
