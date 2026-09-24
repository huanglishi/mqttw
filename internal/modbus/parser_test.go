package modbus

import (
	"math"
	"testing"
)

func point(name, regType, dataType string, addr, qty int32) PointConfig {
	return PointConfig{Name: name, RegisterType: regType, Address: addr, Quantity: qty, DataType: dataType}
}

func TestParseUint16(t *testing.T) {
	// 大端 0x1234
	p := point("u16", RegisterHolding, DataTypeUint16, 0, 1)
	v, err := ParseRegisters([]byte{0x12, 0x34}, p, ByteOrderBig, WordOrderHighFirst)
	if err != nil || v != 0x1234 {
		t.Fatalf("uint16 big-endian: got %v err %v", v, err)
	}
	// 小端字节序：低字节在前
	v, err = ParseRegisters([]byte{0x12, 0x34}, p, ByteOrderLittle, WordOrderHighFirst)
	if err != nil || v != 0x3412 {
		t.Fatalf("uint16 little-endian: got %v err %v", v, err)
	}
}

func TestParseInt16(t *testing.T) {
	p := point("i16", RegisterHolding, DataTypeInt16, 0, 1)
	v, err := ParseRegisters([]byte{0xFF, 0xFE}, p, ByteOrderBig, WordOrderHighFirst)
	if err != nil || v != -2 {
		t.Fatalf("int16: got %v err %v", v, err)
	}
}

func TestParseFloat32WordOrders(t *testing.T) {
	p := point("f", RegisterHolding, DataTypeFloat32, 0, 2)
	// 25.5 的 IEEE754: 0x41CC0000
	raw := []byte{0x41, 0xCC, 0x00, 0x00}
	v, err := ParseRegisters(raw, p, ByteOrderBig, WordOrderHighFirst)
	if err != nil || math.Abs(v-25.5) > 1e-6 {
		t.Fatalf("float32 high_first: got %v err %v", v, err)
	}
	// low_first: 寄存器顺序交换 → 原始字节应是 00 00 41 CC
	raw = []byte{0x00, 0x00, 0x41, 0xCC}
	v, err = ParseRegisters(raw, p, ByteOrderBig, WordOrderLowFirst)
	if err != nil || math.Abs(v-25.5) > 1e-6 {
		t.Fatalf("float32 low_first: got %v err %v", v, err)
	}
	// little 字节序 + high_first: 每个 word 内低字节在前 → 字节 CC 41 00 00
	raw = []byte{0xCC, 0x41, 0x00, 0x00}
	v, err = ParseRegisters(raw, p, ByteOrderLittle, WordOrderHighFirst)
	if err != nil || math.Abs(v-25.5) > 1e-6 {
		t.Fatalf("float32 little byte + high_first: got %v err %v", v, err)
	}
}

func TestParseUint32(t *testing.T) {
	p := point("u32", RegisterInput, DataTypeUint32, 0, 2)
	v, err := ParseRegisters([]byte{0x00, 0x01, 0x00, 0x02}, p, ByteOrderBig, WordOrderHighFirst)
	if err != nil || v != 0x00010002 {
		t.Fatalf("uint32: got %v err %v", v, err)
	}
}

func TestParseBCD(t *testing.T) {
	p := point("bcd", RegisterHolding, DataTypeBCD, 0, 1)
	v, err := ParseRegisters([]byte{0x12, 0x34}, p, ByteOrderBig, WordOrderHighFirst)
	if err != nil || v != 1234 {
		t.Fatalf("bcd 1234: got %v err %v", v, err)
	}
	// 8 位 BCD：两个寄存器
	p2 := point("bcd2", RegisterHolding, DataTypeBCD, 0, 2)
	v, err = ParseRegisters([]byte{0x00, 0x01, 0x23, 0x45}, p2, ByteOrderBig, WordOrderHighFirst)
	if err != nil || v != 12345 {
		t.Fatalf("bcd 8位: got %v err %v", v, err)
	}
	// 非法 BCD
	_, err = ParseRegisters([]byte{0x1A, 0x34}, p, ByteOrderBig, WordOrderHighFirst)
	if err == nil {
		t.Fatal("非法 BCD 应报错")
	}
}

func TestParseCoil(t *testing.T) {
	p := point("c", RegisterCoil, DataTypeBit, 0, 1)
	v, err := ParseRegisters([]byte{0x01}, p, ByteOrderBig, WordOrderHighFirst)
	if err != nil || v != 1 {
		t.Fatalf("coil on: got %v err %v", v, err)
	}
	v, err = ParseRegisters([]byte{0x00}, p, ByteOrderBig, WordOrderHighFirst)
	if err != nil || v != 0 {
		t.Fatalf("coil off: got %v err %v", v, err)
	}
	// 线圈不允许非 bit 类型
	p2 := point("c2", RegisterCoil, DataTypeUint16, 0, 1)
	if _, err := ParseRegisters([]byte{0x01}, p2, ByteOrderBig, WordOrderHighFirst); err == nil {
		t.Fatal("线圈使用 uint16 类型应报错")
	}
}

func TestParseErrors(t *testing.T) {
	p := point("f", RegisterHolding, DataTypeFloat32, 0, 2)
	if _, err := ParseRegisters([]byte{0x12, 0x34}, p, ByteOrderBig, WordOrderHighFirst); err == nil {
		t.Fatal("float32 数据不足应报错")
	}
	p2 := point("x", RegisterHolding, "badtype", 0, 1)
	if _, err := ParseRegisters([]byte{0x12, 0x34}, p2, ByteOrderBig, WordOrderHighFirst); err == nil {
		t.Fatal("未知类型应报错")
	}
	// 空数据
	if _, err := ParseRegisters(nil, p2, ByteOrderBig, WordOrderHighFirst); err == nil {
		t.Fatal("空数据应报错")
	}
}

func TestApplyScale(t *testing.T) {
	p := PointConfig{Scale: 0.1, Offset: -5}
	if got := ApplyScale(250, p); got != 20 {
		t.Fatalf("scale: got %v want 20", got)
	}
}
