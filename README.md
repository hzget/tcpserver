# tcpserver

A tcpserver framework for users to dev backend
services.
Applications that currently use this framework:
[mmo game][mmo game], ...

## how to use it

It can be used to register handlers
for different client request, add hook funcs
on the event that the conn is on/off, and so on.
An example to use it:

```golang
func main() {
	s := tcpserver.NewServer()
	// register handlers for different msg from the client
	s.AddRouter(core.MSG_C_Talk, &core.ChatRouter{})
	s.AddRouter(core.MSG_C_Move, &core.MoveRouter{})
	// hook func when the player is online
	s.SetOnConnStart(OnConnStart)
	// hook func when the player is offline
	s.SetOnConnStop(OnConnStop)
	s.Serve()
}
```

AddRouter() works just like Handle() and HandleFunc() in net/http package.

```golang
// in net/http package
func Handle(pattern string, handler Handler)
func HandleFunc(pattern string, handler func(ResponseWriter, *Request))

// in this package
func (mhr *msghandler) AddRouter(msgID uint32, r Router)
```

The msg types and corresponding routers are defined
by the user. Router is an interface that shall be implemented
by the user. It contains three methods:

* PreHandle(Request) error
* Handle(Request) error
* PostHandle(Request) error

## features

* reader and writer split
* a workerpool to consume taskqueue
* connection management
* tcpdata pack/unpack
* message pack/unpack
* register router handlers for different msg type
* register hook funcs on conn start/stop
* add/get properties for specific connection

```golang

protocol:

    --------------------------------------
    |            tcp package             |
    --------------------------------------
    |   size    |   msgid   |  msg data  |
    | (4 bytes) | (4 bytes) |            |
    --------------------------------------
                | <----  raw data  ----> |

workerpool:

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
    taskq:  |conn&msg| --->  |conn&msg| --->  |conn&msg| --->  worker
            ----------       ----------       ----------
               req              req              req

```

[mmo game]: https://github.com/hzget/mmo
