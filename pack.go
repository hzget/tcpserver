package tcpserver

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Packer interface {
	PackTcp(tdata TcpData) ([]byte, error)
	UnPackTcp(r io.Reader) (TcpData, error)
	PackMessage(msg Message) ([]byte, error)
	UnPackMessage(r io.Reader) (Message, error)

	Pack(msg Message) ([]byte, error)
	UnPack(r io.Reader) (Message, error)
}

type packer struct{}

func (*packer) PackTcp(tdata TcpData) ([]byte, error) {

	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.LittleEndian, tdata.Size()); err != nil {
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

	var data = make([]byte, size)
	if err := binary.Read(r, binary.LittleEndian, &data); err != nil {
		return nil, err
	}

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

func (*packer) UnPackMessage(r io.Reader) (Message, error) {

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

	return p.UnPackMessage(bytes.NewBuffer(tdata.Data()))
}
