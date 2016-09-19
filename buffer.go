package hbuffer

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

var ErrNotEnoughLength = errors.New("byte array not enough length.")

type Endian int

const (
	BigEndian Endian = iota
	LittleEndian
)

type Buffer struct {
	buf      []byte
	endian   binary.ByteOrder
	position uint64
	length   uint64
	capacity uint64
}

func NewBuffer() *Buffer {
	return &Buffer{buf: []byte{}, endian: binary.BigEndian}
}

func NewBufferWithLength(l uint64) *Buffer {
	b := NewBuffer()
	b.grow(l)
	return b
}

func NewBufferWithBytes(bs []byte) *Buffer {
	b := NewBuffer()
	b.SetBytes(bs)
	return b
}

func (b *Buffer) SetBytes(bs []byte) {
	b.buf = bs
	b.length = uint64(len(bs))
	b.capacity = b.length
	b.position = 0
}

func (b *Buffer) SetEndian(e Endian) {
	if e == BigEndian {
		b.endian = binary.BigEndian
	} else {
		b.endian = binary.LittleEndian
	}
}

func (b *Buffer) SetPosition(position uint64) {
	if position > b.length {
		position = b.length
	}
	b.position = position
}

func (b *Buffer) GetPosition() uint64 {
	return b.position
}

func (b *Buffer) Available() uint64 {
	return b.length - b.position
}

func (b *Buffer) Len() uint64 {
	return b.length
}

func (b *Buffer) WriteByte(bt byte) {
	b.grow(1)
	b.buf[b.position] = bt
	b.position += 1
}

func (b *Buffer) WriteShort(i int) {
	b.grow(2)
	b.endian.PutUint16(b.buf[b.position:b.position+2], uint16(i))
	b.position += 2
}

func (b *Buffer) WriteInt32(i int32) {
	b.grow(4)
	b.endian.PutUint32(b.buf[b.position:b.position+4], uint32(i))
	b.position += 4
}

func (b *Buffer) WriteUint32(i uint32) {
	b.grow(4)
	b.endian.PutUint32(b.buf[b.position:b.position+4], i)
	b.position += 4
}

func (b *Buffer) WriteUint64(i uint64) {
	b.grow(8)
	b.endian.PutUint64(b.buf[b.position:b.position+8], i)
	b.position += 8
}

func (b *Buffer) WriteInt64(i int64) {
	b.grow(8)
	b.endian.PutUint64(b.buf[b.position:b.position+8], uint64(i))
	b.position += 8
}

func (b *Buffer) WriteFloat32(f float32) {
	b.grow(4)
	b.endian.PutUint32(b.buf[b.position:b.position+4], math.Float32bits(f))
	b.position += 4
}

func (b *Buffer) WriteFloat64(f float64) {
	b.grow(8)
	b.endian.PutUint64(b.buf[b.position:b.position+8], math.Float64bits(f))
	b.position += 8
}

func (b *Buffer) WriteBytes(bytes []byte) {
	l := uint64(len(bytes))
	b.grow(l)
	copy(b.buf[b.position:], bytes)
	b.position += l
}

func (b *Buffer) WriteBool(boo bool) {
	if b.Available() < 1 {
		b.grow(1)
	}
	if boo {
		b.buf[b.position] = 1
	} else {
		b.buf[b.position] = 0
	}
	b.position++
}

func (b *Buffer) WriteString(s string) {
	l := uint64(len(s))
	b.WriteUint64(l)
	b.grow(l)
	copy(b.buf[b.position:], s)
	b.position += l
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.Available() == 0 {
		return 0, ErrNotEnoughLength
	}
	c := b.buf[b.position]
	b.position++
	return c, nil
}

func (b *Buffer) ReadShort() (int, error) {
	bs, err := b.ReadBytes(2)
	if err != nil {
		return 0, err
	}
	i := b.endian.Uint16(bs)
	return int(i), nil
}

func (b *Buffer) ReadBool() (bool, error) {
	if b.Available() == 0 {
		return false, ErrNotEnoughLength
	}
	c := b.buf[b.position]
	b.position++
	return c == 1, nil
}

func (b *Buffer) ReadUint32() (uint32, error) {
	bs, err := b.ReadBytes(4)
	if err != nil {
		return 0, err
	}
	i := b.endian.Uint32(bs)
	return i, nil
}

func (b *Buffer) ReadInt32() (int32, error) {
	i, err := b.ReadUint32()
	if err != nil {
		return 0, err
	}
	return int32(i), err
}

func (b *Buffer) ReadUint64() (uint64, error) {
	bs, err := b.ReadBytes(8)
	if err != nil {
		return 0, err
	}
	i := b.endian.Uint64(bs)
	return i, nil
}

func (b *Buffer) ReadInt64() (int64, error) {
	i, err := b.ReadUint64()
	if err != nil {
		return 0, err
	}
	return int64(i), nil
}

func (b *Buffer) ReadFloat32() (float32, error) {
	i, err := b.ReadUint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(i), nil
}

func (b *Buffer) ReadFloat64() (float64, error) {
	i, err := b.ReadUint64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(i), nil
}

func (b *Buffer) ReadString() (string, error) {
	l, err := b.ReadUint64()
	if err != nil {
		return "", err
	}
	bs, err := b.ReadBytes(l)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

//ReadBytes read only bytes
func (b *Buffer) ReadBytes(size uint64) ([]byte, error) {
	if b.Available() < size {
		return nil, ErrNotEnoughLength
	}
	b.position += size
	return b.buf[b.position-size : b.position], nil
}

func (b *Buffer) ReadBytesAtPosition(position uint64, size uint64) ([]byte, error) {
	p := b.position
	b.SetPosition(position)
	bs, err := b.ReadBytes(size)
	b.SetPosition(p)
	return bs, err
}

func (b *Buffer) ReadAll(r io.Reader) error {
	p := b.position
	for {
		b.grow(512)
		m, e := r.Read(b.buf[b.position:cap(b.buf)])
		b.position += uint64(m)
		b.length = b.position
		if e == io.EOF {
			break
		}
		if e != nil {
			return e
		}
	}
	b.position = p
	return nil
}

//CopyBytes copy bytes to new slice
func (b *Buffer) CopyBytes(size uint64) ([]byte, error) {
	if b.Available() < size {
		return nil, ErrNotEnoughLength
	}
	bs := make([]byte, size)
	copy(bs, b.buf[b.position:b.position+size])
	b.position += size
	return bs, nil
}

func (b *Buffer) GetBytes() []byte {
	return b.buf[0:b.length]
}

func (b *Buffer) Reset() {
	b.position = 0
	b.length = 0
}

func (b *Buffer) Back(position uint64) {
	b.position = position
	b.length = position
}

func (b *Buffer) grow(n uint64) uint64 {
	if b.length+n > uint64(cap(b.buf)) {
		buf := make([]byte, (2*cap(b.buf) + int(n)))
		copy(buf, b.buf)
		b.buf = buf
	}
	b.length += n
	return b.length
}
