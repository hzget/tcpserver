package tcpserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type Packer interface {
	PackTcp(tdata TcpData) ([]byte, error)
	UnPackTcp(r io.Reader) (TcpData, error)
	PackMessage(msg Message) ([]byte, error)
	UnPackMessage(rawData []byte) (Message, error)

	Pack(msg Message) ([]byte, error)
	UnPack(r io.Reader) (Message, error)
}

type packer struct{}

func NewPacker() Packer {
	return &packer{}
}

/*
	protocol:

		size means one of the following:
			* size of msgid + msg data (default)
			* size of msg data

		controlled by config.protocol.tcpsizeadjust

	--------------------------------------
	|            tcp package             |
	--------------------------------------
	|   size    |   msgid   |  msg data  |
	| (4 bytes) | (4 bytes) |            |
	--------------------------------------
	            | <----  raw data  ----> |

*/

func (*packer) PackTcp(tdata TcpData) ([]byte, error) {

	if !tdata.IsValid() {
		return nil, fmt.Errorf("fail to pack invalid tcp data %v", tdata)
	}

	size := tdata.Size() - config.protocol.tcpsizeadjust
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.LittleEndian, size); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, tdata.Data()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (*packer) UnPackTcp(r io.Reader) (TcpData, error) {

	var size uint32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return nil, err
	}

	if size > config.tcpserver.maxpacksize {
		return nil, fmt.Errorf("packsize %d is too big", size)
	}

	log.Printf("unpacktcp bytes - size - %d", size)
	size += config.protocol.tcpsizeadjust
	var data = make([]byte, size)
	if err := binary.Read(r, binary.LittleEndian, &data); err != nil {
		return nil, err
	}

	log.Printf("unpacktcp bytes - data - %v", data)
	return &tcpdata{size, data}, nil
}

func (*packer) PackMessage(msg Message) ([]byte, error) {

	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.LittleEndian, msg.Id()); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.LittleEndian, msg.Data()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (*packer) UnPackMessage(rawData []byte) (Message, error) {

	r := bytes.NewBuffer(rawData)
	var id uint32
	if err := binary.Read(r, binary.LittleEndian, &id); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return &message{id, data}, nil
}

func (p *packer) Pack(msg Message) ([]byte, error) {

	mpack, err := p.PackMessage(msg)
	if err != nil {
		return nil, err
	}

	// attention: will int overflows uint32?
	tdata := &tcpdata{uint32(len(mpack)), mpack}
	return p.PackTcp(tdata)
}

func (p *packer) UnPack(r io.Reader) (Message, error) {
	tdata, err := p.UnPackTcp(r)
	if err != nil {
		return nil, err
	}

	return p.UnPackMessage(tdata.Data())
}
