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
			boo1, err := b.ReadBool()
			So(err, ShouldBeNil)
			So(boo1, ShouldBeTrue)
			boo2, err := b.ReadBool()
			So(err, ShouldBeNil)
			So(boo2, ShouldBeFalse)
			So(b.GetPosition(), ShouldEqual, 2)
			So(b.Len(), ShouldEqual, 2)
		})
		Convey("Check Write & Read Int32.", func() {
			b.WriteInt32(133)
			b.SetPosition(0)
			i, err := b.ReadInt32()
			So(err, ShouldBeNil)
			So(i, ShouldEqual, 133)
			So(b.GetPosition(), ShouldEqual, 4)
			So(b.Len(), ShouldEqual, 4)
		})
		Convey("Check Write & Read Uint32.", func() {
			b.WriteUint32(133)
			b.SetPosition(0)
			i, err := b.ReadUint32()
			So(err, ShouldBeNil)
			So(i, ShouldEqual, 133)
			So(b.GetPosition(), ShouldEqual, 4)
			So(b.Len(), ShouldEqual, 4)
		})
		Convey("Check Write & Read Uint64.", func() {
			b.WriteUint64(133)
			b.SetPosition(0)
			i, err := b.ReadUint64()
			So(err, ShouldBeNil)
			So(i, ShouldEqual, 133)
			So(b.GetPosition(), ShouldEqual, 8)
			So(b.Len(), ShouldEqual, 8)
		})
		Convey("Check Write & Read Int64.", func() {
			b.WriteInt64(133)
			b.SetPosition(0)
			i, err := b.ReadInt64()
			So(err, ShouldBeNil)
			So(i, ShouldEqual, 133)
			So(b.GetPosition(), ShouldEqual, 8)
			So(b.Len(), ShouldEqual, 8)
		})
		Convey("Check Write & Read Float32.", func() {
			b.WriteFloat32(133.33)
			b.SetPosition(0)
			i, err := b.ReadFloat32()
			So(err, ShouldBeNil)
			So(i, ShouldEqual, 133.33)
			So(b.GetPosition(), ShouldEqual, 4)
			So(b.Len(), ShouldEqual, 4)
		})
		Convey("Check Write & Read Float64.", func() {
			b.WriteFloat64(133.33)
			b.SetPosition(0)
			i, err := b.ReadFloat64()
			So(err, ShouldBeNil)
			So(i, ShouldEqual, 133.33)
			So(b.GetPosition(), ShouldEqual, 8)
			So(b.Len(), ShouldEqual, 8)
		})
		Convey("Check Write & Read String.", func() {
			src := "test_abc一二三"
			b.WriteString(src)
			b.SetPosition(0)
			s, err := b.ReadString()
			So(err, ShouldBeNil)
			So(s, ShouldEqual, src)
			So(b.GetPosition(), ShouldEqual, uint64(len(src)+8))
			So(b.Len(), ShouldEqual, uint64(len(src)+8))
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

			boo, err := b.ReadBool()
			So(err, ShouldBeNil)
			So(boo, ShouldBeTrue)

			i32, err := b.ReadInt32()
			So(err, ShouldBeNil)
			So(i32, ShouldEqual, 123)

			i64, err := b.ReadInt64()
			So(err, ShouldBeNil)
			So(i64, ShouldEqual, 124)

			ui32, err := b.ReadUint32()
			So(err, ShouldBeNil)
			So(ui32, ShouldEqual, 125)

			ui64, err := b.ReadUint64()
			So(err, ShouldBeNil)
			So(ui64, ShouldEqual, 126)

			f32, err := b.ReadFloat32()
			So(err, ShouldBeNil)
			So(f32, ShouldEqual, 122.33)

			f64, err := b.ReadFloat64()
			So(err, ShouldBeNil)
			So(f64, ShouldEqual, 122.44)

			s, err := b.ReadString()
			So(err, ShouldBeNil)
			So(s, ShouldEqual, "test_abc")
		})
	})
}

func TestBufferErrs(t *testing.T) {
	Convey("Check Buffer Funcs.", t, func() {
		b := NewBuffer()
		Convey("Check Read Bool Err.", func() {
			_, err := b.ReadBool()
			So(err, ShouldNotBeNil)
		})
		Convey("Check Read Int32 Err.", func() {
			_, err := b.ReadInt32()
			So(err, ShouldNotBeNil)
		})
		Convey("Check Read Uint32 Err.", func() {
			_, err := b.ReadUint32()
			So(err, ShouldNotBeNil)
		})
		Convey("Check Read Uint64 Err.", func() {
			_, err := b.ReadUint64()
			So(err, ShouldNotBeNil)
		})
		Convey("Check Read Int64 Err.", func() {
			_, err := b.ReadInt64()
			So(err, ShouldNotBeNil)
		})
		Convey("Check Read Float32 Err.", func() {
			_, err := b.ReadFloat32()
			So(err, ShouldNotBeNil)
		})
		Convey("Check Read Float64 Err.", func() {
			_, err := b.ReadFloat64()
			So(err, ShouldNotBeNil)
		})
		Convey("Check Read String Err.", func() {
			_, err := b.ReadString()
			So(err, ShouldNotBeNil)
		})
	})
}
