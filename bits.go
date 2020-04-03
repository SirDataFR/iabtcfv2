package iabtcfv2

type Bits struct {
	position uint
	bytes    []byte
}

var (
	bytePows = []byte{128, 64, 32, 16, 8, 4, 2, 1}
)

func newBits(bytes []byte) *Bits {
	return &Bits{bytes: bytes, position: 0}
}

func (b *Bits) readBool() bool {
	byteIndex := b.position / 8
	bitIndex := b.position % 8
	b.position++

	return (b.bytes[byteIndex] & bytePows[bitIndex]) != 0
}

func (b *Bits) writeBool(v bool) {
	byteIndex := b.position / 8
	shift := (byteIndex+1)*8 - b.position - 1
	b.position++

	if v {
		b.bytes[byteIndex] |= 1 << uint(shift)
	} else {
		b.bytes[byteIndex] &^= 1 << uint(shift)
	}
}

func (b *Bits) readInt(n uint) int {
	v := 0
	for i, shift := uint(0), n-1; i < n; i++ {
		if b.readBool() {
			v += 1 << uint(shift)
		}
		shift--
	}

	return v
}

func (b *Bits) writeInt(v int, n uint) {
	b.writeNumber(int64(v), n)
}

func (b *Bits) writeNumber(v int64, n uint) {
	startOffset := int(b.position)
	for i := int(n) - 1; i >= 0; i-- {
		index := startOffset + i
		byteIndex := index / 8
		shift := (byteIndex+1)*8 - index - 1
		b.bytes[byteIndex] |= byte(v%2) << uint(shift)
		v /= 2
	}

	b.position += n
}
