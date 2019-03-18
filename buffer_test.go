package hbuffer

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBuffer(t *testing.T) {
	Convey("Check Buffer Funcs.", t, func() {
		b := NewBuffer()
		Convey("Check Write & Read Bool.", func() {
			b.WriteBool(true)
			b.WriteBool(false)
			b.SetPosition(0)
			boo1 := b.ReadBool()
			So(boo1, ShouldBeTrue)
			boo2 := b.ReadBool()
			So(boo2, ShouldBeFalse)
			So(b.GetPosition(), ShouldEqual, 2)
			So(b.Len(), ShouldEqual, 2)
		})
		Convey("Check Write & Read Int32.", func() {
			b.WriteInt32(133)
			b.SetPosition(0)
			i, _ := b.ReadInt32()
			So(i, ShouldEqual, 133)
			So(b.GetPosition(), ShouldEqual, 2) // varint len
			So(b.Len(), ShouldEqual, 2)
		})
		Convey("Check Write & Read Uint32.", func() {
			b.WriteUint32(133)
			b.SetPosition(0)
			i, _ := b.ReadUint32()
			So(i, ShouldEqual, 133)
			So(b.GetPosition(), ShouldEqual, 2)
			So(b.Len(), ShouldEqual, 2)
		})
		Convey("Check Write & Read Uint64.", func() {
			b.WriteUint64(133)
			b.SetPosition(0)
			i, _ := b.ReadUint64()
			So(i, ShouldEqual, 133)
			So(b.GetPosition(), ShouldEqual, 2)
			So(b.Len(), ShouldEqual, 2)
		})
		Convey("Check Write & Read Int64.", func() {
			b.WriteInt64(133)
			b.SetPosition(0)
			i, _ := b.ReadInt64()
			So(i, ShouldEqual, 133)
			So(b.GetPosition(), ShouldEqual, 2)
			So(b.Len(), ShouldEqual, 2)
		})
		Convey("Check Write & Read Float32.", func() {
			b.WriteFloat32(133.33)
			b.SetPosition(0)
			i, _ := b.ReadFloat32()
			So(i, ShouldEqual, 133.33)
			So(b.GetPosition(), ShouldEqual, 4)
			So(b.Len(), ShouldEqual, 4)
		})
		Convey("Check Write & Read Float64.", func() {
			b.WriteFloat64(133.33)
			b.SetPosition(0)
			i, _ := b.ReadFloat64()
			So(i, ShouldEqual, 133.33)
			So(b.GetPosition(), ShouldEqual, 8)
			So(b.Len(), ShouldEqual, 8)
		})
		Convey("Check Write & Read String.", func() {
			src := "test_abc一二三"
			b.WriteString(src)
			b.SetPosition(0)
			s, _ := b.ReadString()
			So(s, ShouldEqual, src)
			So(b.GetPosition(), ShouldEqual, len(src)+1)
			So(b.Len(), ShouldEqual, len(src)+1)
		})
		Convey("Check Write & Get Bytes.", func() {
			bytes := []byte{1, 2, 3, 4, 5}
			b.WriteBytes(bytes)
			bs := b.GetBytes()
			So(len(bs), ShouldEqual, len(bytes))
			for i := range bs {
				So(bs[i], ShouldEqual, bytes[i])
			}
		})
		Convey("Check Write & Read All Types.", func() {
			b.WriteBool(true)
			b.WriteInt32(123)
			b.WriteInt64(124)
			b.WriteUint32(125)
			b.WriteUint64(126)
			b.WriteFloat32(122.33)
			b.WriteFloat64(122.44)
			b.WriteString("test_abc")
			b.SetPosition(0)

			boo := b.ReadBool()

			So(boo, ShouldBeTrue)

			i32, _ := b.ReadInt32()

			So(i32, ShouldEqual, 123)

			i64, _ := b.ReadInt64()

			So(i64, ShouldEqual, 124)

			ui32, _ := b.ReadUint32()

			So(ui32, ShouldEqual, 125)

			ui64, _ := b.ReadUint64()

			So(ui64, ShouldEqual, 126)

			f32, _ := b.ReadFloat32()

			So(f32, ShouldEqual, 122.33)

			f64, _ := b.ReadFloat64()

			So(f64, ShouldEqual, 122.44)

			s, _ := b.ReadString()

			So(s, ShouldEqual, "test_abc")
		})

		Convey("Check io.Writer & io.Reader", func() {
			bytes := []byte{0, 1, 2, 3}
			n, e := b.Write(bytes)
			So(n, ShouldEqual, 4)
			So(e, ShouldBeNil)

			b.SetPosition(0)
			newBytes := make([]byte, 2)
			n, e = b.Read(newBytes)
			So(n, ShouldEqual, len(newBytes))
			So(e, ShouldBeNil)

			b.SetPosition(0)
			newBytes = make([]byte, 10)
			n, e = b.Read(newBytes)
			So(n, ShouldEqual, b.Len())
			So(e, ShouldBeNil)
		})
	})
}

func TestBufferErrs(t *testing.T) {
	Convey("Check Buffer Funcs.", t, func() {
		b := NewBuffer()
		Convey("Check Read Bytes Err.", func() {
			_, e := b.ReadBytes(1 << 10)
			So(e, ShouldNotBeNil)
		})
	})
}

func TestBufferDeleteBefore(t *testing.T) {
	Convey("Check Buffer Funcs.", t, func() {
		b := NewBuffer()
		Convey("Check DeleteBefor.", func() {
			b.WriteBytes([]byte{1, 2, 3})
			b.DeleteBefore(2)
			So(b.Len(), ShouldEqual, 1)
			bt, _ := b.ReadByte()
			So(bt, ShouldEqual, 3)
		})
		Convey("Check DeleteBefore All.", func() {
			b.WriteBytes([]byte{1, 2, 3})
			b.DeleteBefore(5)
			So(b.Len(), ShouldEqual, 0)
		})
	})
}
