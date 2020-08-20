package hbuffer

import (
	"encoding/binary"
	"errors"
	"io"
	"math"
)

const (
	stepReadLen = 1 << 8
	defInitLen  = 1 << 8
)

type Endian int

const (
	BigEndian Endian = iota
	LittleEndian
)

var errNoAvailableBytes = errors.New("buffer: no available bytes")

type Buffer struct {
	buf      []byte
	endian   binary.ByteOrder
	position int
	length   int
	capacity int
}

func NewBuffer() *Buffer {
	return &Buffer{buf: make([]byte, defInitLen), endian: binary.BigEndian}
}

func NewBufferWithLength(l int) *Buffer {
	return NewBuffer().Grow(l)
}

func NewBufferWithBytes(bs []byte) *Buffer {
	return NewBuffer().SetBytes(bs)
}

// NewBufferWithHead，新建带有 4 位长度 head 的 Buffer，head 用于存放除去 head 的数据长度
func NewBufferWithHead() *Buffer {
	return NewBuffer().WriteEndianUint32(0)
}

// UpdateHead 更新 head，设置 head 为此 buffer 实际携带的数据长度，即 buffer.length-4
// 用于 BufferWithHead
// 注意: updateHead 后，游标指在 head 后的第一位，一般在 updateHead 后不再继续写入数据
func (b *Buffer) UpdateHead() *Buffer {
	return b.SetPosition(0).WriteEndianUint32(uint32(b.length - 4))
}

func (b *Buffer) SetBytes(bs []byte) *Buffer {
	b.buf = bs
	b.length = len(bs)
	b.capacity = b.length
	b.position = 0
	return b
}

func (b *Buffer) SetEndian(e Endian) *Buffer {
	if e == BigEndian {
		b.endian = binary.BigEndian
	} else {
		b.endian = binary.LittleEndian
	}
	return b
}

func (b *Buffer) SetPosition(position int) *Buffer {
	if position > b.length {
		position = b.length
	}
	b.position = position
	return b
}

func (b *Buffer) GetPosition() int {
	return b.position
}

func (b *Buffer) Available() int {
	return b.length - b.position
}

func (b *Buffer) Len() int {
	return b.length
}

func (b *Buffer) Write(bytes []byte) (int, error) {
	b.WriteBytes(bytes)
	return len(bytes), nil
}

func (b *Buffer) WriteByte(bt byte) *Buffer {
	b.willWriteLen(1)
	b.buf[b.position-1] = bt
	return b
}

func (b *Buffer) WriteInt(i int) *Buffer {
	return b.writeVarInt(int64(i))
}

func (b *Buffer) WriteInt32(i int32) *Buffer {
	return b.writeVarInt(int64(i))
}

func (b *Buffer) WriteUint32(i uint32) *Buffer {
	return b.writeUvarInt(uint64(i))
}

func (b *Buffer) WriteEndianUint32(i uint32) *Buffer {
	b.willWriteLen(4)
	b.endian.PutUint32(b.buf[b.position-4:], i)
	return b
}

func (b *Buffer) WriteUint64(i uint64) *Buffer {
	return b.writeUvarInt(i)
}

func (b *Buffer) WriteInt64(i int64) *Buffer {
	return b.writeVarInt(i)
}

func (b *Buffer) WriteFloat32(f float32) *Buffer {
	b.willWriteLen(4)
	b.endian.PutUint32(b.buf[b.position-4:], math.Float32bits(f))
	return b
}

func (b *Buffer) WriteFloat64(f float64) *Buffer {
	b.willWriteLen(8)
	b.endian.PutUint64(b.buf[b.position-8:], math.Float64bits(f))
	return b
}

func (b *Buffer) WriteBytes(bytes []byte) *Buffer {
	l := len(bytes)
	b.willWriteLen(l)
	copy(b.buf[b.position-l:], bytes)
	return b
}

func (b *Buffer) WriteBool(boo bool) *Buffer {
	b.willWriteLen(1)
	if boo {
		b.buf[b.position-1] = 1
	} else {
		b.buf[b.position-1] = 0
	}
	return b
}

func (b *Buffer) WriteString(s string) *Buffer {
	return b.WriteUint32(uint32(len(s))).WriteBytes([]byte(s))
}

func (b *Buffer) willWriteLen(l int) {
	b.Grow(l)
	b.position += l
	if b.length < b.position {
		b.length = b.position
	}
}

func (b *Buffer) growPosition(g int) *Buffer {
	b.position += g
	if b.length < b.position {
		b.length = b.position
	}
	return b
}

func (b *Buffer) readVarInt() (int64, error) {
	return binary.ReadVarint(b)
}

func (b *Buffer) readUvarint() (uint64, error) {
	return binary.ReadUvarint(b)
}

func (b *Buffer) writeVarInt(i int64) *Buffer {
	b.Grow(binary.MaxVarintLen64)
	return b.growPosition(binary.PutVarint(b.buf[b.position:], i))
}

func (b *Buffer) writeUvarInt(i uint64) *Buffer {
	b.Grow(binary.MaxVarintLen64)
	return b.growPosition(binary.PutUvarint(b.buf[b.position:], i))
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.Available() < 1 {
		return 0, errNoAvailableBytes
	}
	c := b.buf[b.position]
	b.position++
	return c, nil
}

func (b *Buffer) ReadBool() (bool, error) {
	bt, err := b.ReadByte()
	return bt == 1, err
}

func (b *Buffer) ReadUint32() (uint32, error) {
	i, e := b.readUvarint()
	return uint32(i), e
}

func (b *Buffer) ReadEndianUint32() (uint32, error) {
	if bt, e := b.ReadBytes(4); e != nil {
		return 0, e
	} else {
		return b.endian.Uint32(bt), nil
	}
}

func (b *Buffer) ReadInt() (int, error) {
	i, e := b.readVarInt()
	return int(i), e
}

func (b *Buffer) ReadInt32() (int32, error) {
	i, e := b.readVarInt()
	return int32(i), e
}

func (b *Buffer) ReadUint64() (uint64, error) {
	return b.readUvarint()
}

func (b *Buffer) ReadInt64() (int64, error) {
	return b.readVarInt()
}

func (b *Buffer) ReadFloat32() (float32, error) {
	if bt, e := b.ReadBytes(4); e != nil {
		return 0, e
	} else {
		return math.Float32frombits(b.endian.Uint32(bt)), nil
	}
}

func (b *Buffer) ReadFloat64() (float64, error) {
	if bt, e := b.ReadBytes(8); e != nil {
		return 0, e
	} else {
		return math.Float64frombits(b.endian.Uint64(bt)), nil
	}
}

func (b *Buffer) ReadString() (string, error) {
	if sz, e := b.ReadUint32(); e != nil {
		return "", e
	} else if sb, e := b.ReadBytes(int(sz)); e != nil {
		return "", e
	} else {
		return string(sb), nil
	}
}

// ReadBytes read only bytes
func (b *Buffer) ReadBytes(size int) ([]byte, error) {
	if b.Available() < size {
		return nil, errNoAvailableBytes
	}
	b.position += size
	return b.buf[b.position-size : b.position], nil
}

func (b *Buffer) ReadBytesAtPosition(position, size int) ([]byte, error) {
	p := b.position
	b.SetPosition(position)
	bs, e := b.ReadBytes(size)
	if e != nil {
		return nil, e
	}
	b.SetPosition(p)
	return bs, nil
}

func (b *Buffer) Read(bytes []byte) (int, error) {
	size := len(bytes)
	available := b.Available()
	if size > available {
		size = available
	}
	bt, _ := b.ReadBytes(size)
	copy(bytes, bt)
	return size, nil
}

func (b *Buffer) ReadAll(r io.Reader) error {
	for {
		_, e := b.ReadFromReader(r)
		if e == io.EOF {
			break
		}
		if e != nil {
			return e
		}
	}
	return nil
}

func (b *Buffer) ReadFromReader(r io.Reader) (int, error) {
	b.Grow(stepReadLen)
	n, e := r.Read(b.buf[b.position:])
	if e != nil {
		return 0, e
	}
	if b.length < b.position+n {
		b.length = b.position + n
	}
	return n, nil
}

func (b *Buffer) ReadFull(r io.Reader, l int) (int, error) {
	b.Grow(l)
	n, e := io.ReadFull(r, b.buf[b.position:b.position+l])
	if e != nil {
		return n, e
	}
	if b.length < b.position+n {
		b.length = b.position + n
	}
	return n, e
}

func (b *Buffer) GetBytes() []byte {
	return b.buf[0:b.length]
}

func (b *Buffer) CopyBytes() []byte {
	bs := make([]byte, b.length)
	copy(bs, b.buf[:b.length])
	return bs
}

func (b *Buffer) GetRestOfBytes() []byte {
	return b.buf[b.position:b.length]
}

func (b *Buffer) CopyRestOfBytes() []byte {
	bs := make([]byte, b.length-b.position)
	copy(bs, b.buf[b.position:b.length])
	return bs
}

func (b *Buffer) Reset() *Buffer {
	b.position = 0
	b.length = 0
	return b
}

func (b *Buffer) Back(position int) *Buffer {
	b.position = position
	b.length = position
	return b
}

func (b *Buffer) DeleteBefore(position int) *Buffer {
	if position >= b.length { // delete all
		b.Reset()
	} else {
		copy(b.buf, b.buf[position:])
		b.length = b.length - position
		b.position = 0
	}
	return b
}

func (b *Buffer) Grow(n int) *Buffer {
	need := b.length + n
	capBuf := cap(b.buf)
	if need > capBuf {
		newCap := 2 * capBuf
		for newCap < need {
			newCap *= 2
		}
		buf := make([]byte, newCap)
		copy(buf, b.buf)
		b.buf = buf
	}
	return b
}

func (b *Buffer) Cap() int {
	return cap(b.buf)
}

func WriteByte(bt byte) *Buffer {
	return NewBuffer().WriteByte(bt)
}

func WriteInt(i int) *Buffer {
	return NewBuffer().WriteInt(i)
}

func WriteInt32(i int32) *Buffer {
	return NewBuffer().WriteInt32(i)
}

func WriteUint32(i uint32) *Buffer {
	return NewBuffer().WriteUint32(i)
}

func WriteEndianUint32(i uint32) *Buffer {
	return NewBuffer().WriteEndianUint32(i)
}

func WriteUint64(i uint64) *Buffer {
	return NewBuffer().WriteUint64(i)
}

func WriteInt64(i int64) *Buffer {
	return NewBuffer().WriteInt64(i)
}

func WriteFloat32(f float32) *Buffer {
	return NewBuffer().WriteFloat32(f)
}

func WriteFloat64(f float64) *Buffer {
	return NewBuffer().WriteFloat64(f)
}

func WriteBytes(bytes []byte) *Buffer {
	return NewBuffer().WriteBytes(bytes)
}

func WriteBool(boo bool) *Buffer {
	return NewBuffer().WriteBool(boo)
}

func WriteString(s string) *Buffer {
	return NewBuffer().WriteString(s)
}
