package network

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
)

type Buffer struct {
	data []byte
	pos  int
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		data: make([]byte, size),
		pos:  0,
	}
}

func (b *Buffer) Bytes() []byte {
	return b.data[:b.pos]
}

func (b *Buffer) WriteByte(value byte) error {
	if b.pos >= len(b.data) {
		return fmt.Errorf("WriteByte buffer overflow")
	}
	b.data[b.pos] = value
	b.pos++
	return nil
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.pos >= len(b.data) {
		return 0, errors.New("buffer underflow: no byte available")
	}
	value := b.data[b.pos]
	b.pos++
	return value, nil
}

func (b *Buffer) WriteFShort(value float32) error {
	fixed := int16(value * 32)
	b.WriteShort(fixed)
	return nil
}

func (b *Buffer) ReadFShort() (float32, error) {
	fixed, err := b.ReadShort()
	if err != nil {
		return 0, err
	}
	return float32(fixed) / 32, nil
}

func (b *Buffer) WriteSByte(value int8) error {
	return b.WriteByte(byte(value))
}

func (b *Buffer) ReadSByte() (int8, error) {
	v, err := b.ReadByte()
	return int8(v), err
}

func (b *Buffer) WriteShort(value int16) {
	if b.pos+2 > len(b.data) {
		panic("buffer overflow")
	}
	binary.BigEndian.PutUint16(b.data[b.pos:], uint16(value))
	b.pos += 2
}

func (b *Buffer) ReadShort() (int16, error) {
	if b.pos+2 > len(b.data) {
		return 0, errors.New("buffer underflow: not enough bytes for short")
	}
	value := int16(binary.BigEndian.Uint16(b.data[b.pos:]))
	b.pos += 2
	return value, nil
}

func (b *Buffer) WriteString(value string) error {
	bytes := []byte(value)
	if len(bytes) > 64 {
		bytes = bytes[:64]
	}
	for i := 0; i < len(bytes); i++ {
		err := b.WriteByte(bytes[i])
		if err != nil {
			return err
		}
	}
	for i := len(bytes); i < 64; i++ {
		err := b.WriteByte(0x20)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Buffer) ReadString() (string, error) {
	if b.pos+64 > len(b.data) {
		return "", errors.New("buffer underflow: not enough bytes for string")
	}
	strBytes := b.data[b.pos : b.pos+64]
	b.pos += 64
	return strings.TrimRight(string(strBytes), " "), nil
}

func (b *Buffer) WriteByteArray(data []byte) error {
	if len(data) > 1024 {
		data = data[:1024]
	}
	for i := 0; i < len(data); i++ {
		err := b.WriteByte(data[i])
		if err != nil {
			return err
		}
	}
	for i := len(data); i < 1024; i++ {
		err := b.WriteByte(0x00)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Buffer) ReadByteArray() ([]byte, error) {
	if b.pos+1024 > len(b.data) {
		return nil, errors.New("buffer underflow: not enough bytes for byte array")
	}
	data := b.data[b.pos : b.pos+1024]
	b.pos += 1024
	return data, nil
}

func (b *Buffer) Data() []byte {
	return b.data
}

func (b *Buffer) Reset() {
	b.pos = 0
}

func (b *Buffer) ResetWithData(data []byte) {
	b.data = data
	b.pos = 0
}
