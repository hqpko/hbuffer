package hbuffer

import (
	"testing"
)

func TestBuffer(t *testing.T) {
	b := NewBuffer()
	// read&write bool
	b.WriteBool(true)
	b.WriteBool(false)
	b.SetPosition(0)
	boo1, _ := b.ReadBool()
	shouldBeBool(boo1, true, t)
	boo2, _ := b.ReadBool()
	shouldBeBool(boo2, false, t)
	shouldEqual(b.GetPosition(), 2, t)
	shouldEqual(b.Len(), 2, t)

	// read&write int32
	b.Reset()
	b.WriteInt32(133)
	b.SetPosition(0)
	i32, _ := b.ReadInt32()
	shouldEqual(i32, int32(133), t)

	// read&write uint32
	b.WriteUint64(133)
	b.SetPosition(0)
	u32, _ := b.ReadUint32()
	shouldEqual(u32, uint32(133), t)

	// read&write int64
	b.Reset()
	b.WriteInt64(133)
	b.SetPosition(0)
	i64, _ := b.ReadInt64()
	shouldEqual(i64, int64(133), t)

	// read&write uint64
	b.Reset()
	b.WriteUint64(133)
	b.SetPosition(0)
	u64, _ := b.ReadUint32()
	shouldEqual(u64, uint32(133), t)

	// read&write float32
	b.Reset()
	b.WriteFloat32(133.33)
	b.SetPosition(0)
	f32, _ := b.ReadFloat32()
	shouldEqual(f32, float32(133.33), t)
	shouldEqual(b.GetPosition(), 4, t)
	shouldEqual(b.Len(), 4, t)

	// read&write float64
	b.Reset()
	b.WriteFloat64(133.33)
	b.SetPosition(0)
	f64, _ := b.ReadFloat64()
	shouldEqual(f64, float64(133.33), t)
	shouldEqual(b.GetPosition(), 8, t)
	shouldEqual(b.Len(), 8, t)

	// read&write string
	b.Reset()
	src := "test_abc一二三"
	b.WriteString(src)
	b.SetPosition(0)
	s, _ := b.ReadString()
	shouldEqual(src, s, t)

	// io.reader & io.writer
	b.Reset()
	bytes := []byte{0, 1, 2, 3}
	n, e := b.Write(bytes)
	shouldEqual(n, len(bytes), t)
	shouldEqual(e, nil, t)

	b.SetPosition(0)
	newBytes := make([]byte, 2)
	n, e = b.Read(newBytes)
	shouldEqual(n, len(newBytes), t)
	shouldEqual(e, nil, t)

	b.SetPosition(0)
	newBytes = make([]byte, 10)
	n, e = b.Read(newBytes)
	shouldEqual(n, len(bytes), t)
	shouldEqual(e, nil, t)
}

func TestBufferErrs(t *testing.T) {
	b := NewBuffer()
	_, e := b.ReadBytes(1 << 10)
	shouldNotEqual(e, nil, t)
}

func shouldBeBool(b, equal bool, t *testing.T) {
	if b != equal {
		t.Fatal("not be bool")
	}
}

func shouldEqual(a, b interface{}, t *testing.T) {
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
}

func shouldNotEqual(a, b interface{}, t *testing.T) {
	if a == b {
		t.Fatalf("%v == %v", a, b)
	}
}
