package tcpserver

type TcpData interface {
	Size() uint32
	Data() []byte
	SetSize(size uint32)
	SetData(data []byte)
}

//  tcpdata:   size + data
//  size is used for fixing the tcp-packing-issue
type tcpdata struct {
	size uint32
	data []byte
}

func (tdata *tcpdata) Size() uint32 {
	return tdata.size
}

func (tdata *tcpdata) Data() []byte {
	return tdata.data
}

func (tdata *tcpdata) SetSize(size uint32) {
	tdata.size = size
}

func (tdata *tcpdata) SetData(data []byte) {
	tdata.data = data
}

type Message interface {
	Id() uint32
	Data() []byte
	SetId(id uint32)
	SetData(data []byte)
}

type message struct {
	id   uint32
	data []byte
}

func (msg *message) Id() uint32 {
	return msg.id
}

func (msg *message) Data() []byte {
	return msg.data
}

func (msg *message) SetId(id uint32) {
	msg.id = id
}

func (msg *message) SetData(data []byte) {
	msg.data = data
}
