[![Build Status](https://travis-ci.org/hqpko/hbuffer.svg?branch=master)](https://travis-ci.org/hqpko/hbuffer)

# Buffer
go Buffer.
使用方便、可复用

```
package main

import (
	"fmt"

	"github.com/hqpko/hbuffer"
)

func main() {
	by := hbuffer.NewBuffer()

	by.WriteBool(true)
	by.WriteInt32(123)
	by.WriteInt64(124)
	by.WriteUint32(125)
	by.WriteUint64(126)
	by.WriteFloat32(122.33)
	by.WriteFloat64(122.44)
	by.WriteString("test_abc一二三")

	//read from position 0
	by.SetPosition(0)

	boo, err := by.ReadBool()
	fmt.Println(boo, err) //true,nil

	i32, err := by.ReadInt32()
	fmt.Println(i32, err) //123,nil

	i64, err := by.ReadInt64()
	fmt.Println(i64, err) //124,nil

	ui32, err := by.ReadUint32()
	fmt.Println(ui32, err) //125,nil

	ui64, err := by.ReadUint64()
	fmt.Println(ui64, err) //126,nil

	f32, err := by.ReadFloat32()
	fmt.Println(f32, err) //122.33,nil

	f64, err := by.ReadFloat64()
	fmt.Println(f64, err) //122.44,nil

	s, err := by.ReadString()
	fmt.Println(s, err) //test_abc一二三,nil
}

```
