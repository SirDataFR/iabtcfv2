package iabtcfv2

type Bits struct {
	Position uint
	Bytes    []byte
}

var (
	bytePows = []byte{128, 64, 32, 16, 8, 4, 2, 1}
)

func NewBits(bytes []byte) *Bits {
	return &Bits{Bytes: bytes, Position: 0}
}

func (b *Bits) ReadBool() bool {
	byteIndex := b.Position / 8
	bitIndex := b.Position % 8
	b.Position++

	return (b.Bytes[byteIndex] & bytePows[bitIndex]) != 0
}

func (b *Bits) WriteBool(v bool) {
	byteIndex := b.Position / 8
	shift := (byteIndex+1)*8 - b.Position - 1
	b.Position++

	if v {
		b.Bytes[byteIndex] |= 1 << uint(shift)
	} else {
		b.Bytes[byteIndex] &^= 1 << uint(shift)
	}
}

func (b *Bits) ReadInt(n uint) int {
	v := 0
	for i, shift := uint(0), n-1; i < n; i++ {
		if b.ReadBool() {
			v += 1 << uint(shift)
		}
		shift--
	}

	return v
}

func (b *Bits) WriteInt(v int, n uint) {
	b.WriteNumber(int64(v), n)
}

func (b *Bits) WriteNumber(v int64, n uint) {
	startOffset := int(b.Position)
	for i := int(n) - 1; i >= 0; i-- {
		index := startOffset + i
		byteIndex := index / 8
		shift := (byteIndex+1)*8 - index - 1
		b.Bytes[byteIndex] |= byte(v%2) << uint(shift)
		v /= 2
	}

	b.Position += n
}
