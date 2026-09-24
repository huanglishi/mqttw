package modbus

import "testing"

func TestGroupPointsMergeContinuous(t *testing.T) {
	points := []PointConfig{
		{Name: "a", RegisterType: RegisterHolding, Address: 0, Quantity: 1},
		{Name: "b", RegisterType: RegisterHolding, Address: 1, Quantity: 1},
		{Name: "c", RegisterType: RegisterHolding, Address: 2, Quantity: 2},
	}
	groups := groupPoints(points)
	if len(groups) != 1 {
		t.Fatalf("期望 1 组，实际 %d", len(groups))
	}
	g := groups[0]
	if len(g.batches) != 1 {
		t.Fatalf("期望 1 批，实际 %d", len(g.batches))
	}
	b := g.batches[0]
	if b.start != 0 || b.quantity != 4 || len(b.points) != 3 {
		t.Fatalf("批参数错误: start=%d qty=%d points=%d", b.start, b.quantity, len(b.points))
	}
}

func TestGroupPointsSplitOnGap(t *testing.T) {
	points := []PointConfig{
		{Name: "a", RegisterType: RegisterHolding, Address: 0, Quantity: 1},
		{Name: "b", RegisterType: RegisterHolding, Address: 10, Quantity: 1},
		{Name: "c", RegisterType: RegisterHolding, Address: 11, Quantity: 1},
	}
	groups := groupPoints(points)
	g := groups[0]
	if len(g.batches) != 2 {
		t.Fatalf("期望 2 批（地址 0 与 10-11 不连续），实际 %d", len(g.batches))
	}
	if g.batches[0].start != 0 || g.batches[1].start != 10 {
		t.Fatalf("批起始地址错误: %d, %d", g.batches[0].start, g.batches[1].start)
	}
}

func TestGroupPointsSplitOnLimit(t *testing.T) {
	// 126 个连续寄存器 → 拆 2 批（上限 125）
	points := make([]PointConfig, 0, 126)
	for i := 0; i < 126; i++ {
		points = append(points, PointConfig{Name: "p", RegisterType: RegisterHolding, Address: int32(i), Quantity: 1})
	}
	groups := groupPoints(points)
	if len(groups[0].batches) != 2 {
		t.Fatalf("126 个寄存器期望拆 2 批，实际 %d", len(groups[0].batches))
	}
	if groups[0].batches[0].quantity != 125 || groups[0].batches[1].quantity != 1 {
		t.Fatalf("拆分数量错误: %d, %d", groups[0].batches[0].quantity, groups[0].batches[1].quantity)
	}
}

func TestGroupPointsCoilSeparate(t *testing.T) {
	points := []PointConfig{
		{Name: "c", RegisterType: RegisterCoil, Address: 0, Quantity: 1},
		{Name: "d", RegisterType: RegisterCoil, Address: 1, Quantity: 1},
		{Name: "h", RegisterType: RegisterHolding, Address: 0, Quantity: 1},
	}
	groups := groupPoints(points)
	if len(groups) != 2 {
		t.Fatalf("期望 2 组（coil/holding），实际 %d", len(groups))
	}
	// coil 组合并为一批
	for _, g := range groups {
		if g.regType == RegisterCoil && (len(g.batches) != 1 || g.batches[0].quantity != 2) {
			t.Fatalf("coil 组合并错误: %+v", g)
		}
	}
}

func TestSliceChunkRegisterOffset(t *testing.T) {
	b := batch{
		regType: RegisterHolding,
		points: []PointConfig{
			{Name: "a", Quantity: 1},
			{Name: "b", Quantity: 2},
			{Name: "c", Quantity: 1},
		},
	}
	raw := []byte{0x00, 0x01, 0x00, 0x02, 0x00, 0x03, 0x00, 0x04, 0x00, 0x05}
	// b 偏移 = 1*2 = 2 → raw[2:6] = 00 02 00 03
	chunk, err := sliceChunk(raw, b, PointConfig{Name: "b", Quantity: 2}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunk) != 4 || chunk[0] != 0x00 || chunk[1] != 0x02 || chunk[2] != 0x00 || chunk[3] != 0x03 {
		t.Fatalf("寄存器切片错误: % x", chunk)
	}
}

func TestSliceChunkCoilBit(t *testing.T) {
	b := batch{regType: RegisterCoil, points: []PointConfig{{Name: "c0"}, {Name: "c1"}, {Name: "c2"}}}
	// 线圈 0/1/2 全 ON → 第 1 字节 0x07；coil2 在 bit2
	chunk, err := sliceChunk([]byte{0x07}, b, PointConfig{Name: "c2"}, 2)
	if err != nil {
		t.Fatal(err)
	}
	if chunk[0] != 0x01 {
		t.Fatalf("coil2 位提取错误: % x", chunk)
	}
	// coil8 在第 2 字节 bit0（raw 覆盖 16 个线圈）
	chunk, err = sliceChunk([]byte{0x07, 0x01}, b, PointConfig{Name: "c8"}, 8)
	if err != nil {
		t.Fatal(err)
	}
	if chunk[0] != 0x01 {
		t.Fatalf("coil8 位提取错误: % x", chunk)
	}
	// 越界：raw 只有 1 字节，但点位在第 8 位 → 应报错
	if _, err = sliceChunk([]byte{0x07}, b, PointConfig{Name: "c8"}, 8); err == nil {
		t.Fatal("线圈字节越界应报错")
	}
}
