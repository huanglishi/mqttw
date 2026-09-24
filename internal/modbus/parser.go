package modbus

import (
	"encoding/binary"
	"fmt"
	"math"
)

// ParseRegisters 将 Modbus 读取到的原始寄存器字节解析为物理值。
//
// raw: 读取返回的原始字节（寄存器按 Modbus 标准为大端字节序，1 个寄存器占 2 字节，
// 线圈/离散输入类为按位打包的字节流，此时取第 1 字节的最低 bit）。
// point: 点位配置（寄存器类型/数据类型/数量/缩放）。
// byteOrder/wordOrder: 设备级字节序与字序配置。
//
// 返回的值为换算前 raw 经类型解析后的数值（scale/offset 由上层统一应用）。
func ParseRegisters(raw []byte, point PointConfig, byteOrder, wordOrder string) (float64, error) {
	switch point.RegisterType {
	case RegisterCoil, RegisterDiscreteInput:
		if len(raw) < 1 {
			return 0, fmt.Errorf("寄存器数据为空")
		}
		if point.DataType != "" && point.DataType != DataTypeBit {
			return 0, fmt.Errorf("线圈/离散输入仅支持 bit 类型，当前为 %s", point.DataType)
		}
		// 批量读取时返回按位打包字节，点位偏移处的 bit 已由上层切片为单字节
		return float64(raw[0] & 0x01), nil
	}

	if len(raw) < 2 {
		return 0, fmt.Errorf("寄存器数据不足 2 字节")
	}

	switch point.DataType {
	case DataTypeUint16:
		return float64(u16(raw, byteOrder)), nil
	case DataTypeInt16:
		return float64(int16(u16(raw, byteOrder))), nil
	case DataTypeUint32:
		buf, err := assembleWords(raw, byteOrder, wordOrder)
		if err != nil {
			return 0, err
		}
		return float64(binary.BigEndian.Uint32(buf)), nil
	case DataTypeInt32:
		buf, err := assembleWords(raw, byteOrder, wordOrder)
		if err != nil {
			return 0, err
		}
		return float64(int32(binary.BigEndian.Uint32(buf))), nil
	case DataTypeFloat32:
		buf, err := assembleWords(raw, byteOrder, wordOrder)
		if err != nil {
			return 0, err
		}
		return float64(math.Float32frombits(binary.BigEndian.Uint32(buf))), nil
	case DataTypeBCD:
		return parseBCD(raw, byteOrder)
	default:
		return 0, fmt.Errorf("不支持的数据类型 %q", point.DataType)
	}
}

// ApplyScale 应用缩放系数与偏移：real = raw*scale + offset
func ApplyScale(raw float64, point PointConfig) float64 {
	return raw*point.Scale + point.Offset
}

// u16 按字节序读取一个 word（2 字节）
func u16(raw []byte, byteOrder string) uint16 {
	if byteOrder == ByteOrderLittle {
		return uint16(raw[0]) | uint16(raw[1])<<8
	}
	return binary.BigEndian.Uint16(raw)
}

// assembleWords 将多个寄存器的原始字节按字序与字节序组装为 4 字节大端缓冲。
// Modbus 寄存器本身按大端传输（高字节在前），因此：
//   - wordOrder 决定寄存器排列顺序（high_first: 首寄存器为高字）
//   - byteOrder 决定每个 word 内部字节是否需要交换（little 表示低字节在前）
func assembleWords(raw []byte, byteOrder, wordOrder string) ([]byte, error) {
	if len(raw) < 4 {
		return nil, fmt.Errorf("32 位数据类型需要 2 个寄存器（4 字节），实际 %d 字节", len(raw))
	}
	w0h, w0l := raw[0], raw[1]
	w1h, w1l := raw[2], raw[3]

	var b0, b1, b2, b3 byte
	switch wordOrder {
	case WordOrderHighFirst:
		b0, b1 = pick(w0h, w0l, byteOrder)
		b2, b3 = pick(w1h, w1l, byteOrder)
	case WordOrderLowFirst:
		b0, b1 = pick(w1h, w1l, byteOrder)
		b2, b3 = pick(w0h, w0l, byteOrder)
	default:
		return nil, fmt.Errorf("不支持的字序 %q", wordOrder)
	}
	return []byte{b0, b1, b2, b3}, nil
}

// pick 按字节序取出 word 的两个字节（big: 高字节在前）
func pick(hi, lo byte, byteOrder string) (byte, byte) {
	if byteOrder == ByteOrderLittle {
		return lo, hi
	}
	return hi, lo
}

// parseBCD 解析 BCD 编码：每个半字节 0-9，quantity=1 → 4 位十进制；quantity=2 → 8 位十进制。
func parseBCD(raw []byte, byteOrder string) (float64, error) {
	buf := make([]byte, len(raw))
	copy(buf, raw)
	if byteOrder == ByteOrderLittle {
		for i := 0; i+1 < len(buf); i += 2 {
			buf[i], buf[i+1] = buf[i+1], buf[i]
		}
	}
	if len(buf) > 4 {
		return 0, fmt.Errorf("BCD 最多支持 2 个寄存器（4 字节），实际 %d 字节", len(buf))
	}
	var val uint64
	for _, b := range buf {
		hi, lo := b>>4, b&0x0f
		if hi > 9 || lo > 9 {
			return 0, fmt.Errorf("非法的 BCD 数据 %02X", b)
		}
		val = val*100 + uint64(hi)*10 + uint64(lo)
	}
	return float64(val), nil
}
