package protocol

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
)

type Packet struct {
	buff []byte
}

func (this *Packet) Serialize() []byte {
	return this.buff
}

func (this *Packet) GetLength() uint16 {
	return binary.BigEndian.Uint16(this.buff[0:2])
}

func (this *Packet) GetMsgId() uint16 {
	return binary.BigEndian.Uint16(this.buff[2:4])
}

func (this *Packet) GetBody() []byte {
	return this.buff[4:]
}

func NewPacket(buff []byte, msgId uint16, hasLengthField bool) Packet {
	p := Packet{}

	if hasLengthField {
		p.buff = buff
	} else {
		p.buff = make([]byte, 4+len(buff))

		binary.BigEndian.PutUint16(p.buff[0:2], uint16(len(buff))+4)
		binary.BigEndian.PutUint16(p.buff[2:4], msgId)
		copy(p.buff[4:], buff)
	}
	return p
}

func ReadPacket(conn *net.TCPConn) (*Packet, error) {
	var (
		lengthBytes []byte = make([]byte, 2)
		msgIdBytes  []byte = make([]byte, 2)
		length      uint16
	)

	p := &Packet{}

	// read length
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return p, err
	}
	if length = binary.BigEndian.Uint16(lengthBytes); length > 1024 || length <= 0 {
		return p, errors.New("the size of packet is larger than the limit or is empty")
	}
	//read msgId
	if _, err := io.ReadFull(conn, msgIdBytes); err != nil {
		return p, err
	}

	buff := make([]byte, length)
	copy(buff[0:2], lengthBytes)
	copy(buff[2:4], msgIdBytes)

	// read body ( buff = lengthBytes + body )
	if _, err := io.ReadFull(conn, buff[4:]); err != nil {
		return p, err
	}
	msgId := binary.BigEndian.Uint16(msgIdBytes)
	pp := NewPacket(buff, msgId, true)
	return &pp, nil
}

func ReadPacketCopy(conn *net.TCPConn) (Packet, error) {
	var (
		lengthBytes []byte = make([]byte, 2)
		msgIdBytes  []byte = make([]byte, 2)
		length      uint16
	)

	p := Packet{}

	// read length
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return p, err
	}
	if length = binary.BigEndian.Uint16(lengthBytes); length > 1024 || length <= 0 {
		return p, errors.New("the size of packet is larger than the limit or is empty")
	}
	//read msgId
	if _, err := io.ReadFull(conn, msgIdBytes); err != nil {
		return p, err
	}

	buff := make([]byte, length)
	copy(buff[0:2], lengthBytes)
	copy(buff[2:4], msgIdBytes)

	// read body ( buff = lengthBytes + body )
	if _, err := io.ReadFull(conn, buff[4:]); err != nil {
		return p, err
	}
	msgId := binary.BigEndian.Uint16(msgIdBytes)
	return NewPacket(buff, msgId, true), nil
}
